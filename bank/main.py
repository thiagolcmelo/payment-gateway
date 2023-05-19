from datetime import datetime


from fastapi import BackgroundTasks, FastAPI, Request, Response, status
import httpx
from pydantic import BaseModel

from db.memory import (
    MemoryDB,
    Card,
    PaymentStatus,
    create_memory_db,
)

################################################################################
from logging.config import dictConfig
import logging
from config.logger import LogConfig

dictConfig(LogConfig().dict())
logger = logging.getLogger("bank-simulator")
################################################################################


class PaymentRequest(BaseModel):
    amount: float
    currency: str
    purchase_time: datetime
    validation_method: str
    card: Card
    merchant: str


class PaymentResponse(BaseModel):
    id: str
    success: bool
    message: str


class UpdatePaymentRequest(BaseModel):
    id: str
    success: bool
    message: str


class UpdatePaymentResponse(BaseModel):
    acknowledge: bool


async def process_payment(payment_id: int, host: str, db_manager: MemoryDB) -> None:
    shopper = await db_manager.find_shopper_by_payment_id(payment_id)
    merchants = await db_manager.find_shopper_merchants(shopper)
    payment = await db_manager.find_payment_by_id(payment_id)

    success, message = True, "success"
    if shopper.balance < payment.amount:
        message = "not enough balance"
        success = False
    elif payment.merchant not in merchants:
        message = "merchant unauthorized"
        success = False

    async with httpx.AsyncClient() as client:
        json_data = {
            "id": payment.uuid_id,
            "success": success,
            "message": message,
        }
        r = await client.put(
            f"http://{host}:8080/payment", json=json_data, timeout=10.0
        )
        r_data = r.json()

        if r.status_code == httpx.codes.OK and r_data.get("acknowledge", False):
            await db_manager.decrement_shopper_balance(shopper, payment.amount)
            await db_manager.mark_payment_status(payment_id, PaymentStatus.SUCCESS)
            db_manager.logger.info(f"{payment.uuid_id} - SUCCESS")
        else:
            await db_manager.mark_payment_status(payment_id, PaymentStatus.FAIL)
            db_manager.logger.info(f"{payment.uuid_id} - FAIL")


app = FastAPI()


@app.on_event("startup")
async def startup():
    # Store the database connection in the application state
    app.state.db_connection = create_memory_db("data/shoppers.json")
    app.state.db_helper = MemoryDB(app.state.db_connection, logger)


@app.on_event("shutdown")
async def shutdown():
    # Close the database connection when the application is shutting down
    app.state.db_connection.close()


@app.post("/payment")
async def create_payment(
    payment_request: PaymentRequest,
    background_tasks: BackgroundTasks,
    req: Request,
    resp: Response,
) -> PaymentResponse:
    response = PaymentResponse(id="", success=False, message="error")

    try:
        card = await app.state.db_helper.fill_card_id(payment_request.card)
        if card.id is None:
            response.message = "card not found"
            raise Exception(response.message)
        shopper = await app.state.db_helper.find_shopper_by_card(card)

        if shopper is None:
            response.message = "card does not match a shopper"
            raise Exception(response.message)
        elif shopper.currency != payment_request.currency:
            response.message = "shopper currency is not correct"
            raise Exception(response.message)

        payment_id, payment_uuid = await app.state.db_helper.create_payment_for_shopper(
            shopper,
            card,
            payment_request.amount,
            payment_request.purchase_time,
            payment_request.validation_method,
            payment_request.merchant,
        )
        logger.info(f"{payment_uuid} - CREATED")

        if payment_id == "":
            response.message = "could not create payment"
            raise Exception(response.message)

        background_tasks.add_task(
            process_payment, payment_id, req.client.host, app.state.db_helper
        )
        await app.state.db_helper.mark_payment_status(payment_id, PaymentStatus.PENDING)

        response.id = payment_uuid
        response.message = "payment request created"
        response.success = True
        resp.status_code = status.HTTP_201_CREATED
        logger.info(f"{payment_uuid} - PENDING")
    except Exception as err:
        resp.status_code = status.HTTP_400_BAD_REQUEST
        logger.error(err)

    return response


@app.put("/payment")
async def update_payment(
    update: UpdatePaymentRequest,
    resp: Response,
) -> UpdatePaymentResponse:
    id = update.id
    message = update.message
    success = update.success
    logger.info(f"acknowledging message: ({id}, {success}, {message})")
    resp.status_code = status.HTTP_200_OK
    return UpdatePaymentResponse(acknowledge=True)

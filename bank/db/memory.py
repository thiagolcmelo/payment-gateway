import asyncio
from datetime import datetime
from enum import Enum
import json
from logging import Logger
import sqlite3
from typing import List, Optional, Tuple
import uuid
import httpx

from pydantic import BaseModel


class PaymentStatus(Enum):
    CREATED = 0
    PENDING = 1
    SUCCESS = 2
    FAIL = 3


class Shopper(BaseModel):
    id: int
    name: str
    description: str
    currency: str
    balance: float


class Card(BaseModel):
    id: Optional[int]
    number: str
    name: str
    expire_month: int
    expire_year: int
    cvv: int


class Payment(BaseModel):
    id: int
    uuid_id: str
    amount: float
    purchase_time: datetime
    validation_method: str
    card_id: int
    merchant: str
    shopper_id: int
    created_at: datetime
    status: PaymentStatus


class MemoryDB:
    def __init__(self, conn: sqlite3.Connection, logger: Logger) -> None:
        self.conn = conn
        self.database_lock = asyncio.Lock()
        self.logger = logger

    async def fill_card_id(self, card: Card) -> Card:
        await self.database_lock.acquire()
        try:
            cursor = self.conn.cursor()
            cursor.execute(
                "SELECT id FROM cards WHERE number=? AND name=? AND expire_month=? AND expire_year=? AND cvv=?",
                (
                    card.number,
                    card.name,
                    card.expire_month,
                    card.expire_year,
                    card.cvv,
                ),
            )
            card_id = cursor.fetchone()
            if card_id is not None and len(card_id) > 0:
                card.id = card_id[0]
            else:
                card.id = None
        finally:
            self.database_lock.release()
        return card

    async def decrement_shopper_balance(self, shopper: Shopper, amount: float) -> None:
        await self.database_lock.acquire()
        try:
            new_balance = shopper.balance - amount
            cursor = self.conn.cursor()
            cursor.execute(
                "UPDATE shoppers SET balance=? WHERE id=?", (new_balance, shopper.id)
            )
            self.conn.commit()
        finally:
            self.database_lock.release()

    async def find_shopper_by_card(self, card: Card) -> Shopper:
        await self.database_lock.acquire()
        shopper = None
        try:
            cursor = self.conn.cursor()
            cursor.execute("SELECT shopper_id FROM cards WHERE id=?", (card.id,))
            shopper_id = cursor.fetchone()
            if shopper_id is not None and len(shopper_id) > 0:
                shopper_id = shopper_id[0]
                cursor.execute(
                    "SELECT id, name, description, currency, balance FROM shoppers WHERE id=?",
                    (shopper_id,),
                )
                _shopper = cursor.fetchone()
                if _shopper is not None:
                    shopper = Shopper(
                        id=int(_shopper[0]),
                        name=_shopper[1],
                        description=_shopper[2],
                        currency=_shopper[3],
                        balance=float(_shopper[4]),
                    )
        finally:
            self.database_lock.release()
        return shopper

    async def find_shopper_by_payment_id(self, payment_id: int) -> Shopper:
        await self.database_lock.acquire()
        shopper = None
        try:
            cursor = self.conn.cursor()
            cursor.execute("SELECT shopper_id FROM payments WHERE id=?", (payment_id,))
            shopper_id = cursor.fetchone()
            if shopper_id is not None and len(shopper_id) > 0:
                shopper_id = shopper_id[0]
                cursor.execute(
                    "SELECT id, name, description, currency, balance FROM shoppers WHERE id=?",
                    (shopper_id,),
                )
                _shopper = cursor.fetchone()
                if _shopper is not None:
                    shopper = Shopper(
                        id=int(_shopper[0]),
                        name=_shopper[1],
                        description=_shopper[2],
                        currency=_shopper[3],
                        balance=float(_shopper[4]),
                    )
        finally:
            self.database_lock.release()
        return shopper

    async def find_payment_by_id(self, payment_id: int) -> Payment:
        await self.database_lock.acquire()
        payment = None
        try:
            cursor = self.conn.cursor()
            cursor.execute(
                "SELECT id, uuid_id, amount, purchase_time, validation_method, card_id, merchant, shopper_id, created_at, status FROM payments WHERE id=?",
                (payment_id,),
            )
            row = cursor.fetchone()
            if row is not None and len(row) > 0:
                payment = Payment(
                    id=int(row[0]),
                    uuid_id=row[1],
                    amount=float(row[2]),
                    purchase_time=datetime.strptime(row[3], "%Y%m%dT%H%M%S.%f"),
                    validation_method=row[4],
                    card_id=int(row[5]),
                    merchant=row[6],
                    shopper_id=int(row[7]),
                    created_at=datetime.strptime(row[8], "%Y%m%dT%H%M%S.%f"),
                    status=PaymentStatus(int(row[9])),
                )
        finally:
            self.database_lock.release()
        return payment

    async def find_shopper_merchants(self, shopper: Shopper) -> List[str]:
        await self.database_lock.acquire()
        merchants = []
        try:
            cursor = self.conn.cursor()
            cursor.execute(
                "SELECT merchant FROM auto_approve_merchants WHERE shopper_id=?",
                (shopper.id,),
            )
            _merchats = cursor.fetchall()
            for row in _merchats:
                merchants.append(row[0])
        finally:
            self.database_lock.release()
        return merchants

    async def mark_payment_status(self, payment_id: int, status: PaymentStatus) -> None:
        await self.database_lock.acquire()
        try:
            cursor = self.conn.cursor()
            cursor.execute(
                "UPDATE payments SET status=? WHERE id=?",
                (int(status.value), payment_id),
            )
            self.conn.commit()
        finally:
            self.database_lock.release()

    async def create_payment_for_shopper(
        self,
        shopper: Shopper,
        card: Card,
        amount: float,
        purchase_time: datetime,
        validation_method: str,
        merchant: str,
    ) -> Tuple[str]:
        await self.database_lock.acquire()
        payment_id, payment_uuid = None, None

        try:
            id = uuid.uuid1()
            cursor = self.conn.cursor()
            cursor.execute(
                """
                INSERT INTO payments
                    (uuid_id, amount, purchase_time, validation_method, card_id, merchant, shopper_id, created_at, status)
                    VALUES
                    (?, ?, ?, ?, ?, ?, ?, ?, ?)""",
                (
                    str(id),
                    amount,
                    purchase_time.strftime("%Y%m%dT%H%M%S.%f"),
                    validation_method,
                    card.id,
                    merchant,
                    shopper.id,
                    datetime.now().strftime("%Y%m%dT%H%M%S.%f"),
                    int(PaymentStatus.CREATED.value),
                ),
            )
            payment_id = cursor.lastrowid
            self.conn.commit()

            cursor.execute(
                "SELECT uuid_id FROM payments WHERE id=?",
                (payment_id,),
            )
            row = cursor.fetchone()
            if row is not None:
                payment_uuid = row[0]
        finally:
            self.database_lock.release()

        return payment_id, payment_uuid

    async def process_payment(self, payment_id: int, host: str) -> None:
        shopper = await self.find_shopper_by_payment_id(payment_id)
        merchants = await self.find_shopper_merchants(shopper)
        payment = await self.find_payment_by_id(payment_id)

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

            if (
                r.status_code == httpx.codes.OK
                and message is None
                and r_data.get("acknowledge", False)
            ):
                await self.decrement_shopper_balance(shopper, payment.amount)
                await self.mark_payment_status(payment_id, PaymentStatus.SUCCESS)
                self.logger.info(f"{payment.uuid_id} - SUCCESS")
            else:
                await self.mark_payment_status(payment_id, PaymentStatus.FAIL)
                self.logger.info(f"{payment.uuid_id} - FAIL")


def create_memory_db(json_data: str) -> sqlite3.Connection:
    dummy_data = []
    with open(json_data, "r") as f:
        dummy_data = json.loads(f.read())

    # Create a connection to the in-memory SQLite database
    conn = sqlite3.connect(":memory:")

    # Perform necessary database setup, such as creating tables or indexes
    cursor = conn.cursor()
    cursor.execute(
        """
    CREATE TABLE IF NOT EXISTS shoppers (
        id INTEGER PRIMARY KEY,
        name TEXT,
        description TEXT,
        currency TEXT,
        balance REAL
    )"""
    )
    cursor.execute(
        """
    CREATE TABLE IF NOT EXISTS cards (
        id INTEGER PRIMARY KEY,
        number TEXT,
        name TEXT,
        expire_month INTEGER,
        expire_year INTEGER,
        cvv INTEGER,
        shopper_id INTEGER,
        FOREIGN KEY (shopper_id) REFERENCES shoppers(id)
    )"""
    )
    cursor.execute(
        """
    CREATE TABLE IF NOT EXISTS payments (
        id INTEGER PRIMARY KEY,
        uuid_id TEXT,
        amount REAL,
        purchase_time TEXT,
        validation_method TEXt,
        card_id INTEGER,
        merchant TEXT,
        shopper_id INTEGER,
        created_at TEXT,
        status INTEGER,
        FOREIGN KEY (card_id) REFERENCES cards(id),
        FOREIGN KEY (shopper_id) REFERENCES shoppers(id)
    )"""
    )
    cursor.execute(
        """
    CREATE TABLE IF NOT EXISTS auto_approve_merchants (
        id TEXT PRIMARY KEY,
        merchant TEXT,
        shopper_id INTEGER,
        FOREIGN KEY (shopper_id) REFERENCES shoppers(id)
    )"""
    )

    for shopper in dummy_data:
        cursor.execute(
            "INSERT INTO shoppers (name, description, currency, balance) VALUES (?, ?, ?, ?)",
            (
                shopper["name"],
                shopper["description"],
                shopper["currency"],
                shopper["balance"],
            ),
        )
        shopper_id = cursor.lastrowid
        cursor.execute(
            "INSERT INTO cards (number, name, expire_month, expire_year, cvv, shopper_id) VALUES (?, ?, ?, ?, ?, ?)",
            (
                shopper["card"]["number"],
                shopper["card"]["name"],
                shopper["card"]["expire_month"],
                shopper["card"]["expire_year"],
                shopper["card"]["cvv"],
                shopper_id,
            ),
        )
        for merchant in shopper["auto_approve"]:
            cursor.execute(
                "INSERT INTO auto_approve_merchants (merchant, shopper_id) VALUES (?, ?)",
                (
                    merchant,
                    shopper_id,
                ),
            )

    conn.commit()
    return conn

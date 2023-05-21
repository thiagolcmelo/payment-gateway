# Bank Simulator

This is a simple **Acquiring Bank** simulator. It comes with some **Shoppers** in memory and uses them to decide upon payment requests.

There is one endpoint exposed over HTTP:

- `POST /payment HTTP/1.1` to create a payment. If the payload is correct, it will reply with a success message and trigger a background task to process the payment.

The background task will attempt to inform to the **Payment Gateway** if the processing was successful. If it fail in contacting the **Payment Gateway** or if it does not receive a valid reply acknowledging the message, it sets the payment to fail.

The idea is that the payload (please see below) must contain all necessary information to identify a **Shopper**. There is for instance a field `validation_method` to simulate a way to contact the **Shopper** and verify the purchase, for instance **sms**, **push**, **email**, etc. There is very few validation as it is mainly conceptual.

A payment is expected to have the following data:

```json
{
    "amount": 10.0,
    "currency": "USD",
    "purchase_time": "2023-05-18T10:00:00.000",
    "validation_method": "sms",
    "card": {
        "number": "1111-2222-3333-4444",
        "name": "shopper 0",
        "expire_month": 10,
        "expire_year": 2050,
        "cvv": 123
    },
    "merchant": "Merchant 0 Ltd."
}
```

The name of the **Merchant** is used in a *auto approve* mechanism. If the **SHopper** has a **Merchant** among those set to *auto approve*, the payment proceeds.

## Data Model

Internally there is an in memory database with the following structure:

```sql
CREATE TABLE IF NOT EXISTS shoppers (
    id INTEGER PRIMARY KEY,
    name TEXT,
    description TEXT,
    currency TEXT,
    balance REAL
);

CREATE TABLE IF NOT EXISTS cards (
    id INTEGER PRIMARY KEY,
    number TEXT,
    name TEXT,
    expire_month INTEGER,
    expire_year INTEGER,
    cvv INTEGER,
    shopper_id INTEGER,
    FOREIGN KEY (shopper_id) REFERENCES shoppers(id)
);

CREATE TABLE IF NOT EXISTS payments (
    id INTEGER PRIMARY KEY,
    uuid_id TEXT,
    amount REAL,
    currency TEXT,
    purchase_time TEXT,
    validation_method TEXt,
    card_id INTEGER,
    merchant TEXT,
    shopper_id INTEGER,
    created_at TEXT,
    status INTEGER,
    FOREIGN KEY (card_id) REFERENCES cards(id),
    FOREIGN KEY (shopper_id) REFERENCES shoppers(id)
);

CREATE TABLE IF NOT EXISTS auto_approve_merchants (
    id TEXT PRIMARY KEY,
    merchant TEXT,
    shopper_id INTEGER,
    FOREIGN KEY (shopper_id) REFERENCES shoppers(id)
)
```




## Testing

Please run it as follows:

```bash
$ uvicorn main:app --reload
```

Then please try to request a payment:

```bash
$ curl -H "Content-Type: application/json" -X POST -d @data/request.json http://127.0.0.1:8000/payment
{"id":"5bd5aa5a-f5e4-11ed-84f9-8c859093fdeb","success":true,"message":"payment request created"}
```

In the stdout, something similar to this will be printed:

```
INFO:     [2023-05-19 02:32:57] create_payment: 1541d784-f5e5-11ed-8a90-8c859093fdeb - CREATED
INFO:     [2023-05-19 02:32:57] create_payment: 1541d784-f5e5-11ed-8a90-8c859093fdeb - PENDING
INFO:     127.0.0.1:62857 - "POST /payment HTTP/1.1" 201 Created
INFO:     [2023-05-19 02:32:57] update_payment: acknowledging message: (1541d784-f5e5-11ed-8a90-8c859093fdeb, success)
INFO:     127.0.0.1:62858 - "PUT /payment HTTP/1.1" 200 OK
INFO:     [2023-05-19 02:32:57] process_payment: 1541d784-f5e5-11ed-8a90-8c859093fdeb - SUCCESS
```

## **Acquiring Bank** configuration

It will determine if a **Shopper** matches a given card, if they authorized a given **Merchant**, if it has enough balance, if the currency is correct, and so on.

The content is stored in `data/shoppers.json`, but for simplicity:

```json
[
    {
        "name": "shopper 0",
        "description": "shopper with valid card and large balance in USD",
        "card": {
            "number": "1111-2222-3333-4444",
            "name": "shopper 0",
            "expire_month": 10,
            "expire_year": 2050,
            "cvv": 123
        },
        "auto_approve": [
            "Merchant 0 Ltd.",
            "Merchant 1 Ltd.",
            "Merchant 2 Ltd.",
            "Merchant 3 Ltd.",
            "Merchant 4 Ltd."
        ],
        "currency": "USD",
        "balance": 1000000.00
    },
    {
        "name": "shopper 1",
        "description": "shopper with valid card and small balance in USD",
        "card": {
            "number": "5555-6666-7777-8888",
            "name": "shopper 1",
            "expire_month": 10,
            "expire_year": 2040,
            "cvv": 456
        },
        "auto_approve": [
            "Merchant 0 Ltd.",
            "Merchant 1 Ltd.",
            "Merchant 2 Ltd."
        ],
        "currency": "USD",
        "balance": 100.00
    },
    {
        "name": "shopper 2",
        "description": "shopper with valid card and medium balance in USD",
        "card": {
            "number": "9999-1010-1111-1212",
            "name": "shopper 2",
            "expire_month": 3,
            "expire_year": 2045,
            "cvv": 789
        },
        "auto_approve": [
            "Merchant 3 Ltd.",
            "Merchant 4 Ltd."
        ],
        "currency": "USD",
        "balance": 10000.00
    },
    {
        "name": "shopper 4",
        "description": "shopper with valid card and large balance in EUR",
        "card": {
            "number": "1313-1414-1515-1616",
            "name": "shopper 5",
            "expire_month": 1,
            "expire_year": 2070,
            "cvv": 987
        },
        "auto_approve": [
            "Merchant 0 Ltd."
        ],
        "currency": "EUR",
        "balance": 1000000.00
    },
    {
        "name": "shopper 5",
        "description": "shopper with valid card and large balance in GBP",
        "card": {
            "number": "1717-1818-1919-2020",
            "name": "shopper 5",
            "expire_month": 1,
            "expire_year": 2070,
            "cvv": 987
        },
        "auto_approve": [
            "Merchant 0 Ltd.",
            "Merchant 1 Ltd.",
            "Merchant 2 Ltd.",
            "Merchant 3 Ltd.",
            "Merchant 4 Ltd."
        ],
        "currency": "EUR",
        "balance": 1000000.00
    }
]
```
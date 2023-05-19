# Bank Simulator

This is a simple **Acquiring Bank** simulator. It comes with some **Shoppers** in memory and uses them to decide upon payment requests.

There are two endpoints exposed over HTTP:

- `POST /payment HTTP/1.1` to create a payment.
- `PUT /payment HTTP/1.1` only for testing, allows the background task that processes payments to pretend it is receiving an acknowledge message from the **Payment Gateway**.

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

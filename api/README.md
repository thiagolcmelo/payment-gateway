# Payment API

This API connects all the services.

## Testing

This service is not covered by tests yet.

It can be executed manually as follows:

```bash
$ go run handlers.go middleware.go main.go
```

If all the other services are running on their default ports, the following commands will permit testing:

```bash
$ curl -u merchant0:password0 http://localhost:8080/login
{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im1lcmNoYW50MCIsImlkIjoiZTEyMTEzNTEtYmI5MS00NDFmLTllYTAtM2IyNDMxODlkZWM2IiwiZXhwIjoiMjAyMy0wNS0yMVQwMDowOToyMi4zNDAwNiswMTowMCJ9.lFs4rj02Zolnc0U9tne7rSBXcYMNfb9YQuk4pG0Zmxg"}
```

Or

```bash
$ TOKEN=$(curl -u merchant0:password0 http://localhost:8080/login 2>/dev/null | jq --raw-output .token)
$ curl -X POST -d @data/payment.json -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" http://127.0.0.1:8080/payment
{"id":"2b862843-fe6a-4798-bd9f-bf1de4fc385b"}
```

Looking at the server logs, two requests should happen:

```
[GIN] 2023/05/20 - 00:19:55 | 200 |   17.286534ms |       127.0.0.1 | POST     "/payment"
[GIN] 2023/05/20 - 00:19:55 | 200 |   17.923023ms |       127.0.0.1 | PUT      "/payment"
```

The first is the payment request and the second is the bank confirmation.

After that, it is possible to read the payment information:

```bash
$ curl -X GET -H "Authorization: Bearer $TOKEN" http://127.0.0.1:8080/payment/2b862843-fe6a-4798-bd9f-bf1de4fc385b 2>/dev/null | jq .
{
  "id": "2b862843-fe6a-4798-bd9f-bf1de4fc385b",
  "merchant_id": "e1211351-bb91-441f-9ea0-3b243189dec6",
  "amount": 10,
  "currency": "USD",
  "purchase_time": "2023-05-18T10:00:00Z",
  "validation_method": "push",
  "card": {
    "number": "1111-2222-3333-4444",
    "name": "shopper 0",
    "expire_month": 10,
    "expire_year": 2050,
    "cvv": 123
  },
  "metadata": "shopper 0",
  "status": "SUCCESS",
  "bank_payment_id": "aa0dd29e-f69b-11ed-8560-8c859093fdeb",
  "bank_request_time": "2023-05-20T00:19:55.186Z",
  "bank_response_time": "2023-05-20T00:19:55.198Z",
  "bank_message": "success"
}
```


If everything goes well, the web server out put should be similar to:

![High Level](../images/golden-path-logs.png)
# Ledger

It is a conceptual Ledger for storing payments. It exposes the following endpoints over gRPC:

- `CreatePayment` to create a new payment.
- `ReadPayment` to get all information about a payment.
- `ReadPaymentUsingBankReference` to get all information about a payment using the bank reference.
- `UpdatePaymentToPending` to inform that a payment was sent to an **Acquiring Bank**.
- `UpdatePaymentToSuccess` to inform that a payment was successfully executed by an **Acquiring Bank**.
- `UpdatePaymentToFail` to set the payment as a failure, if refused by the bank, a message is expected to infrom the reason.

The validation of fields is very simple, it is not truly checking credit card number, or currencies. The idea is to get something minimal working, but that can be extended and improved later.

## Testing

### Unit tests

Please run the unit tests as follows:

```bash
$ go test -v ./... -coverprofile=coverage.out
```

And check the coverage (currently ~74%):

```bash
$ go tool cover -func=coverage.out
```

### Manual testing

These tests should later be automated as integrations tests, but for now please use the [grpcurl](https://github.com/fullstorydev/grpcurl) tool to verify that the endpoints are working as expected.

- **Create a valid payment**

```bash
$ grpcurl -plaintext -d '{"merchant_id": "e1211351-bb91-441f-9ea0-3b243189dec6", "amount": 150.0, "currency": "USD", "purchase_time_utc": "2023-05-18T05:00:10.000", "validation_method": "sms", "card": {"number": "1111-2222-3333-4444", "name": "name surname", "expire_month": 10, "expire_year": 2099, "cvv": 123}, "metadata": "shopper:123"}' "[::1]:50053" ledger.LedgerService/CreatePayment
{
  "id": "35947b97-1cf3-4f5c-aa23-3ee0b9734ff7"
}
```

- **Create an invalid payment - negative amount**

```bash
$ grpcurl -plaintext -d '{"merchant_id": "e1211351-bb91-441f-9ea0-3b243189dec6", "amount": -150.0, "currency": "USD", "purchase_time_utc": "2023-05-18T05:00:10.000", "validation_method": "sms", "card": {"number": "1111-2222-3333-4444", "name": "name surname", "expire_month": 10, "expire_year": 2099, "cvv": 123}, "metadata": "shopper:123"}' "[::1]:50053" ledger.LedgerService/CreatePayment
ERROR:
  Code: Unknown
  Message: negative amount
```

- **Create an invalid payment - card expired**

```bash
$ grpcurl -plaintext -d '{"merchant_id": "e1211351-bb91-441f-9ea0-3b243189dec6", "amount": 150.0, "currency": "USD", "purchase_time_utc": "2023-05-18T05:00:10.000", "validation_method": "sms", "card": {"number": "1111-2222-3333-4444", "name": "name surname", "expire_month": 10, "expire_year": 2020, "cvv": 123}, "metadata": "shopper:123"}' "[::1]:50053" ledger.LedgerService/CreatePayment
ERROR:
  Code: Unknown
  Message: invalid expiration
```

- **Read existing payment**

```bash
$ grpcurl -plaintext -d '{"id": "35947b97-1cf3-4f5c-aa23-3ee0b9734ff7"}' "[::1]:50053" ledger.LedgerService/ReadPayment
{
  "payment": {
    "id": "35947b97-1cf3-4f5c-aa23-3ee0b9734ff7",
    "merchantId": "e1211351-bb91-441f-9ea0-3b243189dec6",
    "amount": 150,
    "currency": "USD",
    "purchaseTimeUtc": "2023-05-18T05:00:10.000",
    "validationMethod": "sms",
    "card": {
      "number": "1111-2222-3333-4444",
      "name": "name surname",
      "expireMonth": 10,
      "expireYear": 2099,
      "cvv": 123
    },
    "metadata": "shopper:123",
    "bankPaymentId": "00000000-0000-0000-0000-000000000000",
    "bankRequestTimeUtc": "0001-01-01T00:00:00.000",
    "bankResponseTimeUtc": "0001-01-01T00:00:00.000"
  }
}
```

- **Read unexisting payment**

```bash
$ grpcurl -plaintext -d '{"id": "e1211351-bb91-441f-9ea0-3b243189dec6"}' "[::1]:50053" ledger.LedgerService/ReadPayment
ERROR:
  Code: Unknown
  Message: there is no payment with given id
```

- **Update payment to pending**

```bash
$ grpcurl -plaintext -d '{"id": "35947b97-1cf3-4f5c-aa23-3ee0b9734ff7", "bank_payment_id": "70580f0a-8478-4aba-8ccf-3de1e1df665c", "bank_request_time_utc": "2023-05-18T05:01:10.000"}' "[::1]:50053" ledger.LedgerService/UpdatePaymentToPending
{

}
```

```bash
$ grpcurl -plaintext -d '{"id": "35947b97-1cf3-4f5c-aa23-3ee0b9734ff7"}' "[::1]:50053" ledger.LedgerService/ReadPayment
{
  "payment": {
    "id": "35947b97-1cf3-4f5c-aa23-3ee0b9734ff7",
    "merchantId": "e1211351-bb91-441f-9ea0-3b243189dec6",
    "amount": 150,
    "currency": "USD",
    "purchaseTimeUtc": "2023-05-18T05:00:10.000",
    "validationMethod": "sms",
    "card": {
      "number": "1111-2222-3333-4444",
      "name": "name surname",
      "expireMonth": 10,
      "expireYear": 2099,
      "cvv": 123
    },
    "metadata": "shopper:123",
    "status": "PENDING",
    "bankPaymentId": "70580f0a-8478-4aba-8ccf-3de1e1df665c",
    "bankRequestTimeUtc": "2023-05-18T05:01:10.000",
    "bankResponseTimeUtc": "0001-01-01T00:00:00.000"
  }
}
```

- **Update payment to success**

```bash
$ grpcurl -plaintext -d '{"id": "35947b97-1cf3-4f5c-aa23-3ee0b9734ff7", "bank_payment_id": "70580f0a-8478-4aba-8ccf-3de1e1df665c", "bank_response_time_utc": "2023-05-18T05:02:10.000", "bank_message": "success"}' "[::1]:50053" ledger.LedgerService/UpdatePaymentToSuccess
{

}
```

```bash
$ grpcurl -plaintext -d '{"id": "35947b97-1cf3-4f5c-aa23-3ee0b9734ff7"}' "[::1]:50053" ledger.LedgerService/ReadPayment
{
  "payment": {
    "id": "35947b97-1cf3-4f5c-aa23-3ee0b9734ff7",
    "merchantId": "e1211351-bb91-441f-9ea0-3b243189dec6",
    "amount": 150,
    "currency": "USD",
    "purchaseTimeUtc": "2023-05-18T05:00:10.000",
    "validationMethod": "sms",
    "card": {
      "number": "1111-2222-3333-4444",
      "name": "name surname",
      "expireMonth": 10,
      "expireYear": 2099,
      "cvv": 123
    },
    "metadata": "shopper:123",
    "status": "SUCCESS",
    "bankPaymentId": "70580f0a-8478-4aba-8ccf-3de1e1df665c",
    "bankRequestTimeUtc": "2023-05-18T05:01:10.000",
    "bankResponseTimeUtc": "2023-05-18T05:02:10.000",
    "bankMessage": "success"
  }
}
```

- **Update payment to fail**

```bash
$ grpcurl -plaintext -d '{"merchant_id": "e1211351-bb91-441f-9ea0-3b243189dec6", "amount": 150.0, "currency": "USD", "purchase_time_utc": "2023-05-18T05:00:10.000", "validation_method": "sms", "card": {"number": "1111-2222-3333-4444", "name": "name surname", "expire_month": 10, "expire_year": 2029, "cvv": 123}, "metadata": "shopper:123"}' "[::1]:50053" ledger.LedgerService/CreatePayment
{
  "id": "dbb2c818-ec5f-4dac-ad2a-a6b021c981bf"
}
```

```bash
$ grpcurl -plaintext -d '{"id": "dbb2c818-ec5f-4dac-ad2a-a6b021c981bf", "bank_payment_id": "70580f0a-8478-4aba-8ccf-3de1e1df665c", "bank_response_time_utc": "2023-05-18T05:02:10.000", "bank_message": "could not reach shopper"}' "[::1]:50053" ledger.LedgerService/UpdatePaymentToFail
{

}
```

```bash
$ grpcurl -plaintext -d '{"id": "dbb2c818-ec5f-4dac-ad2a-a6b021c981bf"}' "[::1]:50053" ledger.LedgerService/ReadPayment
{
  "payment": {
    "id": "dbb2c818-ec5f-4dac-ad2a-a6b021c981bf",
    "merchantId": "e1211351-bb91-441f-9ea0-3b243189dec6",
    "amount": 150,
    "currency": "USD",
    "purchaseTimeUtc": "2023-05-18T05:00:10.000",
    "validationMethod": "sms",
    "card": {
      "number": "1111-2222-3333-4444",
      "name": "name surname",
      "expireMonth": 10,
      "expireYear": 2029,
      "cvv": 123
    },
    "metadata": "shopper:123",
    "status": "FAIL",
    "bankPaymentId": "70580f0a-8478-4aba-8ccf-3de1e1df665c",
    "bankRequestTimeUtc": "0001-01-01T00:00:00.000",
    "bankResponseTimeUtc": "2023-05-18T05:02:10.000",
    "bankMessage": "could not reach shopper"
  }
}
```

- **Read using bank payment id**

```bash
$ grpcurl -plaintext -d '{"merchant_id": "e1211351-bb91-441f-9ea0-3b243189dec6", "amount": 150.0, "currency": "USD", "purchase_time_utc": "2023-05-18T05:00:10.000", "validation_method": "sms", "card": {"number": "1111-2222-3333-4444", "name": "name surname", "expire_month": 10, "expire_year": 2029, "cvv": 123}, "metadata": "shopper:123"}' "[::1]:50053" ledger.LedgerService/CreatePayment
{
  "id": "ac5503cc-3018-4484-90e1-0bcc64c91f63"
}
```

```bash
$ grpcurl -plaintext -d '{"id": "ac5503cc-3018-4484-90e1-0bcc64c91f63", "bank_payment_id": "70580f0a-8478-4aba-8ccf-3de1e1df665c", "bank_response_time_utc": "2023-05-18T05:02:10.000", "bank_message": "success"}' "[::1]:50053" ledger.LedgerService/UpdatePaymentToSuccess
{

}
```

```bash
$ grpcurl -plaintext -d '{"id": "70580f0a-8478-4aba-8ccf-3de1e1df665c"}' "[::1]:50053" ledger.LedgerService/ReadPaymentUsingBankReference
{
  "payment": {
    "id": "ac5503cc-3018-4484-90e1-0bcc64c91f63",
    "merchantId": "e1211351-bb91-441f-9ea0-3b243189dec6",
    "amount": 150,
    "currency": "USD",
    "purchaseTimeUtc": "2023-05-18T05:00:10.000",
    "validationMethod": "sms",
    "card": {
      "number": "1111-2222-3333-4444",
      "name": "name surname",
      "expireMonth": 10,
      "expireYear": 2029,
      "cvv": 123
    },
    "metadata": "shopper:123",
    "status": "SUCCESS",
    "bankPaymentId": "70580f0a-8478-4aba-8ccf-3de1e1df665c",
    "bankRequestTimeUtc": "0001-01-01T00:00:00.000",
    "bankResponseTimeUtc": "2023-05-18T05:02:10.000",
    "bankMessage": "success"
  }
}
```
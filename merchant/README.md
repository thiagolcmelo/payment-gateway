# Merchant

It is a conceptual service to store information related to **Merchants**.

A **Merchant** is modeled here as follows:

- `ID`:       **string**
- `Username`: **string**, unique, not null
- `Password`: **string**, not null, encrypted before storage
- `Name`: **string**, not null
- `Active`:   **bool**
- `MaxQPS`:   **int**, greater or equal to zero

There is an in memory implementation of a possible `Storage` service that can be used to persist **Merchants** in disk.

The service provides 4 endpoints over gRPC:

- `GetMerchant` to retrieve all the information for a **Merchant** given its `ID`.
- `GetQPS` to retrieve `MaxQPS` information for a **Merchant** given its `ID`.
- `MerchantActive` to retrieve `Active` information for a **Merchant** given its `ID`.
- `MerchantExists` to check if a `Username` and `Password` matches any **Merchant**, if so, its `ID` is returned.

## Testing

### Unit test

The initialization will create sample **Merchants** in the memory storage that can be used for testing. Their details can be found in `/data/merchants.json`.

Please run the tests as follows:

```bash
$ go test -v ./... -coverprofile=coverage.out
```

and verify the coverage (currently ~96%):

```bash
$ go tool cover -func=coverage.out
```

### Manual testing

The service can be exposed through gRPC:

```bash
$ go run main.go --ip-version=4 --host=0.0.0.0 --port=50051
```

In a different terminal, please use the [grpcurl](https://github.com/fullstorydev/grpcurl) tool to verify that the endpoints are working as expected.

#### GetMerchant

- Existing merchant:

```bash
$ grpcurl -plaintext -d '{"id": "e1211351-bb91-441f-9ea0-3b243189dec6"}' "0.0.0.0:50051" merchant.MerchantService/GetMerchant | jq .
{
  "id": "e1211351-bb91-441f-9ea0-3b243189dec6",
  "username": "merchant0",
  "password": "$2a$10$TfTXwd7PA.rUrioJrkPbEutsp8WxvJrFDPfOgtTRwolNN3O7m0zKS",
  "name": "Merchant 0 Ltd.",
  "active": true,
  "maxQps": 100
}
```

- Non existing merchant

```bash
$ grpcurl -plaintext -d '{"id": "00000000-0000-0000-0000-000000000000"}' "0.0.0.0:50051" merchant.MerchantService/GetMerchant | jq .
ERROR:
  Code: Unknown
  Message: id does not match any merchant
```

#### GetQPS

- Existing merchant:

```bash
$ grpcurl -plaintext -d '{"id": "e1211351-bb91-441f-9ea0-3b243189dec6"}' "0.0.0.0:50051" merchant.MerchantService/GetQPS | jq .
{
  "maxQps": 100
}
```

- Non existing merchant

```bash
$ grpcurl -plaintext -d '{"id": "00000000-0000-0000-0000-000000000000"}' "0.0.0.0:50051" merchant.MerchantService/GetQPS | jq .
ERROR:
  Code: Unknown
  Message: id does not match any merchant
```

#### MerchantActive

- Existing merchant:

```bash
$ grpcurl -plaintext -d '{"id": "e1211351-bb91-441f-9ea0-3b243189dec6"}' "0.0.0.0:50051" merchant.MerchantService/MerchantActive | jq .
{
  "active": true
}
```

- Non existing merchant

```bash
$ grpcurl -plaintext -d '{"id": "00000000-0000-0000-0000-000000000000"}' "0.0.0.0:50051" merchant.MerchantService/MerchantActive | jq .
ERROR:
  Code: Unknown
  Message: id does not match any merchant
```

#### FindMerchant

- Correct username and password:

```bash
$ grpcurl -plaintext -d '{"username": "merchant0", "password": "password0"}' "0.0.0.0:50051" merchant.MerchantService/FindMerchant | jq .
{
  "exists": true,
  "id": "63e9f8ed-eb67-4448-8448-6b58ec47248d"
}
```

- Correct username and wrong password:

```bash
$ grpcurl -plaintext -d '{"username": "merchant0", "password": "password0000"}' "0.0.0.0:50051" merchant.MerchantService/FindMerchant | jq .
{}
```

- Invalid username:

```bash
$ grpcurl -plaintext -d '{"username": "merchant10", "password": "password0000"}' "0.0.0.0:50051" merchant.MerchantService/FindMerchant | jq .
{}
```

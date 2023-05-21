# Rate Limiter

It is a rate limiter which can be exposed over gRPC. It has one method:

- `Allow` which receives an id from a **Merchant** and returns true or false depending on the usage.

The usage per **Merchant** is kept in memory, and it requires access to the Merchant Service to learn the MaxQPS per **Merchant**.

## Testing

The rate limiter uses [Golang rate package](https://pkg.go.dev/golang.org/x/time/rate) and runs completely in memory.

If a Merchant Service is listening on localhost, port 50051, then the Rate Limiter can be initialized as follows:

```bash
$ go run main.go --merchant-ip="::1" --merchant-port=50051 --port=50052
```

In a different terminal, please use the [grpcurl](https://github.com/fullstorydev/grpcurl) tool to verify that the endpoints are working as expected.

- **Merchant with MaxQPS>0**

```bash
$ grpcurl -plaintext -d '{"id": "e1211351-bb91-441f-9ea0-3b243189dec6"}' "0.0.0.0:50052" ratelimiter.RateLimiterService/Allow | jq .allow
true
```

- **Merchant with MaxQPS=0**

```bash
$ grpcurl -plaintext -d '{"id": "84a8cb14-0e7b-43f5-8ec7-0840147d3d47"}' "0.0.0.0:50052" ratelimiter.RateLimiterService/Allow | jq .allow
null
```

- **Too many requests**

Merchant `6c1285c2-f09e-4a9b-8a6c-4d94695c1a15` has `MaxQPS=10`:

```bash
$ seq 1 20 | xargs -P $(sysctl hw.logicalcpu | cut -d" " -f2) -I % grpcurl -plaintext -d '{"id": "6c1285c2-f09e-4a9b-8a6c-4d94695c1a15"}' "0.0.0.0:50052" ratelimiter.RateLimiterService/Allow | jq .allow
true
true
true
true
true
true
true
true
true
true
null
null
true
null
null
null
null
null
null
null
```

- **Merchant Service offline - merchant in memory**

```bash
$ grpcurl -plaintext -d '{"id": "e1211351-bb91-441f-9ea0-3b243189dec6"}' "0.0.0.0:50052" ratelimiter.RateLimiterService/Allow | jq .allow
true
```

- **Merchant Service offline - merchant not in memory**

```bash
$ grpcurl -plaintext -d '{"id": "54c5b126-9a3b-4ecc-9590-e8bc5e9bd069"}' "0.0.0.0:50052" ratelimiter.RateLimiterService/Allow | jq .allow
ERROR:
  Code: Unknown
  Message: error reading max qps: rpc error: code = Unavailable desc = connection error: desc = "transport: Error while dialing: dial tcp 0.0.0.0:50051: connect: connection refused"
```
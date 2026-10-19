# Swing Common Libraries

`swing-common` is a Go package with shared helpers for Swing services. The Go
module path is `github.com/getswing/swing-common`, and the package name is
`sw`.

## Install

```sh
go get github.com/getswing/swing-common
```

Import it in your service:

```go
import sw "github.com/getswing/swing-common"
```

The explicit `sw` alias is optional, but is used in the examples below.

## Logging

Initialize the logger once when the service starts, before writing log entries.
It creates `./logs` if necessary and writes JSON logs to stdout and to a
rotating `./logs/<service-name>.log` file.

```go
package main

import (
	"context"

	sw "github.com/getswing/swing-common"
)

func main() {
	sw.LoggerInit("orders-service", nil)

	ctx := context.Background()
	sw.LoggerInfo(ctx, "service started")
	sw.LoggerInfof(ctx, "listening on port %d", 8080)
	sw.LoggerWarn(ctx, "using fallback configuration")
	sw.LoggerErrorf(ctx, "failed to process order %s", "order-123")
}
```

Use `LoggerInfo`, `LoggerWarn`, and `LoggerError` for plain messages. Their
`...f` variants format arguments using `fmt.Sprintf` conventions. Log entries
include the timestamp, service name, level, source file and function; when the
context has a request ID, it is included as `request_id`.

### Request IDs with Echo

Register `RequestIDMiddleware` to accept an incoming `X-Request-ID` header or
generate an ID when the header is missing. The middleware also returns the ID
in the response header. Pass the request context to logger calls to correlate
entries for the same request.

```go
e := echo.New()
e.Use(sw.RequestIDMiddleware)

e.GET("/orders/:id", func(c echo.Context) error {
	ctx := c.Request().Context()
	sw.LoggerInfof(ctx, "fetching order %s", c.Param("id"))
	return c.NoContent(http.StatusOK)
})
```

For non-Echo code, `WithRequestID(ctx)` creates a new request ID, while
`WithRequestIDFromHeader(ctx, header)` uses the given ID or generates one if
the header value is empty. Use `GetRequestID(ctx)` to read the ID.

## PostgreSQL

`ConnectPostgresDB` returns both the GORM connection and its underlying
`*sql.DB`. It configures the connection pool and installs the library's
JSON-formatted GORM logger.

```go
sw.LoggerInit("orders-service", nil)

cfg := sw.CommonDBConfig{
	DBUser:     os.Getenv("DB_USER"),
	DBPassword: os.Getenv("DB_PASSWORD"),
	DBHost:     os.Getenv("DB_HOST"),
	DBPort:     os.Getenv("DB_PORT"),
	DBName:     os.Getenv("DB_NAME"),
	DBSslMode:  "disable",
}

gormDB, sqlDB, err := sw.ConnectPostgresDB(cfg, "orders-service")
if err != nil {
	return err
}
defer sqlDB.Close()

if err := sw.DBPing(context.Background(), sqlDB); err != nil {
	return err
}
```

## RabbitMQ

`Connect` opens a RabbitMQ connection and channel. Close both when finished.
`NewPublisher` and `Consume` provide JSON publishing and queue consumption:

```go
conn, ch, err := sw.Connect(os.Getenv("AMQP_URL"))
if err != nil {
	return err
}
defer conn.Close()
defer ch.Close()

publisher := sw.NewPublisher(ch)
if err := publisher.PublishJSON(ctx, "orders", map[string]string{
	"order_id": "order-123",
}); err != nil {
	return err
}

deliveries, err := sw.Consume(ch, "orders")
if err != nil {
	return err
}
for delivery := range deliveries {
	// Decode and handle delivery.Body, then acknowledge after successful work.
	if err := delivery.Ack(false); err != nil {
		return err
	}
}
```

`Consume` uses manual acknowledgements; acknowledge a message only after it has
been processed successfully. `MQPing(ch)` is available as a channel health
check.

## HTTP requests

`DoRequest` sends a JSON request and returns the response body, HTTP status
code, and any transport or body-reading error:

```go
body, status, err := sw.DoRequest(
	ctx,
	http.MethodGet,
	"https://example.com/api/orders",
	nil,
	map[string]string{"Authorization": "Bearer " + token},
	3, // maximum attempts
)
if err != nil {
	return err
}
if status < http.StatusOK || status >= http.StatusMultipleChoices {
	return fmt.Errorf("orders API returned status %d: %s", status, body)
}
```

The helper sets `Content-Type: application/json`, JSON-encodes a non-nil
request body, and retries transport failures and 5xx responses. A non-2xx HTTP
response is returned with its status code and does not by itself produce an
error.

## Echo responses and validation

Use `SuccessResponse` and `ErrorResponse` to return the package's standard
response shapes from an Echo handler:

```go
return sw.SuccessResponse(c, http.StatusOK, "order found", order)
// or
return sw.ErrorResponse(c, http.StatusNotFound, "order not found")
```

To enable struct-tag validation with Echo, configure the validator once:

```go
e.Validator = &sw.CustomValidator{Validator: validator.New()}
```

Then validate request DTOs using Echo's `c.Validate(dto)`; validation errors
are returned as HTTP 422 responses.

## Other helpers

- `UuidToStringValue`, `StringPointerToStringValue`,
  `BoolPointerToBoolValue`, `IntPointerToInt32Value`, and
  `TimePointerToTimestamp` convert Go values/pointers to protobuf wrapper or
  timestamp values. Nil pointers (and zero times for `TimePointerToTimestamp`)
  return nil.
- `PathToS3File(path)` prefixes a path with the library's `__s3Url__` marker.
- `WithContext(ctx)` is Echo middleware that replaces each request's context
  with the supplied context.

There is currently no ETS API in this Go module. If “ETS” refers to a specific
feature you expected here, it is not exposed by the current library.

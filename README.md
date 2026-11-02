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

### Request IDs over gRPC

The request ID travels as the `x-request-id` gRPC metadata key
(`GRPCRequestIDKey`). Register the client interceptors so every outgoing call
forwards the ID from `ctx`:

```go
conn, err := grpc.NewClient(
	addr,
	grpc.WithTransportCredentials(insecure.NewCredentials()),
	grpc.WithUnaryInterceptor(sw.RequestIDUnaryClientInterceptor()),
	grpc.WithStreamInterceptor(sw.RequestIDStreamClientInterceptor()),
)
```

Register the server interceptors so handlers get the caller's ID in `ctx`, or
a generated one when the caller didn't send it. The ID is also returned in the
response header metadata.

```go
grpcServer := grpc.NewServer(
	grpc.ChainUnaryInterceptor(sw.RequestIDUnaryServerInterceptor()),
	grpc.ChainStreamInterceptor(sw.RequestIDStreamServerInterceptor()),
)
```

For a single call without an interceptor, use
`sw.OutgoingRequestIDContext(ctx)` on the client and
`sw.IncomingRequestIDContext(ctx)` on the server.
## Tracing (OpenTelemetry)

Tracing shows where the time of a single request goes. When an API is slow,
open its trace in your tracing backend (Jaeger, Tempo, Grafana, etc.): every
step is a span on a timeline, so the longest span is the bottleneck — a slow
SQL query, a slow downstream API, or slow business logic.

The library sends spans over OTLP/HTTP to an OpenTelemetry collector (or any
backend that accepts OTLP, such as Jaeger on port `4318`).

### 1. Initialize once at startup

```go
func main() {
	ctx := context.Background()
	sw.LoggerInit("player-service", nil)

	shutdownTracing, err := sw.TracingInit(ctx, sw.TracingConfig{
		ServiceName:    "player-service",
		ServiceVersion: os.Getenv("APP_VERSION"),
		Environment:    os.Getenv("APP_ENV"), // e.g. "staging", "production"
		Endpoint:       os.Getenv("OTEL_COLLECTOR_HOST"), // e.g. "localhost:4318"
		Insecure:       true,  // plain HTTP to the collector
		SampleRatio:    1,     // record every trace; lower it (e.g. 0.1) under high traffic
		Disabled:       os.Getenv("TRACING_DISABLED") == "true",
	})
	if err != nil {
		sw.LoggerErrorf(ctx, "[Telemetry] failed to initialize tracing: %v", err)
		panic(err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = shutdownTracing(ctx) // flushes buffered spans
	}()
	// ...
}
```

`TracingConfig` fields:

| Field | Meaning |
| --- | --- |
| `ServiceName` | Name shown in the tracing UI. Defaults to the name given to `LoggerInit`. |
| `ServiceVersion`, `Environment` | Optional; useful to filter traces by deploy or environment. |
| `Endpoint` | Collector `host:port`. If empty, the standard `OTEL_EXPORTER_OTLP_ENDPOINT` / `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT` env vars are used. |
| `Insecure` | Use plain HTTP instead of HTTPS to the collector. |
| `SampleRatio` | Fraction of new traces to keep (`0 < x <= 1`). `0` means `1`. Incoming requests that are already sampled by the caller are always kept. |
| `Disabled` | Do not export anything; all tracing calls become cheap no-ops. |

If `TracingInit` is never called, every helper below is a no-op, so the code
is safe to leave in place in tests and scripts.

### 2. Trace every HTTP request (Echo)

Register `TracingMiddleware` **before** `RequestIDMiddleware`:

```go
e := echo.New()
e.Use(sw.TracingMiddleware)
e.Use(sw.RequestIDMiddleware)
```

Each request gets a root span named `<METHOD> <route>` (for example
`GET /api/v1/players/:id`) with the method, route, path, client IP, status code
and `request_id`. 5xx responses and returned errors mark the span as failed.
The trace ID is returned to the caller in the `X-Trace-ID` response header, so
a slow or failing request reported by the client/QA can be looked up directly.

If the caller sent a W3C `traceparent` header, the span joins the caller's
trace instead of starting a new one.

### 3. Add spans in your own code

The request span alone only tells you that the API was slow. Add child spans
around the steps you want to measure (usecase, repository, external calls).
Always pass the `ctx` returned by the span into the next layer, otherwise the
child spans will not be connected.

```go
// Handler
func (h *PlayerHandler) GetPlayer(c echo.Context) error {
	ctx := c.Request().Context() // already contains the request span
	player, err := h.usecase.GetPlayer(ctx, c.Param("id"))
	// ...
}

// Usecase
func (u *playerUsecase) GetPlayer(ctx context.Context, id string) (_ *Player, err error) {
	ctx, span := sw.StartSpan(ctx, "player.usecase.get_player",
		attribute.String("player.id", id),
	)
	defer func() { sw.EndSpan(span, err) }()

	player, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	// ...
}
```

For a small block, `TraceFunc` starts and ends the span for you:

```go
err := sw.TraceFunc(ctx, "player.calculate_points", func(ctx context.Context) error {
	return u.calculatePoints(ctx, player)
})
```

Available helpers:

| Function | Use |
| --- | --- |
| `StartSpan(ctx, name, attrs...)` | Start a child span; returns the new `ctx` and the span. |
| `EndSpan(span, err)` | End the span; when `err != nil` it is recorded and the span is marked as an error. |
| `TraceFunc(ctx, name, fn, attrs...)` | Run `fn` inside a span. |
| `GetTraceID(ctx)`, `GetSpanID(ctx)` | Read the current IDs (e.g. to return them in an error response). |
| `InjectTraceHeaders(ctx, header)` | Put the trace context into outgoing `http.Header`. |
| `ExtractTraceHeaders(ctx, header)` | Continue a trace from incoming `http.Header` (non-Echo servers). |

Attributes use `go.opentelemetry.io/otel/attribute`
(`attribute.String`, `attribute.Int`, `attribute.Bool`, ...).

### What is traced automatically

- **Database queries**: with `ConnectPostgresDB`, every GORM query run with a
  traced context becomes a `db.query` span containing the SQL text, row count
  and exact duration. Use `gormDB.WithContext(ctx)` in repositories, otherwise
  the query cannot be linked to the request.
- **Outgoing HTTP**: `DoRequest` creates an `HTTP <METHOD>` client span (covering
  all retries) and sends the `traceparent` header, so if the downstream service
  also uses this library both services appear in one trace.
- **Logs**: `LoggerInfo`/`LoggerError`/... and the GORM logger add `trace_id`
  and `span_id` to each JSON entry when the context has a span, so you can jump
  from a trace to its logs and back.

gRPC and RabbitMQ are not instrumented by this library. For gRPC, the global
provider and propagator set by `TracingInit` are picked up by
`go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc`
(`grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))` and
`grpc.WithStatsHandler(otelgrpc.NewClientHandler())` on clients).

### Finding a bottleneck

1. Get the trace ID from the `X-Trace-ID` response header or the `trace_id`
   field in the logs, or search the tracing UI by service + route and sort by
   duration.
2. Open the trace and look at the waterfall. Find the span that takes most of
   the total time and has no long child spans — that is where the time is
   actually spent.
3. Typical patterns:
   - one long `db.query` → missing index or heavy query; check its SQL text;
   - many short `db.query` spans in a row → N+1 queries, batch or preload them;
   - long `HTTP GET` span → slow downstream service (open its part of the trace);
   - gaps with no child span → untraced code; add `StartSpan` there.

### Conventions

- Span names: `<layer-or-domain>.<action>` in lower snake case, e.g.
  `player.usecase.get_player`, `player.repo.find_by_id`, `payment.client.charge`.
  Keep names low-cardinality — put IDs in attributes, not in the name.
- Do not put secrets or personal data (passwords, OTPs, tokens, phone numbers)
  in attributes. Note that `db.query` spans contain the SQL text as logged by
  GORM.
- Always end spans (`defer`), and always pass the span's `ctx` downstream.
- For background goroutines that should outlive the request, start a new span
  from `context.WithoutCancel(ctx)` so the work stays in the same trace without
  being cancelled with the request.

### Local collector for development

```sh
docker run --rm -p 16686:16686 -p 4318:4318 jaegertracing/all-in-one:latest
```

Set `Endpoint: "localhost:4318"` and `Insecure: true`, then open
http://localhost:16686 to browse traces.

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

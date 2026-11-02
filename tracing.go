package sw

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

const tracerName = "github.com/getswing/swing-common"

type TracingConfig struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	// Endpoint is the OTLP/HTTP collector host:port, e.g. "localhost:4318".
	// When empty, the exporter reads OTEL_EXPORTER_OTLP_(TRACES_)ENDPOINT.
	Endpoint string
	Insecure bool
	// SampleRatio is the fraction of new traces to record (0 < x <= 1).
	// Zero or out of range defaults to 1 (record everything).
	SampleRatio float64
	// Disabled installs no exporter; spans become no-ops.
	Disabled bool
}

// TracingInit installs a global tracer provider and W3C trace-context
// propagator. Call the returned shutdown func before the service exits so
// buffered spans are flushed.
func TracingInit(ctx context.Context, cfg TracingConfig) (func(context.Context) error, error) {
	if cfg.ServiceName == "" {
		cfg.ServiceName = ServiceName
	}

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	if cfg.Disabled {
		return func(context.Context) error { return nil }, nil
	}

	var opts []otlptracehttp.Option
	if cfg.Endpoint != "" {
		opts = append(opts, otlptracehttp.WithEndpoint(cfg.Endpoint))
	}
	if cfg.Insecure {
		opts = append(opts, otlptracehttp.WithInsecure())
	}
	exporter, err := otlptracehttp.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("create otlp trace exporter: %w", err)
	}

	attrs := []attribute.KeyValue{attribute.String("service.name", cfg.ServiceName)}
	if cfg.ServiceVersion != "" {
		attrs = append(attrs, attribute.String("service.version", cfg.ServiceVersion))
	}
	if cfg.Environment != "" {
		attrs = append(attrs, attribute.String("deployment.environment.name", cfg.Environment))
	}

	ratio := cfg.SampleRatio
	if ratio <= 0 || ratio > 1 {
		ratio = 1
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewSchemaless(attrs...)),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(ratio))),
	)
	otel.SetTracerProvider(provider)

	return provider.Shutdown, nil
}

// StartSpan starts a child span of whatever span is in ctx. Always call
// EndSpan (or span.End) on the returned span.
func StartSpan(ctx context.Context, name string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	return otel.Tracer(tracerName).Start(ctx, name, trace.WithAttributes(attrs...))
}

// EndSpan records err (if any) on the span, marks it as failed and ends it.
func EndSpan(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	span.End()
}

// TraceFunc wraps fn in a span named name.
func TraceFunc(ctx context.Context, name string, fn func(ctx context.Context) error, attrs ...attribute.KeyValue) error {
	ctx, span := StartSpan(ctx, name, attrs...)
	err := fn(ctx)
	EndSpan(span, err)
	return err
}

func GetTraceID(ctx context.Context) string {
	sc := trace.SpanContextFromContext(ctx)
	if !sc.HasTraceID() {
		return ""
	}
	return sc.TraceID().String()
}

func GetSpanID(ctx context.Context) string {
	sc := trace.SpanContextFromContext(ctx)
	if !sc.HasSpanID() {
		return ""
	}
	return sc.SpanID().String()
}

// InjectTraceHeaders writes the trace context of ctx into outgoing headers so
// the downstream service continues the same trace.
func InjectTraceHeaders(ctx context.Context, header http.Header) {
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(header))
}

// ExtractTraceHeaders returns ctx carrying the trace context from incoming headers.
func ExtractTraceHeaders(ctx context.Context, header http.Header) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(header))
}

// TracingMiddleware creates a server span per Echo request, named
// "<METHOD> <route>". Register it before RequestIDMiddleware so log entries
// written by handlers carry the trace ID.
func TracingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		route := c.Path()
		if route == "" {
			route = req.URL.Path
		}

		ctx := ExtractTraceHeaders(req.Context(), req.Header)
		ctx, span := otel.Tracer(tracerName).Start(ctx, req.Method+" "+route,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("http.request.method", req.Method),
				attribute.String("http.route", route),
				attribute.String("url.path", req.URL.Path),
				attribute.String("client.address", c.RealIP()),
				attribute.String("user_agent.original", req.UserAgent()),
			),
		)
		defer span.End()

		if reqID := req.Header.Get("X-Request-ID"); reqID != "" {
			span.SetAttributes(attribute.String("request_id", reqID))
		}
		if traceID := GetTraceID(ctx); traceID != "" {
			c.Response().Header().Set("X-Trace-ID", traceID)
		}

		c.SetRequest(req.WithContext(ctx))

		err := next(c)
		if err != nil {
			span.RecordError(err)
			// Let Echo write the error response so the final status is known.
			c.Error(err)
		}

		status := c.Response().Status
		span.SetAttributes(attribute.Int("http.response.status_code", status))
		if status >= http.StatusInternalServerError {
			span.SetStatus(codes.Error, http.StatusText(status))
		}
		return nil
	}
}

// traceQuery records an already finished GORM query as a span, so slow
// queries show up inside the request trace.
func traceQuery(ctx context.Context, dbSystem string, begin time.Time, elapsed time.Duration, sql string, rows int64, err error) {
	if !trace.SpanContextFromContext(ctx).IsValid() {
		return
	}

	attrs := []attribute.KeyValue{
		attribute.String("db.query.text", sanitizeSQL(sql)),
		attribute.Int64("db.response.returned_rows", rows),
	}
	if dbSystem != "" {
		attrs = append(attrs, attribute.String("db.system.name", dbSystem))
	}

	_, span := otel.Tracer(tracerName).Start(ctx, "db.query",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithTimestamp(begin),
		trace.WithAttributes(attrs...),
	)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	span.End(trace.WithTimestamp(begin.Add(elapsed)))
}

// dbSystemName maps a GORM dialector name to the OpenTelemetry
// db.system.name value; unknown names are passed through unchanged.
func dbSystemName(dialector string) string {
	switch dialector {
	case "postgres":
		return "postgresql"
	case "sqlserver":
		return "microsoft.sql_server"
	default:
		return dialector
	}
}

var (
	sqlStringLiteral  = regexp.MustCompile(`'(?:[^']|'')*'`)
	sqlNumericLiteral = regexp.MustCompile(`\b\d+(?:\.\d+)?\b`)
)

// sanitizeSQL replaces literal values with "?" because GORM hands the logger
// SQL with the bound values (phone numbers, tokens, ...) already inlined.
func sanitizeSQL(sql string) string {
	sql = sqlStringLiteral.ReplaceAllString(sql, "?")
	return sqlNumericLiteral.ReplaceAllString(sql, "?")
}

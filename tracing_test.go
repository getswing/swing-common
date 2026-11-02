package sw

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestTracingExportsRequestSpans(t *testing.T) {
	exported := make(chan string, 4)
	collector := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		exported <- string(body)
	}))
	defer collector.Close()

	var downstreamTraceparent string
	downstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		downstreamTraceparent = r.Header.Get("traceparent")
		w.WriteHeader(http.StatusOK)
	}))
	defer downstream.Close()

	t.Setenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", collector.URL+"/v1/traces")
	shutdown, err := TracingInit(context.Background(), TracingConfig{ServiceName: "test-service"})
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	e.Use(TracingMiddleware)
	e.GET("/players/:id", func(c echo.Context) error {
		ctx := c.Request().Context()
		if err := TraceFunc(ctx, "usecase.get_player", func(ctx context.Context) error {
			_, _, err := DoRequest(ctx, http.MethodGet, downstream.URL, nil, nil, 1)
			return err
		}); err != nil {
			return err
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "boom")
	})

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/players/1", nil))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d", rec.Code)
	}
	traceID := rec.Header().Get("X-Trace-ID")
	if traceID == "" {
		t.Fatal("missing X-Trace-ID header")
	}
	if !strings.Contains(downstreamTraceparent, traceID) {
		t.Fatalf("downstream traceparent %q does not carry trace %s", downstreamTraceparent, traceID)
	}

	if err := shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
	body := <-exported
	for _, name := range []string{"GET /players/:id", "usecase.get_player", "HTTP GET"} {
		if !strings.Contains(body, name) {
			t.Errorf("exported spans missing %q", name)
		}
	}
}

func TestEndSpanWithoutProviderIsNoop(t *testing.T) {
	_, span := StartSpan(context.Background(), "noop")
	EndSpan(span, errors.New("ignored"))
	if GetTraceID(context.Background()) != "" {
		t.Fatal("expected empty trace id")
	}
}

func TestSanitizeSQL(t *testing.T) {
	got := sanitizeSQL(`SELECT * FROM "otps" WHERE phone_number = '0812''34' AND table2.id = 15 LIMIT 1`)
	want := `SELECT * FROM "otps" WHERE phone_number = ? AND table2.id = ? LIMIT ?`
	if got != want {
		t.Fatalf("got  %s\nwant %s", got, want)
	}
}

func TestDBSystemName(t *testing.T) {
	for in, want := range map[string]string{"postgres": "postgresql", "mysql": "mysql", "sqlserver": "microsoft.sql_server"} {
		if got := dbSystemName(in); got != want {
			t.Errorf("dbSystemName(%q) = %q, want %q", in, got, want)
		}
	}
}

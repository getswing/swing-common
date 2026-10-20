package sw

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestRequestIDUnaryClientInterceptor(t *testing.T) {
	ctx := WithRequestIDFromHeader(context.Background(), "req-123")

	var got []string
	invoker := func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		md, _ := metadata.FromOutgoingContext(ctx)
		got = md.Get(GRPCRequestIDKey)
		return nil
	}

	if err := RequestIDUnaryClientInterceptor()(ctx, "/svc/M", nil, nil, nil, invoker); err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != "req-123" {
		t.Fatalf("outgoing x-request-id = %v, want [req-123]", got)
	}
}

func TestRequestIDUnaryServerInterceptor(t *testing.T) {
	handler := func(ctx context.Context, req any) (any, error) {
		return GetRequestID(ctx), nil
	}

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(GRPCRequestIDKey, "req-456"))
	got, _ := RequestIDUnaryServerInterceptor()(ctx, nil, &grpc.UnaryServerInfo{}, handler)
	if got != "req-456" {
		t.Fatalf("request id = %v, want req-456", got)
	}

	got, _ = RequestIDUnaryServerInterceptor()(context.Background(), nil, &grpc.UnaryServerInfo{}, handler)
	if got == "" {
		t.Fatal("expected a generated request id when metadata is missing")
	}
}

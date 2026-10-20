package sw

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// GRPCRequestIDKey is the gRPC metadata key carrying the request ID.
// gRPC metadata keys are lowercase; this is the equivalent of the X-Request-ID header.
const GRPCRequestIDKey = "x-request-id"

// OutgoingRequestIDContext attaches the request ID in ctx to the outgoing gRPC metadata.
func OutgoingRequestIDContext(ctx context.Context) context.Context {
	id := GetRequestID(ctx)
	if id == "" {
		return ctx
	}
	if md, ok := metadata.FromOutgoingContext(ctx); ok && len(md.Get(GRPCRequestIDKey)) > 0 {
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, GRPCRequestIDKey, id)
}

// IncomingRequestIDContext stores the request ID from incoming gRPC metadata in ctx,
// generating a new one when the caller didn't send it.
func IncomingRequestIDContext(ctx context.Context) context.Context {
	var id string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if v := md.Get(GRPCRequestIDKey); len(v) > 0 {
			id = v[0]
		}
	}
	return WithRequestIDFromHeader(ctx, id)
}

// RequestIDUnaryClientInterceptor forwards the request ID on every unary call.
// Usage: grpc.NewClient(addr, grpc.WithUnaryInterceptor(sw.RequestIDUnaryClientInterceptor()))
func RequestIDUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return invoker(OutgoingRequestIDContext(ctx), method, req, reply, cc, opts...)
	}
}

// RequestIDStreamClientInterceptor forwards the request ID on every streaming call.
func RequestIDStreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return streamer(OutgoingRequestIDContext(ctx), desc, cc, method, opts...)
	}
}

// RequestIDUnaryServerInterceptor puts the caller's request ID into the handler context
// so sw.Logger* calls log the same ID as the caller.
// Usage: grpc.NewServer(grpc.UnaryInterceptor(sw.RequestIDUnaryServerInterceptor()))
func RequestIDUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		ctx = IncomingRequestIDContext(ctx)
		_ = grpc.SetHeader(ctx, metadata.Pairs(GRPCRequestIDKey, GetRequestID(ctx)))
		return handler(ctx, req)
	}
}

// RequestIDStreamServerInterceptor is the streaming variant of RequestIDUnaryServerInterceptor.
func RequestIDStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := IncomingRequestIDContext(ss.Context())
		_ = ss.SetHeader(metadata.Pairs(GRPCRequestIDKey, GetRequestID(ctx)))
		return handler(srv, &requestIDServerStream{ServerStream: ss, ctx: ctx})
	}
}

type requestIDServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *requestIDServerStream) Context() context.Context { return s.ctx }

package client

import (
	"google.golang.org/grpc"
	// pb "getswing.app/player-service/"
)

type GrpcClient struct {
	conn *grpc.ClientConn
	// client pb.GreeterClient
}

// NewGrpcClient membuat koneksi ke gRPC server lain
func NewGrpcClient(target string) (*GrpcClient, error) {
	// conn, err := grpc.Dial(
	// 	target,
	// 	grpc.WithTransportCredentials(insecure.NewCredentials()),
	// 	grpc.WithBlock(),
	// )
	// if err != nil {
	// 	return nil, err
	// }

	// c := pb.NewGreeterClient(conn)
	// return &GrpcClient{conn: conn, client: c}, nil

	return nil, nil
}

// Close menutup koneksi
func (g *GrpcClient) Close() {
	g.conn.Close()
}

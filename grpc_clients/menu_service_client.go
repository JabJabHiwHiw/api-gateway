package grpc_clients

import (
    pb "food-service/proto"
    "google.golang.org/grpc"
)

func NewMenuServiceClient(addr string) (pb.MenuServiceClient, *grpc.ClientConn, error) {
    conn, err := grpc.Dial(addr, grpc.WithInsecure())
    if err != nil {
        return nil, nil, err
    }
    client := pb.NewMenuServiceClient(conn)
    return client, conn, nil
}
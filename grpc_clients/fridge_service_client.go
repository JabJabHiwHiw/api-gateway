package grpc_clients

import (
    pb "food-service/proto"
    "google.golang.org/grpc"
)

func NewFridgeItemServiceClient(addr string) (pb.FridgeItemServiceClient, *grpc.ClientConn, error) {
    conn, err := grpc.Dial(addr, grpc.WithInsecure())
    if err != nil {
        return nil, nil, err
    }
    client := pb.NewFridgeItemServiceClient(conn)
    return client, conn, nil
}
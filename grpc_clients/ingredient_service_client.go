package grpc_clients

import (
    pb "food-service/proto"
    "google.golang.org/grpc"
)

func NewIngredientServiceClient(addr string) (pb.IngredientServiceClient, *grpc.ClientConn, error) {
    conn, err := grpc.Dial(addr, grpc.WithInsecure())
    if err != nil {
        return nil, nil, err
    }
    client := pb.NewIngredientServiceClient(conn)
    return client, conn, nil
}
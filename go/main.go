package main

import (
	"context"
	"fmt"
	"net"

	"grpc-server/pb"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedUserServiceServer
}

func (s *server) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	fmt.Printf("Menerima request untuk User ID: %d\n", req.GetId())
	return &pb.UserResponse{
		Id:    req.GetId(),
		Name:  "Shafa Nabila",
		Phone: "081234567890",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Printf("Gagal mendengarkan port: %v\n", err)
		return
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{})

	fmt.Printf("gRPC Server berjalan di port :50051...\n")
	if err := s.Serve(lis); err != nil {
		fmt.Printf("Gagal menjalankan server: %v\n", err)
	}
}
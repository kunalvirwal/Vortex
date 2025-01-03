package master

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	pb "github.com/kunalvirwal/Vortex/proto/factory"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedContainerFactoryServer
}

// func (s *server) Apply(ctx context.Context, body *pb.RequestBody) (*pb.BoolResponse, error) {
// 	data := body.GetData()

// }

func (s *server) Down(ctx context.Context, body *pb.NameHolder) (*pb.NameHolder, error) {
	msg := body.GetName()
	if msg != "kill" {
		return nil, errors.New("invalid command recieved")
	}
	go func() {
		timer := time.NewTimer(3 * time.Second)
		<-timer.C
		os.Exit(0)
	}()
	return &pb.NameHolder{Name: "Stopping Vortex...Bye"}, nil
}

func StartGrpcServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Println("failed to listen:", err)
		os.Exit(1)
	}
	grpcServer := grpc.NewServer(
	// Can define grpc interceptors(middlewares) here
	// grpc.StreamInterceptor(func{}),
	// grpc.UnaryInterceptor(func{}),
	)
	pb.RegisterContainerFactoryServer(grpcServer, &server{})
	grpcServer.Serve(lis)
}

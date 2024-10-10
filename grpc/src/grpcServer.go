package src

import (
	pb "OCluster_runner/grpc/src/orunner"
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
	"net"
	"os"
)

type server struct {
	pb.UnimplementedOClusterServer
}

func (s *server) Health(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{Status: true}, nil
}

func (s server) FileUploader(stream grpc.BidiStreamingServer[pb.FileChunkRequest, pb.FileChunkResponse]) error {
	randomUUID, err := uuid.NewUUID()
	file, err := os.OpenFile(randomUUID.String(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	var file_name string

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			file.Close()
			err = os.Rename(randomUUID.String(), file_name)
			return nil
		}
		file_name = req.FileName
		if err != nil {
			return err
		}
		_, err = file.Write(req.Chunk)
		fmt.Println("chunk saved")
		res := pb.FileChunkResponse{Percent: 1, Status: true}
		stream.Send(&res)
	}
}

func RunGRPC() {
	listen, err := net.Listen("tcp", ":50052")
	if err != nil {
		fmt.Println(err)
	}

	s := grpc.NewServer()
	pb.RegisterOClusterServer(s, &server{})
	reflection.Register(s)
	fmt.Println("listening at", listen.Addr())

	err = s.Serve(listen)
	if err != nil {
		fmt.Println(err)
	}
}

package main

import (
	"flag"
	//	"fmt"
	"os"
	"strconv"

	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	//	"google.golang.org/grpc/credentials"
	//	"google.golang.org/grpc/grpclog"

	//	"github.com/golang/protobuf/proto"

	"google.golang.org/grpc/reflection"
	pb "sohu.com/grpcTransfer/transfer"
)

var (
	port = flag.Int("port", 10000, "The server port")
)

// server is used to implement server interface.
type server struct{}

// implement upload function
func (s *server) Upload(ctx context.Context, upRequest *pb.UploadRequest) (*pb.UploadReply, error) {
	f, err := os.Create(upRequest.FileName)
	if err != nil {
		return &pb.UploadReply{Status: false}, err //错误情况下uploadReplay为false
	}
	defer f.Close()
	f.Write([]byte(upRequest.FileContents))
	return &pb.UploadReply{Status: true}, nil //正确情况下uploadReply为ture值
}

// implement download function
func (s *server) Download(ctx context.Context, downRequest *pb.DownloadRequest) (*pb.DownloadReply, error) {
	//to be done
	f, err := os.Open(downRequest.GetFileName())
	defer f.Close()
	if err != nil {
		return &pb.DownloadReply{FileExistence: false, FileContents: err.Error()}, nil
	}
	fileStatus, _ := f.Stat()
	if fileStatus.IsDir() {
		return &pb.DownloadReply{FileExistence: false, FileContents: "failed: you should input a file, not a folder"}, nil
	}
	data := make([]byte, 1024)
	count, err := f.Read(data)
	if err != nil {
		return &pb.DownloadReply{FileExistence: false, FileContents: err.Error()}, nil
	}
	//	fmt.Printf("read %d bytes: %q\n", count, data[:count])
	fileStr := string(data[:count])
	return &pb.DownloadReply{FileExistence: true, FileContents: fileStr}, nil
}
func main() {
	flag.Parse()
	log.Print(*port)
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterFileTransferServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

package main

import (
	"flag"
	//	"fmt"
	"log"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "sohu.com/grpcTransfer/transfer"
)

var (
	uploadName   = flag.String("upload", "", "the name of the file you want to upload")
	downloadName = flag.String("download", "", "the name of the file you want to download")
	fileContent  = flag.String("content", "", "the content of the file you want to upload")
	addressNport = flag.String("P", "127.0.0.1:8981", "The server IPAddress and port")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addressNport, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewFileTransferClient(conn)

	// 设置判断参数是否为空的布尔值
	uploadNameEmpty := (*uploadName == "")
	downloadNameEmpty := (*downloadName == "")
	fileContentEmpty := (*fileContent == "")

	// 判断参数结构分别进行上传、下载和报错
	if !uploadNameEmpty && downloadNameEmpty && !fileContentEmpty {
		r, err := c.Upload(context.Background(), &pb.UploadRequest{FileContents: *fileContent, FileName: *uploadName})
		if err != nil {
			log.Fatalf("fail to upload: %v", err)
		}
		if r.GetStatus() == true {
			log.Println("upload success!")
		}

	} else if uploadNameEmpty && !downloadNameEmpty && fileContentEmpty {
		r, err := c.Download(context.Background(), &pb.DownloadRequest{FileName: *downloadName})
		if err != nil {
			log.Fatalf("fail to download: %v", err)
		}
		//文件不存在或为文件夹时
		if !r.FileExistence {
			log.Fatalf("fail to download: %v", r.FileContents)
		}
		//服务器上存在时，写入本地
		f, err := os.Create(*downloadName)
		if err != nil {
			log.Fatalf("failed to create file:%v", err)
		}
		defer f.Close()
		f.Write([]byte(r.FileContents))
	} else {
		//命令行参数错误
		log.Fatalln("Operation failed: Arguments in wrong format\nYou should either include\n\n uploadName and fileContent\n\n or \n\n downloadName ")
	}

}

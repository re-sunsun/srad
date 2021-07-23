package main

import (
	pb "SunsunSRAD/ssrad/rpc"
	"context"
	"google.golang.org/grpc"
	"log"
)

//注册测试
func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}

	gc := pb.NewRDRpcClient(conn)

	request := pb.RegisterRequest{
		Servername:  "用户注册服务",
		Serverstate: "127.0.0.1:1234561",
	}

	response, err := gc.Register(context.Background(), &request)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(response)
}

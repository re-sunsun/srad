package main

import (
	pb "SunsunSRAD/ssrad/rpc"
	"context"
	"google.golang.org/grpc"
	"log"
)

//发现测试
func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}

	gc := pb.NewRDRpcClient(conn)

	request := pb.DiscoverRequest{
		Servername: "用户注册服务",
	}

	response, err := gc.Discover(context.Background(), &request)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(response)

}

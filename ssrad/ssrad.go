package ssrad

import (
	pb "SunsunSRAD/ssrad/rpc"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type ServersPool struct {
	Servers []Server
}
type Server struct {
	ServerName  string
	ServerState []string
}

func Default() *ServersPool {
	sp := &ServersPool{
		Servers: make([]Server, 0),
	}
	return sp
}

var Servers []Server

type RpcServer struct{}

//注册服务
func (r *RpcServer) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	fmt.Println("接收到服务注册请求\n服务名:", request.Servername, "\n服务信息", request.Serverstate)
	//遍历查找是否存在同名服务
	for _, v := range Servers {
		if v.ServerName == request.Servername {
			v.ServerState = append(v.ServerState, request.Serverstate)

			response := pb.RegisterResponse{
				Msg: "已有同名服务,注册成功",
			}
			return &response, nil
		}
	}
	//未找到同名服务,则添加同名服务
	//新建服务
	server := Server{
		ServerName:  request.Servername,
		ServerState: make([]string, 0),
	}
	//写入服务信息
	server.ServerState = append(server.ServerState, request.Serverstate)
	Servers = append(Servers, server)
	response := pb.RegisterResponse{
		Msg: "未有同名服务,新建并注册成功",
	}
	return &response, nil
}

//发现服务
func (r *RpcServer) Discover(ctx context.Context, request *pb.DiscoverRequest) (*pb.DiscoverResponse, error) {
	fmt.Println("接收到服务发现请求\n服务名:", request.Servername)
	//遍历查找服务
	for _, v := range Servers {
		if v.ServerName == request.Servername {
			response := pb.DiscoverResponse{
				Msg:         "找到该服务了，给你一个",
				Serverstate: v.ServerState[len(v.ServerState)-1],
			}
			return &response, nil
		}
	}
	response := pb.DiscoverResponse{
		Msg:         "没有找到该服务",
		Serverstate: "",
	}

	return &response, nil
}

//开始监听
func (sp *ServersPool) Run(Port string) {
	sp.Servers = make([]Server, 0)
	Servers = sp.Servers
	listen, err := net.Listen("tcp", Port)
	if err != nil {
		log.Fatalln(err)
	}
	gs := grpc.NewServer()
	pb.RegisterRDRpcServer(gs, &RpcServer{})
	reflection.Register(gs)
	err = gs.Serve(listen)
	if err != nil {
		log.Fatalln(err)
	}
}

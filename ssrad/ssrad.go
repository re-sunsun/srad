package ssrad

import (
	"SunsunSRAD/ssrad/cache"
	pb "SunsunSRAD/ssrad/rpc"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type ServicesPool struct {
}

func Default() *ServicesPool {
	sp := &ServicesPool{}
	return sp
}

type RpcServer struct{}

//注册服务
func (r *RpcServer) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	fmt.Println("接收到服务注册请求\n服务名称:", request.ServiceName, "\n服务提供者id", request.ServiceId, "\n服务地址", request.ServiceAddress)

	//本来想利用redis的过期时间与hash来实现一个实时的服务信息表，即f-k-v-ex这样的格式
	//但是redis自带的过期时间只支持顶级key
	//所以还是采用了k-v的形式,服务名.提供者id-服务地址
	conn := cache.Pool.Get()
	defer conn.Close()

	conn.Do("SETEX", request.ServiceName+"."+request.ServiceId, 15, request.ServiceAddress)

	response := pb.RegisterResponse{
		Msg: "注册成功",
	}

	return &response, nil
	////遍历查找是否存在同名服务
	//for _, v := range Servers {
	//	if v.ServerName == request.Servername {
	//		v.ServerState = append(v.ServerState, request.Serverstate)
	//
	//		response := pb.RegisterResponse{
	//			Msg: "已有同名服务,注册成功",
	//		}
	//		return &response, nil
	//	}
	//}
	////未找到同名服务,则添加同名服务
	////新建服务
	//server := Server{
	//	ServerName:  request.Servername,
	//	ServerState: make([]string, 0),
	//}
	////写入服务信息
	//server.ServerState = append(server.ServerState, request.Serverstate)
	//Servers = append(Servers, server)
	//response := pb.RegisterResponse{
	//	Msg: "未有同名服务,新建并注册成功",
	//}
}

//发现服务
func (r *RpcServer) Discover(ctx context.Context, request *pb.DiscoverRequest) (*pb.DiscoverResponse, error) {
	//fmt.Println("接收到服务发现请求\n服务名:", request.Servername)
	////遍历查找服务
	//for _, v := range Servers {
	//	if v.ServerName == request.Servername {
	//		//该部分将在后续实现负载均衡时改进
	//		response := pb.DiscoverResponse{
	//			Msg:         "找到该服务了，给你一个",
	//			Serverstate: v.ServerState[len(v.ServerState)-1],
	//		}
	//		return &response, nil
	//	}
	//}
	//response := pb.DiscoverResponse{
	//	Msg:         "没有找到该服务",
	//	Serverstate: "",
	//}
	return nil, nil
}

func (r *RpcServer) HeartbeatDetection(ctx context.Context, request *pb.HeartbeatDetectionRequest) (*pb.HeartbeatDetectionResponse, error) {
	fmt.Println("更新ttl" + request.ServiceName + request.ServiceId)

	conn := cache.Pool.Get()
	defer conn.Close()

	conn.Do("EXPIRE", request.ServiceName+"."+request.ServiceId, 15)

	response := pb.HeartbeatDetectionResponse{
		Msg: "ttl更新成功",
	}

	return &response, nil
}

//开始监听
func (sp *ServicesPool) Run(Port string) {

	//预载redis
	cache.RedisInit()

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

package ssradc

import (
	pb "SunsunSRAD/ssradc/rpc"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"time"
)

//ssradc - ssrad clint
type SSRADClint struct {
	conn *grpc.ClientConn
}

var exit chan string

func Default(IpAndPort string) (*SSRADClint, error) {
	sc := &SSRADClint{}

	conn, err := grpc.Dial(IpAndPort, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	sc.conn = conn

	return sc, nil
}

func (sc *SSRADClint) Register(ServiceName string, ServiceId string, ServiceAddress string) error {
	//录入信息
	rpcClient := pb.NewRDRpcClient(sc.conn)
	request := pb.RegisterRequest{
		ServiceName:    ServiceName,
		ServiceId:      ServiceId,
		ServiceAddress: ServiceAddress,
	}
	_, err := rpcClient.Register(context.Background(), &request)
	if err != nil {
		return err
	}

	exit = make(chan string)

	go sc.HeartbeatDetection(ServiceName, ServiceId)

	return nil
}

func (sc *SSRADClint) Discover() {
}

func (sc *SSRADClint) HeartbeatDetection(ServiceName string, ServiceId string) {
	ticker := time.NewTicker(time.Duration(5) * time.Second)
	for {
		select {
		case time := <-ticker.C:
			fmt.Println("心跳保持:", time)
			//更新信息
			rpcClient := pb.NewRDRpcClient(sc.conn)
			request := pb.HeartbeatDetectionRequest{
				ServiceName: ServiceName,
				ServiceId:   ServiceId,
			}
			response, err := rpcClient.HeartbeatDetection(context.Background(), &request)
			if err != nil {
				println(err)
			}
			fmt.Println(response)
		case str := <-exit:
			fmt.Println(str)
			return
		}
	}
}

func (sc *SSRADClint) Close() {
	exit <- "time to go"
}

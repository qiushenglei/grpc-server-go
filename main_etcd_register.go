package main

import (
	"flag"
	"fmt"
	"github.com/qiushenglei/grpc-server-go/proto"
	"github.com/qiushenglei/grpc-server-go/rpcserver"
	"github.com/qiushenglei/grpc-server-go/service"
	"google.golang.org/grpc"
)

var (
	port1 = flag.Int("p", 50002, "server port1")
)

const (
	key1 string = "vector_rpc_server1"
)

func main() {
	//解析标签
	flag.Parse()
	//etcd服务注册中心需要用到的配置
	config := &rpcserver.RpcServiceConfig{
		//etcd当中key的前缀
		Key: key1,
		//rpc监听的地址和端口号
		ServerAddress: fmt.Sprintf("127.0.0.1:%d", *port1),
		//etcd集群地址
		Endpoints: []string{"127.0.0.1:2380"},
	}

	//注册etcd中心以及grpc服务
	if server, err := rpcserver.NewRpcService(config, func(server1 *grpc.Server) {
		proto.RegisterUserServiceServer(server1, &service.UserService{})
	}); err == nil {
		if err := server.Run(); err != nil {
			panic(err)
		}
		fmt.Println("started")
	}
}

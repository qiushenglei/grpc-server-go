package main

import (
	"flag"
	"github.com/qiushenglei/grpc-server-go/proto"
	"github.com/qiushenglei/grpc-server-go/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
)

var tls = flag.Int("t", 0, "need tls")
var certFile = flag.String("cert", "", "cert file path")
var keyFile = flag.String("key", "", "key file path")

func main1() {

	flag.Parse()

	//监听
	lis, _ := net.Listen("tcp", "127.0.0.1:8899")

	// 证书验证，没有的话可以省略
	var opts []grpc.ServerOption
	if *tls == 1 {
		if *certFile == "" {
			*certFile = "x509/server_cert.pem"
		}
		if *keyFile == "" {
			*keyFile = "x509/server_key.pem"
		}
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	// 起GRPC服务

	// 生成 grpc server结构体
	grpcServer := grpc.NewServer(opts...)

	// 生成Instance结构体
	ins := &service.UserService{
		Engine: service.NewMysqlEngine(),
	}

	// 给rpc服务注册,这个方法是proto-go生成的
	proto.RegisterUserServiceServer(grpcServer, ins)

	// 开启服务，  这个地方我们不需要for死循环 grpc当中会自动帮助我们实时监听
	grpcServer.Serve(lis)
}

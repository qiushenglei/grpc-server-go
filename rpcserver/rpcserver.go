package rpcserver

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type (
	//函数作为值来传递 类型是自定义的类型RpcServiceFunc 传入的值是grpc.NewServer()指针类型*grpc.Server
	RpcServiceFunc func(server *grpc.Server)
	//定义一个结构体来承装register和注册中心和传递过来的函数 这个函数是干嘛的往下看
	RpcService struct {
		//register注册中心
		register *Register
		//作为值传递过来的函数
		//demo.RegisterDemoServiceServer(server1, &rpcserverimpl.DemoServiceServerImpl{})
		//上边RegisterDemoServiceServer是形成的 pb.go文件里面的哦  注册服务 第二个参数是你实现接口的struct结构体
		//总的来说是为了实现注册rpc服务的一部分
		rpcServiceFunc RpcServiceFunc
	}
	//将一些配置搞到里面去 etcd服务注册中心用到的比如key的前缀   rpc服务地址  etcd集群地址配置
	RpcServiceConfig struct {
		Key           string
		ServerAddress string
		Endpoints     []string
	}
)

// 实例化各类操作
func NewRpcService(conf *RpcServiceConfig, rpcServiceFun RpcServiceFunc) (*RpcService, error) {

	//链接etcd服务
	client3, err := clientv3.New(
		clientv3.Config{
			Endpoints: conf.Endpoints, //etcd集群地址配置
			//你还可以设置更过的参数 如果开启了auth验证 也可以配置username和password
		},
	)
	if err != nil {
		return nil, err
	}
	fmt.Errorf("md fuck %s", "ssfsdfs")

	// clientv3有问题,没有连接成功也不返回错误，判断一次status
	res, err := client3.Status(context.Background(), conf.Endpoints[0])
	if err != nil {
		fmt.Println(res)
		return nil, err
	}

	//返回的结构体当中有 注册好的Register注册中心结构体
	return &RpcService{
		//实例化服务注册中心 方便接下来的调用
		register: NewRegister(conf.Key, client3, conf.ServerAddress),
		//返回函数
		rpcServiceFunc: rpcServiceFun,
	}, nil
}

// 运行etcd服务注册中心和监听rpc服务
func (s *RpcService) Run() error {
	//监听rpc服务
	listen, err := net.Listen("tcp", s.register.GetServiceAddress())
	if err != nil {
		return err
	}
	log.Printf("Rpc server listen at:%s:", s.register.GetServiceAddress())
	//etcd服务注册中心  这里真正的调起服务注册中心的reg（）注册方法
	s.register.Reg()
	//etcd的善后工作
	//我们的rpc服务有没有挂掉 服务器有没有宕机......
	//所以我们需要建立一个通道 告诉etcd我的服务挂掉了 服务器宕机了...... 赶紧把我的合同解约删掉key 好让接下里的请求不再调取我的rpc服务
	//所以 我们要接收来自系统的信号
	s.deadNotify()
	//grpc  实例化grpc服务 返回的是 *grpc.Server
	server := grpc.NewServer()
	//grpc  s.rpcServiceFunc()其实就是在调用传递进来的函数  函数要求传入的是 *grpc.Server 所以就把grpc.NewServer()传递进去就可以啦
	//这样就真的将我们实现的DemoServiceServerImpl注册进到了grpc服务当中 就等着下边启动就完事了
	s.rpcServiceFunc(server)
	//server.Serve(listen) 启动grpc服务 进行实时监听 程序不会中断！
	if err := server.Serve(listen); err != nil {
		return err
	}
	return nil

}

// etcd的善后工作
// 我们的rpc服务有没有挂掉 服务器有没有宕机......
// 所以我们需要建立一个通道 告诉etcd我的服务挂掉了 服务器宕机了...... 赶紧把我的合同解约删掉key 好让接下里的请求不再调取我的rpc服务
// 所以 我们要接收来自系统的信号
func (s *RpcService) deadNotify() error {
	ch := make(chan os.Signal, 1)
	//接收系统信号写到ch管道当中去
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		sign := <-ch
		log.Printf("signal.notify %v", sign)
		//服务停掉之后 删除掉注册进etcd中心的key 我的服务都停了还留着你干啥用呢！
		s.register.UnReg()
	}()
	return nil
}

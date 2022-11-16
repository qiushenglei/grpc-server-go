package rpcserver

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"
	"log"
	"time"
)

type (
	Register struct {
		key           string           //前缀
		client3       *clientv3.Client //etcd的链接
		serverAddress string           //服务地址
		stop          chan bool        //假设一个rpc服务宕机了我们需要从etcd当中删除 利用channel管道技术传递bool值
		interval      time.Duration    //心跳周期  我们需要跟etcd保持联系 隔一段时间就去联系一下证明自己还活着
		leaseTime     int64            //租赁的时间 都是有时间限制的
	}
)

func NewRegister(key string, client3 *clientv3.Client, serverAddress string) *Register {
	return &Register{
		key:           key,
		client3:       client3,
		serverAddress: serverAddress,
		//心跳的周期一定要小于租赁的周期  不然会存在真空期
		interval:  3 * time.Second,
		leaseTime: 15,
		stop:      make(chan bool, 1),
	}
}

// 注册
func (r *Register) Reg() {
	k := r.makeKye()
	//心跳
	go func() {
		//每次心跳周期设置的时间都会往通道里面塞入值
		t := time.NewTicker(r.interval)
		//这里起一个死循环 不停的循环
		for {
			//租赁 生成租赁合同哦
			lgs, err := r.client3.Grant(context.TODO(), r.leaseTime)
			if nil != err {
				panic(err)
			}
			//判断key是否存在 存在则更新 不存在则写入
			if _, err := r.client3.Get(context.TODO(), k); err != nil {
				//如果没有发现key值的存在
				if err == rpctypes.ErrKeyNotFound {
					//没有发现key那就往里面添加喽  k就是key+服务器地址   value就是服务器地址 然后再来个租赁周期（得有个过期时间啊）
					if _, err := r.client3.Put(context.TODO(), k, r.serverAddress, clientv3.WithLease(lgs.ID)); err != nil {
						//如果有错误那么我们就直接退出了
						panic(err)
					}
				} else {
					//既然没发现次key的存在 err还不为空 那就是未知错误来处理
					panic(err)
				}
			} else {
				//有key的存在那么我们就去更新这个key
				if _, err := r.client3.Put(context.TODO(), k, r.serverAddress, clientv3.WithLease(lgs.ID)); err != nil {
					panic(err)
				}
			}
			//这里需要对select有深入的理解https://studygolang.com/articles/7203 <-t.C 以及 <-r.stop会一直等待
			select {
			//通道里面没有值 程序会阻塞在这里
			case ttl := <-t.C:
				log.Println(ttl)
			//如果收到了停止信号 则整个协程结束 即心跳结束
			case <-r.stop:
				return
			}

		}
	}()
}

// 取消注册
// 比如我的服务端挂掉了 我需要取消key即删除掉key
func (r *Register) UnReg() {
	//1.首先要停止心跳
	r.stop <- true
	//为了防止多线程下出现死锁的问题  channel管道就是为了协程和协程之间的通讯 上边设置了true那么注册中心里面的心跳程序就死掉了 初始化一下r.stop
	r.stop = make(chan bool, 1)
	k := r.makeKye()
	if _, e := r.client3.Delete(context.TODO(), k); e != nil {
		panic(e)
	} else {
		//打印哥日志看看喽
		log.Printf("%s unreg success", k)
	}

	return
}

// 生成key策略
func (r *Register) makeKye() string {
	return fmt.Sprintf("%s_%s", r.key, r.serverAddress)
}

func (r *Register) GetServiceAddress() string {
	return r.serverAddress
}

# go-grpc-server

## 介绍

测试go-grpc服务 和 etcd服务注册中心

## 目录

```bash
│  go.mod
│  go.sum
│  main.go
│  main_etcd_register.go
│  README.md
├─models
│      user.go
│
├─proto
│      user.pb.go
│      user.proto
│      user_grpc.pb.go
│
├─rpcserver
│      register.go
│      rpcserver.go
│
├─service
│      userService.go
│
└─third_party
    └─google
```

- `main.go` 起grpc服务
- `main_etcd_register.go` 起grpc服务，并注册到etcd中
- `proto`目录存放`.proto`文件和生成的`.pb`文件
- `rpcserver`目录存放的是rpc服务启动过程
- `service`目录是具体服务
- `third_party`目录是proto引入的三方包，可以存放go本地包

## 启动

- 生成`.pb`文件
    ```bash
    protoc -I google --proto_path=proto --go_out=. --go-grpc_out=. user.proto
    ```
- 启动服务
  ```bash
  go build -
- ```
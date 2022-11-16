package service

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql" //不能忘记导入
	"github.com/go-xorm/xorm"
	"github.com/qiushenglei/grpc-server-go/models"
	"github.com/qiushenglei/grpc-server-go/proto"
	"xorm.io/core"
)

// 实现.proto的service定义的接口
type UserService struct {
	// 业务侧是查询sql，所以提前注册到rpc服务内
	Engine *xorm.Engine
	// grpc规范必须要引入，否则我无法实现这个接口，因为她有个must小写方法
	proto.UnimplementedUserServiceServer
}

// 数据库操作引擎
func NewMysqlEngine() *xorm.Engine {
	engine, err := xorm.NewEngine("mysql", "root:root@tcp(127.0.0.1:333306)/testdb?charset=utf8")
	if err != nil {
		panic(err)
	}

	engine.ShowSQL(true)
	engine.Logger().SetLevel(core.LOG_DEBUG)
	engine.SetMaxOpenConns(10)

	//返回引擎
	return engine
}

// 实现proto文件定义的rpc暴露的方法
func (s *UserService) GetUserList(ctx context.Context, request *proto.UserRequest) (*proto.UserListReply, error) {
	var userList []models.User
	err := s.Engine.Where("id > ?", 0).Find(&userList)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(userList)
	fmt.Println("===========================================")

	//因为返回的是[]*proto.User
	res := make([]*proto.User, len(userList))
	//遍历数据库查询结果 然后重新塞入到新的切片当中去
	for _, u := range userList {
		res = append(res, &proto.User{
			Name:  u.Name,
			Phone: u.Phone,
			Sex:   u.Sex,
			Id:    u.ID,
		})
	}

	return &proto.UserListReply{User: res}, err
}

func (s *UserService) GetUser(ctx context.Context, request *proto.UserRequest) (*proto.UserReply, error) {
	// 业务code
	return &proto.UserReply{User: &proto.User{Id: 10, Name: "qiushenglei"}}, nil
}

package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/proto"
	"mxshop-api/user-web/utils/otgrpc"
)

func InitSrvConn()  {
	//连接用户grpc服务
	consulInfo := global.ServerConfig.ConsulInfo
	//grpc consul最小连接数负载均衡
	conn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)

	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 [用户服务]失败", "msg", err.Error())
		return
	}

	//grpc实例，调用接口
	global.UserSrvClient = proto.NewUserClient(conn)
}

func InitSrvConnV2()  {
	//注册中心获取consul
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)

	userSrvHost := ""
	userSrvPort := 0

	client, err := api.NewClient(cfg)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 [consul]失败", "msg", err.Error())
		return
	}

	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%s\"", global.ServerConfig.UserSrvInfo.Name))
	if err != nil {
		zap.S().Fatal("[GetUserList] 获取 [consul服务列表]失败", "msg", err.Error())
		return
	}

	for _, v := range data {
		userSrvHost = v.Address
		userSrvPort = v.Port
		break
	}

	if userSrvHost == "" {
		zap.S().Fatal("[GetUserList] 获取 [consul服务列表]失败")
		return
	}

	//连接用户grpc服务
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 [用户服务失败]", "msg", err.Error())
	}

	//全局变量就不关闭了
	//defer userConn.Close()

	//grpc实例，调用接口
	global.UserSrvClient = proto.NewUserClient(userConn)
}

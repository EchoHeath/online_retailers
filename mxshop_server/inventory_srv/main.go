package main

import (
	"flag"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"mxshop_server/inventory_srv/global"
	"mxshop_server/inventory_srv/handler"
	"mxshop_server/inventory_srv/initialize"
	"mxshop_server/inventory_srv/proto"
	"mxshop_server/inventory_srv/utils"
	"mxshop_server/inventory_srv/utils/register/consul"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 50051, "端口号")
	flag.Parse()

	//初始化配置
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()

	//如果没有从命令行启动，就默认分配端口
	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}
	zap.S().Info("ip: ", *IP)
	zap.S().Info("port: ", *Port)

	server := grpc.NewServer()
	proto.RegisterInventoryServer(server, &handler.InventoryServer{})

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	//注册健康服务检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	//服务注册
	register_client := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	err = register_client.Register(global.ServerConfig.Host, *Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("服务注册失败:", err.Error())
	}
	zap.S().Debugf("启动服务器, 端口： %d", *Port)

	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = register_client.DelRegister(serviceId); err != nil {
		zap.S().Info("注销失败:", err.Error())
	}else{
		zap.S().Info("注销成功:")
	}
}

package main

import (
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/opentracing/opentracing-go"
	uuid "github.com/satori/go.uuid"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"mxshop_server/user_srv/global"
	"mxshop_server/user_srv/handler"
	"mxshop_server/user_srv/initialize"
	"mxshop_server/user_srv/proto"
	"mxshop_server/user_srv/utils"
	"mxshop_server/user_srv/utils/otgrpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 0, "端口号")
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

	// jaeger配置
	jaegerInfo := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "127.0.0.1:6831",
		},
		ServiceName: "mx_shop",
	}
	tracer, closer, err := jaegerInfo.NewTracer(config.Logger(jaeger.StdLogger))
	opentracing.SetGlobalTracer(tracer)

	//集成server并注册
	server := grpc.NewServer(grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(tracer)))
	proto.RegisterUserServer(server, &handler.UserServer{})

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	//注册健康服务检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	//服务注册
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(client)
	}

	//生成注册对象
	registration := new(api.AgentServiceRegistration)
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	registration.ID = serviceID
	registration.Name = global.ServerConfig.Name
	registration.Address = global.ServerConfig.ConsulInfo.Host
	registration.Port = *Port
	registration.Tags = []string{"user", "srv"}

	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, *Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "15s",
	}

	registration.Check = check

	//注册服务到consul
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic("failed to register consul:" + err.Error())
	}

	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	_ = closer.Close()

	if err = client.Agent().ServiceDeregister(serviceID); err != nil {
		zap.S().Info("注销失败")
	}
	zap.S().Info("注销成功")
}

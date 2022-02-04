package main

import (
	"fmt"
	"go.uber.org/zap"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/initialize"
)

func main() {
	//1.初始化logger
	initialize.InitLogger()

	//2.初始化配置文件
	initialize.InitConfig()

	//3.初始化router
	router := initialize.Router()

	//4.初始化翻译
	if err := initialize.InitTrans("zh"); err != nil {
		panic(err)
	}

	//5.初始化srv服务
	initialize.InitSrvConn()

	//6.初始化sentinel
	initialize.Sentinel()

	//debug := initialize.GetEnvInfo("MXSHOP_DEBUG")
	//if !debug { //本地端口无需变动，服务端口需要动态生成并注册到consul
	//	port, err := utils.GetFreePort()
	//	if err == nil {
	//		global.ServerConfig.Port = port
	//	}
	//}

	zap.S().Infof("启动服务器，端口: %d", global.ServerConfig.Port)
	if err := router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动失败:", err.Error())
	}
}

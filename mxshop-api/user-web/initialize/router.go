package initialize

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/user-web/middlewares"
	"mxshop-api/user-web/router"
)

func Router() *gin.Engine {
	Router := gin.Default()
	//配置跨域
	Router.Use(middlewares.Cors())

	//路由配置
	ApiGroup := Router.Group("/u/v1")
	//配置链路追踪trace
	ApiGroup.Use(middlewares.Trace())

	router.InitUserRouter(ApiGroup)
	router.InitBaseRouter(ApiGroup)
	return Router
}

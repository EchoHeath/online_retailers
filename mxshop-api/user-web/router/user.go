package router

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/user-web/api"
	"mxshop-api/user-web/middlewares"
)

func InitUserRouter(router *gin.RouterGroup) {
	userGroup := router.Group("user")
	{
		userGroup.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
		userGroup.POST("pwd_login", api.PassWordLogin)
		userGroup.POST("register", api.Register)
	}
}

package global

import (
	ut "github.com/go-playground/universal-translator"
	"mxshop-api/user-web/config"
	"mxshop-api/user-web/proto"
)

var (
	ServerConfig  *config.ServerConfig = &config.ServerConfig{}
	NacosConfig  *config.NacosConfig = &config.NacosConfig{}
	Trans         ut.Translator
	UserSrvClient proto.UserClient
)

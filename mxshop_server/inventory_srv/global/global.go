package global

import (
	"gorm.io/gorm"
	"mxshop_server/inventory_srv/config"
)

var (
	DB *gorm.DB
	ServerConfig config.ServiceConfig
	NacosConfig config.NacosConfig
)

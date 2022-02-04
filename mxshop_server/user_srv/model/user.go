package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        int32     `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"column:add_time"`
	UpdatedAt time.Time `gorm:"column:update_time"`
	DeletedAt gorm.DeletedAt
	IsDeleted bool
}

type User struct {
	BaseModel
	Mobile   string     `gorm:"index:idx_mobile,unique;type:varchar(11) comment '手机号';not null;default:''"`
	Password string     `gorm:"type:varchar(100) comment '密码';not null;default:''"`
	NickName string     `gorm:"type:varchar(20) comment '昵称';not null;default:''"`
	Birthday *time.Time `gorm:"type:datetime comment '生日'"`
	Gender   string     `gorm:"type:varchar(6) comment 'female女,male男';default:male;not null"`
	Role     int        `gorm:"type:int(1) comment '权限管理,1普通用户,2管理员';not null;default:1"`
}

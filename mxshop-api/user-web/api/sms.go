package api

import (
	"context"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"mxshop-api/user-web/forms"
	"mxshop-api/user-web/global"
	"net/http"
	"strings"
	"time"
)

func GenerateSmsCode(width int) string {
	//生成短信验证码
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

func SendSms(c *gin.Context) {
	sendSmsForm := forms.SendSmsForm{}
	if err := c.ShouldBind(&sendSmsForm); err != nil {
		HandlerValidatorErr(c, err)
		return
	}

	client, err := dysmsapi.NewClientWithAccessKey("cn-beijing", global.ServerConfig.AliSmsInfo.ApiKey, global.ServerConfig.AliSmsInfo.ApiSecrect)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "发送失败",
		})
		return
	}
	smsCode := GenerateSmsCode(6)
	mobile := "1503871****"
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"

	request.PhoneNumbers = mobile                            //接收短信的手机号码
	request.SignName = "阿里大于测试专用"                            //短信签名名称
	request.TemplateCode = "SMS_209335004"                   //短信模板ID
	request.TemplateParam = "{\"code\":\"" + smsCode + "\"}" //短信模板变量对应的实际值，JSON格式

	_, err = client.SendSms(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "发送失败",
		})
		return
	}
	//连接redis
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})

	rdb.Set(context.Background(), mobile, smsCode, time.Duration(global.ServerConfig.RedisInfo.Expire)*time.Second)

	c.JSON(http.StatusOK, gin.H{
		"msg": "发送成功",
	})
}

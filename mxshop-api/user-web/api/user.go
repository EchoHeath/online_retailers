package api

import (
	"context"
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mxshop-api/user-web/forms"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/global/response"
	"mxshop-api/user-web/middlewares"
	"mxshop-api/user-web/models"
	"mxshop-api/user-web/proto"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

//HandleGrpcErrorToHttp 将grpc的code转化为http状态码
func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其他错误" + e.Message(),
				})
			}
			return
		}

	}
}

//GetUserList 获取用户列表
func GetUserList(c *gin.Context) {
	claims, _ := c.Get("claims")
	currentUser, _ := claims.(*models.CustomClaims)
	zap.S().Infof("访问用户： %d", currentUser.ID)

	pnInt, _ := strconv.Atoi(c.DefaultQuery("pn", "0"))
	pSizeInt, _ := strconv.Atoi(c.DefaultQuery("psize", "10"))

	//sentinel限流
	e, b := sentinel.Entry("api", sentinel.WithTrafficType(base.Inbound))
	if b != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"msg": "请求过于频繁，请稍后重试",
		})
		return
	}

	rsp, err := global.UserSrvClient.GetUserList(context.WithValue(context.Background(), "ginContext", c), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询 [用户列表]失败", "msg", err.Error())
		HandleGrpcErrorToHttp(err, c)
		return
	}
	e.Exit()
	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		user := response.UserResponse{
			Id:       value.Id,
			NikeName: value.NickName,
			Birthday: response.JsonTime(time.Unix(int64(value.BirthDay), 0)),
			Gender:   value.Gender,
			Mobile:   value.Mobile,
		}

		result = append(result, user)
	}

	c.JSON(http.StatusOK, result)
}

func HandlerValidatorErr(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"err": removeTopStruct(errs.Translate(global.Trans)),
	})
}

//PassWordLogin 登录
func PassWordLogin(c *gin.Context) {
	passwordLoginForm := forms.PassWordLoginForm{}
	if err := c.ShouldBind(&passwordLoginForm); err != nil {
		HandlerValidatorErr(c, err)
		return
	}

	if !store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, false) { //false：测试用，每次验证验证码不清除，可复用验证码
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "验证码错误",
		})
		return
	}

	zap.S().Infof("grpc host: %s", global.ServerConfig.UserSrvInfo.Host)

	if rsp, err := global.UserSrvClient.GetUserByMobile(c, &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	}); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
			case codes.Unknown:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": "用户不存在",
				})
			default:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": "登录失败",
				})
			}
			return
		}
	} else {
		if passRsp, err := global.UserSrvClient.CheckPassWord(c, &proto.PassWordCheckInfo{
			PassWord:          passwordLoginForm.PassWord,
			EncryptedPassWord: rsp.PassWord,
		}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "登陆失败",
			})
			return
		} else {
			if passRsp.Success {
				//生成token
				j := middlewares.NewJWT()
				claims := models.CustomClaims{
					ID:          uint(rsp.Id),
					NickName:    rsp.NickName,
					AuthorityId: uint(rsp.Role),
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Unix(),
						ExpiresAt: time.Now().Unix() + 60*60*30,
						Issuer:    "test",
					},
				}

				token, err := j.CreateToken(claims)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"msg": "生成token失败",
					})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"id":         rsp.Id,
					"nick_name":  rsp.NickName,
					"token":      token,
					"expired_at": (time.Now().Unix() + 60*60*30) * 1000,
				})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "密码错误",
				})
			}
			return
		}
	}
}

func Register(c *gin.Context) {
	//用户注册
	registerForm := forms.RegisterForm{}
	if err := c.ShouldBind(&registerForm); err != nil {
		HandlerValidatorErr(c, err)
		return
	}

	//连接redis
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})

	res, err := rdb.Get(context.Background(), registerForm.Mobile).Result()
	if err == redis.Nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "验证码错误",
		})
		return
	} else {
		if res != registerForm.Code {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "验证码错误",
			})
			return
		}
	}

	_, err = global.UserSrvClient.CreateUser(c, &proto.CreateUserInfo{
		NickName: registerForm.Mobile,
		PassWord: registerForm.PassWord,
		Mobile:   registerForm.Mobile,
	})
	if err != nil {
		zap.S().Errorw("[Register] 创建 [用户]失败", "msg", err.Error())
		HandleGrpcErrorToHttp(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "创建成功",
	})
	return

}

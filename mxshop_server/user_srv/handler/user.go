package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"mxshop_server/user_srv/global"
	"mxshop_server/user_srv/model"
	"mxshop_server/user_srv/proto"
	"strings"
	"time"
)

type UserServer struct {
	proto.UnimplementedUserServer
	Ctx context.Context
}

func ModelToResponse(user model.User) proto.UserInfoResponse {
	resp := proto.UserInfoResponse{
		Id:       user.ID,
		PassWord: user.Password,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     int32(user.Role),
	}
	if user.Birthday != nil {
		resp.BirthDay = uint64(user.Birthday.Unix())
	}
	return resp
}

//Paginate 分页逻辑
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

//GetUserList get user info lists.
func (u *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	//链路追踪
	parentSpan := opentracing.SpanFromContext(ctx)
	listSpan := opentracing.GlobalTracer().StartSpan("get_user_list", opentracing.ChildOf(parentSpan.Context()))

	//获取用户列表
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	fmt.Println("用户列表")
	rsp := &proto.UserListResponse{}
	rsp.Total = int32(result.RowsAffected)

	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)

	for _, user := range users{
		userInfoRsp := ModelToResponse(user)
		rsp.Data = append(rsp.Data, &userInfoRsp)
	}

	listSpan.Finish()
	return rsp, nil
}

//GetUserByMobile get user info by mobile.
func (u *UserServer) GetUserByMobile(c context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	res := global.DB.Where("mobile = ?", req.Mobile).First(&user)

	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}
	data := ModelToResponse(user)
	return &data, nil
}

//GetUserByID get user info by ID.
func (u *UserServer) GetUserByID(c context.Context, req *proto.IDRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	res := global.DB.First(&user, req.Id)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}
	data := ModelToResponse(user)
	return &data, nil

}

//CreateUser creates user.
func (u *UserServer) CreateUser(c context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	var user model.User
	res := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if res.RowsAffected == 1 {
		return nil, status.Error(codes.AlreadyExists, "用户已存在")
	}

	user.Mobile = req.Mobile
	user.NickName = req.NickName

	option := &password.Options{16, 100, 32, sha512.New}
	salt, pwd := password.Encode(req.PassWord, option)
	user.Password = fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, pwd)

	res = global.DB.Create(&user)
	if res.Error != nil {
		return nil, status.Error(codes.Internal, res.Error.Error())
	}

	data := ModelToResponse(user)
	return &data, nil
}

//UpdateUser updates user info by ID.
func (u *UserServer) UpdateUser(c context.Context, req *proto.UpdateUserInfo) (*emptypb.Empty, error) {
	var user model.User
	res := global.DB.First(&user, req.Id)
	if res.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}

	birthDay := time.Unix(int64(req.BirthDay), 0)
	user.Birthday = &birthDay
	user.Password = req.PassWord
	user.NickName = req.NickName
	res = global.DB.Save(&user)
	if res.Error != nil {
		return nil, status.Error(codes.Internal, res.Error.Error())
	}

	return &emptypb.Empty{}, nil
}

//CheckPassWord check password.
func (u *UserServer) CheckPassWord(c context.Context, req *proto.PassWordCheckInfo) (*proto.CheckResponse, error) {
	options := &password.Options{16, 100, 32, sha512.New}
	passwordInfo := strings.Split(req.EncryptedPassWord, "$")
	check := password.Verify(req.PassWord, passwordInfo[2], passwordInfo[3], options)
	return &proto.CheckResponse{
		Success: check,
	}, nil
}

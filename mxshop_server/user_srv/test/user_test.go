package test

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"mxshop_server/user_srv/proto"
	"testing"
)

func TestGetUserList(t *testing.T) {
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("failed to grpc dial: %v", err)
		return
	}
	defer conn.Close()

	c := proto.NewUserClient(conn)
	ctx := context.Background()
	cases := []struct {
		name    string
		op      func() error
		wantErr bool
	}{
		{
			name: "get_user_list",
			op: func() error {
				data, err := c.GetUserList(ctx, &proto.PageInfo{
					Pn:    1,
					PSize: 2,
				})
				if err != nil {
					return err
				}
				fmt.Println(data)
				return nil
			},
		},
		//{
		//	name: "get_user_by_mobile",
		//	op: func() error {
		//		user, err := c.GetUserByMobile(ctx, &proto.MobileRequest{
		//			Mobile: "13888888894",
		//		})
		//		if err != nil {
		//			return err
		//		}
		//		fmt.Println(user)
		//		return nil
		//	},
		//},
		//{
		//	name: "get_user_by_id",
		//	op: func() error {
		//		user, err := c.GetUserByID(ctx, &proto.IDRequest{
		//			Id: 3,
		//		})
		//		if err != nil {
		//			return err
		//		}
		//		fmt.Println(user)
		//		return nil
		//	},
		//},
		//{
		//	name: "create_user",
		//	op: func() error {
		//		user, err := c.CreateUser(ctx, &proto.CreateUserInfo{
		//			NickName: "test",
		//			Mobile:   "13888888999",
		//			PassWord: "admin123",
		//		})
		//		if err != nil {
		//			return err
		//		}
		//		fmt.Println(user)
		//		return nil
		//	},
		//},
		//{
		//	name: "update_user",
		//	op: func() error {
		//		user, err := c.UpdateUser(ctx, &proto.UpdateUserInfo{
		//			Id:       "10",
		//			NickName: "nickname",
		//		})
		//		if err != nil {
		//			return err
		//		}
		//		fmt.Println(user)
		//		return nil
		//	},
		//},
		//{
		//	name: "check_password",
		//	op: func() error {
		//		word, err := c.CheckPassWord(ctx, &proto.PassWordCheckInfo{
		//			PassWord:          "admin123",
		//			EncryptedPassWord: "$pbkdf2-sha512$eckQ7qCVEIrSh0Kn$6a6f10de8190ebcf2beb9adbc37493f009c16f03bf504d652612a45fd7500fb4",
		//		})
		//		if err != nil {
		//			return err
		//		}
		//		if !word.Success {
		//			return fmt.Errorf("密码错误")
		//		}
		//		return nil
		//	},
		//},
	}

	for _, cc := range cases {
		err := cc.op()
		if cc.wantErr {
			if err == nil {
				t.Errorf("%s: want err;got done", cc.name)
			} else {
				continue
			}
		}

		if err != nil {
			t.Errorf("%s: operation failed: %v", cc.name, err)
			return
		}
	}

}

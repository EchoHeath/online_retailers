package handler

import (
	"context"
	"fmt"
	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"mxshop_server/inventory_srv/global"
	"mxshop_server/inventory_srv/model"
	"mxshop_server/inventory_srv/proto"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServer
}

//SetInv 设置库存
func (i *InventoryServer) SetInv(c context.Context, req *proto.GoodsInvInfo) (*emptypb.Empty, error) {
	var inv model.Inventory
	if result := global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv); result.Error != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "库存查询失败")
	}

	inv.Goods = req.GoodsId
	inv.Stocks = req.Num

	global.DB.Save(&inv)
	return &emptypb.Empty{}, nil
}

func (i *InventoryServer) InvDetail(c context.Context, req *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inv model.Inventory
	if result := global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv); result.RowsAffected == 0 {
		return &proto.GoodsInvInfo{}, status.Errorf(codes.Internal, "库存查询不存在")
	}

	return &proto.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num:     inv.Stocks,
	}, nil
}

func (i *InventoryServer) Sell(c context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "localhost:6379",
	})
	pool := goredis.NewPool(client)
	rs := redsync.New(pool)

	tx := global.DB.Begin()

	for _, goodInfo := range req.GoodsInfo {
		mutex := rs.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))
		if err := mutex.Lock(); err != nil {
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, "获取redis锁错误")
		}

		var inv model.Inventory
		//if result := tx.Clauses(clause.Locking{Strength: "update"}).Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
		//	tx.Rollback()
		//	return nil, status.Errorf(codes.InvalidArgument, "库存查询不存在")
		//}

		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.InvalidArgument, "库存查询不存在")
		}

		if inv.Stocks < goodInfo.Num {
			tx.Rollback()
			return nil, status.Errorf(codes.ResourceExhausted, "库存数量不足")
		}

		inv.Stocks -= goodInfo.Num
		tx.Save(&inv)
		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.ResourceExhausted, "redis解锁失败")
		}
	}

	tx.Commit()
	return &emptypb.Empty{}, nil
}

func (i *InventoryServer) Reback(c context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	tx := global.DB.Begin()

	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.InvalidArgument, "库存查询不存在")
		}

		inv.Stocks += goodInfo.Num
		global.DB.Save(&inv)
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}

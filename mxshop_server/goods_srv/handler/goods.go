package handler

import (
	"mxshop_server/goods_srv/proto"
)

type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

//商品接口
//func (g *GoodsServer) GoodsList(c context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
//	return &proto.GoodsListResponse{}, nil
//}
//func (g *GoodsServer) BatchGetGoods(c context.Context, req *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
//	return &proto.GoodsListResponse{}, nil
//}
//func (g *GoodsServer) CreateGoods(c context.Context, req *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
//	return &proto.GoodsInfoResponse{}, nil
//
//}
//func (g *GoodsServer) DeleteGoods(c context.Context, req *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
//	return &emptypb.Empty{}, nil
//}
//func (g *GoodsServer) UpdateGoods(c context.Context, req *proto.CreateGoodsInfo) (*emptypb.Empty, error) {
//	return &emptypb.Empty{}, nil
//}
//func (g *GoodsServer) GetGoodsDetail(c context.Context, req *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
//	return &proto.GoodsInfoResponse{}, nil
//}

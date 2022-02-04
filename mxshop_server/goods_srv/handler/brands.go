package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"mxshop_server/goods_srv/global"
	"mxshop_server/goods_srv/model"
	"mxshop_server/goods_srv/proto"
)

//品牌和轮播图
func (g *GoodsServer) BrandList(c context.Context, req *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	brandListResponse := proto.BrandListResponse{}
	var brands []model.Brands
	result := global.DB.Find(&brands)
	if result.Error != nil {
		return nil, result.Error
	}

	var brandResponses []*proto.BrandInfoResponse
	for _, brand := range brands {
		brandResponses = append(brandResponses, &proto.BrandInfoResponse{
			Id:   brand.ID,
			Name: brand.Name,
			Logo: brand.Logo,
		})
	}
	brandListResponse.Data = brandResponses
	brandListResponse.Total = int32(result.RowsAffected)

	return &brandListResponse, nil
}

func (g *GoodsServer) CreateBrand(c context.Context, req *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	if result := global.DB.First(&model.Brands{}); result.RowsAffected == 1 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌已存在")
	}
	brand := model.Brands{}
	brand.Name = req.Name
	brand.Logo = req.Logo
	result := global.DB.Create(&brand)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "品牌已存在")
	}
	return &proto.BrandInfoResponse{
		Id:   brand.ID,
		Name: brand.Name,
		Logo: brand.Logo,
	}, nil
}

func (g *GoodsServer) DeleteBrand(c context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {
	if res := global.DB.Delete(&model.Brands{}, req.Id); res.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}
	return &emptypb.Empty{}, nil
}

func (g *GoodsServer) UpdateBrand(c context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {
	brand := &model.Brands{}
	if result := global.DB.First(brand); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	if req.Name != "" {
		brand.Name = req.Name
	}

	if req.Logo != "" {
		brand.Logo = req.Logo
	}

	global.DB.Save(brand)
	return &emptypb.Empty{}, nil
}

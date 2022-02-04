package handler

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"mxshop_server/goods_srv/global"
	"mxshop_server/goods_srv/model"
	"mxshop_server/goods_srv/proto"
)

//商品分类
func (g *GoodsServer) GetAllCategorysList(c context.Context, req *emptypb.Empty) (*proto.CategoryListResponse, error) {
	var categorys []model.Category
	global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categorys)
	b, err := json.Marshal(&categorys)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "")
	}

	return &proto.CategoryListResponse{
		JsonData: string(b),
	}, nil
}

//获取子分类
func (g *GoodsServer) GetSubCategory(c context.Context, req *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	categoryListResponse := proto.SubCategoryListResponse{}

	var category model.Category
	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "分类不存在")
	}

	categoryListResponse.Info = &proto.CategoryInfoResponse{
		Id:             category.ID,
		Name:           category.Name,
		ParentCategory: category.ParentCategoryID,
		Level:          category.Level,
		IsTab:          category.IsTab,
	}

	var subCategorys []model.Category
	var subCategoryResponse []*proto.CategoryInfoResponse
	preload := "SubCategory"
	if category.Level == 1 {
		preload = "SubCategory.SubCategory"
	}

	global.DB.Where(&model.Category{ParentCategoryID: category.ID}).Preload(preload).Find(&subCategorys)
	for _, subCategory := range subCategorys {
		subCategoryResponse = append(subCategoryResponse, &proto.CategoryInfoResponse{
			Id:             subCategory.ID,
			Name:           subCategory.Name,
			ParentCategory: subCategory.ParentCategoryID,
			Level:          subCategory.Level,
			IsTab:          subCategory.IsTab,
		})
	}

	categoryListResponse.SubCategorys = subCategoryResponse
	categoryListResponse.Total = int32(len(subCategoryResponse))

	return &categoryListResponse, nil

} 

//func (g *GoodsServer) CreateCategory(c context.Context, req *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
//}
//func (g *GoodsServer) DeleteCategory(c context.Context, req *proto.DeleteCategoryRequest) (*emptypb.Empty, error) {
//}
//func (g *GoodsServer) UpdateCategory(c context.Context, req *proto.CategoryInfoRequest) (*emptypb.Empty, error) {
//}

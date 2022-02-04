package model

type Inventory struct {
	BaseModel
	Goods   int32 `gorm:"type:int;index" json:"-"`
	Stocks  int32 `gorm:"type:int" json:"-"`
	Version int32 `gorm:"type:int" json:"-"`
}

//type InventoryHistory struct {
//	user   int32
//	goods  int32
//	nums   int32
//	order  int32
//	status int32 //1预扣减，2已支付
//}

package request

type DemoOrderCreateReq struct {
	UserID       int64  `json:"user_id,omitempty"`
	Amount       string `json:"amount,omitempty"`
	OrderGoodsID int64  `json:"order_goods_id,omitempty"`
}

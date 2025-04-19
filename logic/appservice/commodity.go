package appservice

import "github.com/kackerx/go-mall/logic/domainservice"

type CommodityApp struct {
	commoditySvc *domainservice.CommoditySvc
}

func NewCommodityApp(
	commoditySvc *domainservice.CommoditySvc,
) *CommodityApp {
	return &CommodityApp{
		commoditySvc: commoditySvc,
	}
}

func (c *CommodityApp) InitCategoryData() error {
	return c.commoditySvc.InitCategory()
}

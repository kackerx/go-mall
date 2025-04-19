package domainservice

import (
	"encoding/json"

	"github.com/kackerx/go-mall/dal/dao"
	"github.com/kackerx/go-mall/logic/do"
	"github.com/kackerx/go-mall/resources"
)

type CommoditySvc struct {
	dao *dao.CommodityDao
}

func NewCommoditySvc(dao *dao.CommodityDao) *CommoditySvc {
	return &CommoditySvc{dao: dao}
}

func (c *CommoditySvc) InitCategory() error {
	categoryFileHandler, err := resources.LoadResourceFile("category.json")
	if err != nil {
		return err
	}

	var categoryList []*do.CommodityCategory
	if err := json.NewDecoder(categoryFileHandler).Decode(&categoryList); err != nil {
		return err
	}

	return nil
}

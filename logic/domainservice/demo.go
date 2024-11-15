package domainservice

import (
	"context"

	"github.com/kackerx/go-mall/common/errcode"
	"github.com/kackerx/go-mall/common/util"
	"github.com/kackerx/go-mall/dal/dao"
	"github.com/kackerx/go-mall/logic/do"
)

type DemoDomainSvc struct {
	ctx     context.Context
	DemoDao *dao.DemoDao
}

func NewDemoDomainSvc(ctx context.Context, demoDao *dao.DemoDao) *DemoDomainSvc {
	return &DemoDomainSvc{ctx: ctx, DemoDao: demoDao}
}

func (d *DemoDomainSvc) CreateDemoOrder(demoOrder *do.DemoOrder) (*do.DemoOrder, error) {
	demoOrder.Code = "110"
	orderPO, err := d.DemoDao.CreateDemoOrder(demoOrder)
	if err != nil {
		return nil, errcode.Wrap("创建订单失败", err)
	}

	// 事务中写明细表

	err = util.Copy(demoOrder, orderPO)
	return demoOrder, err
}

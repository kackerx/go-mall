package dao

import (
	"context"

	"github.com/kackerx/go-mall/common/util"
	"github.com/kackerx/go-mall/dal/model"
	"github.com/kackerx/go-mall/logic/do"
)

type DemoDao struct {
	ctx context.Context
}

func NewDemoDao(ctx context.Context) *DemoDao {
	return &DemoDao{ctx: ctx}
}

func (i *DemoDao) CreateDemoOrder(demoOrder *do.DemoOrder) (*model.DemoOrder, error) {
	po := new(model.DemoOrder)
	if err := util.Copy(po, demoOrder); err != nil {
		return nil, err
	}

	err := DB().WithContext(i.ctx).Create(po).Error
	return po, err
}

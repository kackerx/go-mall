package appservice

import (
	"context"

	"github.com/kackerx/go-mall/api/reply"
	"github.com/kackerx/go-mall/api/request"
	"github.com/kackerx/go-mall/common/errcode"
	"github.com/kackerx/go-mall/common/logger"
	"github.com/kackerx/go-mall/common/util"
	"github.com/kackerx/go-mall/dal/cache"
	"github.com/kackerx/go-mall/logic/do"
	"github.com/kackerx/go-mall/logic/domainservice"
)

// DemoAppSvc
// ApplicationService层职能
// 1. 调用领域服务的主逻辑方法
// 2. 把接口层的dto --> do
// 3. 执行依赖主逻辑的外围逻辑: 如发送创建成功的消息事件
// 4. 转换获取到的do --> vo
type DemoAppSvc struct {
	ctx           context.Context
	demoDoaminSvc *domainservice.DemoDomainSvc
}

func NewDemoAppSvc(ctx context.Context, demoDoaminSvc *domainservice.DemoDomainSvc) *DemoAppSvc {
	return &DemoAppSvc{ctx: ctx, demoDoaminSvc: demoDoaminSvc}
}

func (das *DemoAppSvc) CreateDemoOrder(orderReq *request.DemoOrderCreateReq) (*reply.DemoOrderResp, error) {
	demoOrder := new(do.DemoOrder)
	if err := util.Copy(demoOrder, orderReq); err != nil {
		return nil, errcode.Wrap("请求转换demoOrder失败", err)
	}

	// redis ex
	cache.SetDemoOrder(das.ctx, demoOrder)
	order, err := cache.GetDemoOrder(das.ctx, demoOrder.Code)
	if err != nil {
		return nil, err
	}

	logger.New(das.ctx).Info("redis data", "data", order)

	demoOrder, err = das.demoDoaminSvc.CreateDemoOrder(demoOrder)
	if err != nil {
		return nil, err
	}

	replyDemoOrder := new(reply.DemoOrderResp)
	if err = util.Copy(replyDemoOrder, demoOrder); err != nil {
		return nil, errcode.Wrap("demoOrderPo转换reply失败", err)
	}

	return replyDemoOrder, nil
}

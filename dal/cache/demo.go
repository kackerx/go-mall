package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kackerx/go-mall/common/enum"
	"github.com/kackerx/go-mall/common/log"
	"github.com/kackerx/go-mall/logic/do"
)

func SetDemoOrder(ctx context.Context, demoOrder *do.DemoOrder) error {
	bs, _ := json.Marshal(demoOrder)
	redisKey := fmt.Sprintf(enum.RedisKeyDemoOrderDetail, demoOrder.Code)
	_, err := Redis().Set(ctx, redisKey, bs, 0).Result()
	if err != nil {
		// ? redis没有gorm的logger接口, 所以在操作处自己打日志
		log.New(ctx).Error("redis set error", "err", err)
		return err
	}

	return nil
}

func GetDemoOrder(ctx context.Context, code string) (*do.DemoOrder, error) {
	redisKey := fmt.Sprintf(enum.RedisKeyDemoOrderDetail, code)
	bs, err := Redis().Get(ctx, redisKey).Bytes()
	if err != nil {
		log.New(ctx).Error("redis set error", "err", err)
		return nil, err
	}

	demoOrder := new(do.DemoOrder)
	json.Unmarshal(bs, demoOrder)

	return demoOrder, nil
}

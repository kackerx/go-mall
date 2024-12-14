package library

import (
	"context"
	"encoding/json"

	"github.com/kackerx/go-mall/common/logger"
	"github.com/kackerx/go-mall/common/util/httptool"
)

const (
	DetailURL = "https://ipwho.is"
)

type WhoisLib struct {
}

type WhoisDetail struct {
	IP      string `json:"ip"`
	Success bool   `json:"success"`
	City    string `json:"city"`
}

func NewWhoisLib() *WhoisLib {
	return &WhoisLib{}
}

func (w *WhoisLib) GetHostIPDetail(ctx context.Context) (detail *WhoisDetail, err error) {
	code, resp, err := httptool.Get(ctx, DetailURL, httptool.WithHeaders(map[string]string{
		"User-Agent": "HTTPie/3.2.3",
		"Host":       "ipwho.is",
	}))
	if err != nil {
		logger.New(ctx).Error("ipwho.is get detail error", "err", err, "status", code)
		return nil, err
	}

	detail = new(WhoisDetail)
	if err := json.Unmarshal(resp, detail); err != nil {
		return nil, err
	}

	return
}

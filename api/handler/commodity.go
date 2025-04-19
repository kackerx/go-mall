package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/kackerx/go-mall/common/app"
	"github.com/kackerx/go-mall/common/errcode"
	"github.com/kackerx/go-mall/logic/appservice"
)

type CommodityHandler struct {
	*Handler
	app *appservice.CommodityApp
}

func NewCommodityHandler(handler *Handler, app *appservice.CommodityApp) *CommodityHandler {
	return &CommodityHandler{Handler: handler, app: app}
}

func (a *CommodityHandler) InitCategory(c *gin.Context) {
	if err := a.app.InitCategoryData(); err != nil {
		app.NewResponse(c).Error(errcode.Wrap("InitCategoryData", err))
	} else {
		app.NewResponse(c).Success("ok")
	}
}

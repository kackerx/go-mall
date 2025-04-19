package router

import (
	"github.com/gin-gonic/gin"

	"github.com/kackerx/go-mall/api/handler"
)

func registerCommodityRoutes(rg *gin.RouterGroup, commodityHandler *handler.CommodityHandler) {
	g := rg.Group("/commodity/")

	g.POST("init-category", commodityHandler.InitCategory)
}

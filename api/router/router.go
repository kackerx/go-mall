package router

import (
	"github.com/gin-gonic/gin"

	"github.com/kackerx/go-mall/api/handler"
	"github.com/kackerx/go-mall/common/middleware"
)

func RegisterRoute(
	engin *gin.Engine,
	userHandler *handler.UserHandler,
	commodityHandler *handler.CommodityHandler,
) {
	engin.Use(middleware.StartTrace(), middleware.LogAccess(), middleware.GinPanicRecovery())
	routeGroup := engin.Group("")

	registerBuildingRoutes(routeGroup)
	registerUserRoutes(routeGroup, userHandler)
	registerCommodityRoutes(routeGroup, commodityHandler)
}

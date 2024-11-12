package router

import (
	"github.com/gin-gonic/gin"

	"github.com/kackerx/go-mall/common/middleware"
)

func RegisterRoute(engin *gin.Engine) {
	engin.Use(middleware.StartTrace(), middleware.LogAccess(), middleware.GinPanicRecovery())
	routeGroup := engin.Group("")

	RegisterBuildingRoutes(routeGroup)
}

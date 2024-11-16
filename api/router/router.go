package router

import (
	"github.com/gin-gonic/gin"

	"github.com/kackerx/go-mall/api/handler"
	"github.com/kackerx/go-mall/common/middleware"
)

func RegisterRoute(engin *gin.Engine, userHandler *handler.UserHandler) {
	engin.Use(middleware.StartTrace(), middleware.LogAccess(), middleware.GinPanicRecovery())
	routeGroup := engin.Group("")

	RegisterBuildingRoutes(routeGroup)
	RegisterUserRoutes(routeGroup, userHandler)
}

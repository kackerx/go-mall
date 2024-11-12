package router

import (
	"github.com/gin-gonic/gin"

	"github.com/kackerx/go-mall/api/handler"
)

func RegisterBuildingRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/building/")

	g.GET("ping", handler.TestErr)
	g.GET("resperr", handler.TestRespErr)
	g.GET("/respsuccess", handler.TestRespSuccess)
}

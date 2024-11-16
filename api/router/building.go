package router

import (
	"github.com/gin-gonic/gin"

	"github.com/kackerx/go-mall/api/handler"
	"github.com/kackerx/go-mall/common/middleware"
)

func RegisterBuildingRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/building/")

	g.GET("ping", handler.TestErr)
	g.GET("resperr", handler.TestRespErr)
	g.GET("/respsuccess", handler.TestRespSuccess)
	g.POST("/demoorder/add", handler.TestCreateDemoOrder)
	g.GET("/whois", handler.TestWhoisLibReq)
	g.GET("/gentoken", handler.TestMakeToken)
	g.GET("/gettoken", middleware.AuthUser(), handler.TestGetToken)
}

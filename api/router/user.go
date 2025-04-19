package router

import (
	"github.com/gin-gonic/gin"

	"github.com/kackerx/go-mall/api/handler"
	"github.com/kackerx/go-mall/common/middleware"
)

func registerUserRoutes(rg *gin.RouterGroup, userHandler *handler.UserHandler) {
	g := rg.Group("/user/")

	g.POST("register", userHandler.RegisterUser)
	g.POST("login", userHandler.LoginUser)
	g.GET("loginout", middleware.AuthUser(), userHandler.LoginoutUser)
}

package router

import (
	"github.com/gin-gonic/gin"

	"github.com/kackerx/go-mall/api/handler"
)

func RegisterUserRoutes(rg *gin.RouterGroup, userHandler *handler.UserHandler) {
	g := rg.Group("/user/")

	g.POST("register", userHandler.RegisterUser)
}

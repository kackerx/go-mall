package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kackerx/go-mall/api/handler"
	"github.com/kackerx/go-mall/api/router"
	"github.com/kackerx/go-mall/common/enum"
	"github.com/kackerx/go-mall/config"
	"github.com/kackerx/go-mall/dal/dao"
	"github.com/kackerx/go-mall/logic/appservice"
	"github.com/kackerx/go-mall/logic/domainservice"
)

func main() {
	conf := config.Conf
	if conf.App.Env == enum.ModeProd {
		gin.SetMode(gin.ReleaseMode)
	}

	e := gin.Default()

	baseHandler := handler.NewHandler()
	userDomainSvc := domainservice.NewUserDomainSvc()
	userAppSvc := appservice.NewUserAppSvc(userDomainSvc)
	userHandler := handler.NewUserHandler(baseHandler, userAppSvc)

	commodityDao := dao.NewCommodityDao()
	commodityDomainSvc := domainservice.NewCommoditySvc(commodityDao)
	commodityApp := appservice.NewCommodityApp(commodityDomainSvc)
	commodityHandler := handler.NewCommodityHandler(baseHandler, commodityApp)

	router.RegisterRoute(e, userHandler, commodityHandler)
	fmt.Println("listen on 9999")
	if err := http.ListenAndServe(":9999", e); err != nil {
		return
	}
}

package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kackerx/go-mall/api/router"
	"github.com/kackerx/go-mall/common/enum"
	"github.com/kackerx/go-mall/config"
)

func main() {
	conf := config.Conf
	if conf.App.Env == enum.ModeProd {
		gin.SetMode(gin.ReleaseMode)
	}

	e := gin.Default()
	router.RegisterRoute(e)

	fmt.Println("listen on 9999")
	if err := http.ListenAndServe(":9999", e); err != nil {
		return
	}
}

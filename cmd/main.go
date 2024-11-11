package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kackerx/go-mall/common/errcode"
	"github.com/kackerx/go-mall/common/log"
	"github.com/kackerx/go-mall/common/middleware"
	"github.com/kackerx/go-mall/config"
)

func main() {
	e := gin.Default()
	e.Use(middleware.StartTrace(), middleware.LogAccess(), middleware.GinPanicRecovery())

	e.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"hello": "world",
		})
	})

	e.GET("/conf", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"dsn": config.Conf.DB.Dsn,
		})
	})

	e.GET("/log", func(c *gin.Context) {
		log.New(c).Info("log test", "key", "keyValue", "val", 2)
		c.JSON(http.StatusOK, gin.H{
			"max_life_time": config.Conf.DB.MaxLiftTime,
		})
	})

	e.POST("/access", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"max_idle": config.Conf.DB.MaxIdle,
		})
	})

	e.GET("/panic", func(c *gin.Context) {
		var m map[int]int
		m[1] = 1
		c.JSON(http.StatusOK, gin.H{"status": "ok", "data": m})
	})

	e.GET("/customerr", func(c *gin.Context) {
		// 使用Wrap包装错误, 生成项目错误
		err := errors.New("dao error")
		appErr := errcode.Wrap("包装错误", err)
		bAppErr := errcode.Wrap("再包装错误", appErr)
		log.New(c).Error("记录错误", "err", bAppErr)

		// 预定义的错误, 追加错误根因
		err = errors.New("domain err")
		apiErr := errcode.ErrServer.WithCause(err)
		log.New(c).Error("API错误", "err", apiErr)
		c.JSON(apiErr.HttpStatusCode(), gin.H{
			"code": apiErr.Code(),
			"msg":  apiErr.Msg(),
		})
	})

	fmt.Println("listen on 9999")
	if err := http.ListenAndServe(":9999", e); err != nil {
		return
	}
}

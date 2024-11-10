package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kackerx/go-mall/common/log"
	"github.com/kackerx/go-mall/config"
)

func main() {
	e := gin.Default()

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
		for i := 0; i < 100; i++ {
			log.ZapTest()
		}
		c.JSON(http.StatusOK, gin.H{
			"max_life_time": config.Conf.DB.MaxLiftTime,
		})
	})

	fmt.Println("listen on 9999")
	if err := http.ListenAndServe(":9999", e); err != nil {
		return
	}
}

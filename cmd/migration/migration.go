package main

import (
	"fmt"

	"github.com/kackerx/go-mall/dal/dao"
	"github.com/kackerx/go-mall/dal/model"
)

func main() {
	if err := dao.DB().Migrator().AutoMigrate(&model.DemoOrder{}); err != nil {
		panic(err)
	}

	fmt.Println("Migrator End")
}

package main

import (
	"fmt"

	"github.com/kackerx/go-mall/common/design_patterns/mediator"
)

func main() {
	fmt.Println("=== 中介者模式示例 - 聊天室系统 ===")
	fmt.Println()

	// 运行示例
	mediator.RunExample()

	fmt.Println()
	fmt.Println("=== 示例结束 ===")
}

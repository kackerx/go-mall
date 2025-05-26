package main

import (
	"fmt"
	"iter"
)

// 演示迭代器中变量的类型和来源
func main() {
	users := []*User{
		{Name: "Alice"},
		{Name: "Bob"},
		{Name: "Charlie"},
	}

	for user := range UserIter(users) {
		fmt.Println(user.Name)
		if user.Name == "Bob" {
			fmt.Println("main end")
			break
		}
	}
}

type User struct {
	Name string
}

func UserIter(us []*User) iter.Seq[*User] {
	return func(yield func(*User) bool) {
		for _, user := range us {
			if !yield(user) {
				fmt.Println("yield end", user.Name)
				break
			}
		}
	}
}

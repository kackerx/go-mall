package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestFoo(t *testing.T) {
	// m := map[string]int{}
	// hehe(m)

	// fmt.Println(m)
	// b := []int{1}
	// fmt.Println(b[:2])

	fmt.Println(bar(1))
}

func hehe(m map[string]int) {
	m["a"] = 1
}

func TestGetName(t *testing.T) {
	tests := []struct {
		Name    string
		Aint    int
		wanterr error
	}{
		{
			Name:    "hehe",
			Aint:    1,
			wanterr: nil,
		},
		{
			Name:    "hehe",
			Aint:    1,
			wanterr: nil,
		},
		{
			Name: "que",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {

		})
	}
}

func bar(i int) (t int) {
	t = i
	return 2
}

type User struct {
	Name string `json:"name"`
	Age  int
}

func NewUser(name string, age int) *User {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(10)

	ch := make(chan struct{}, 1)
	ch <- struct{}{}

	select {
	case <-ctx.Done():
		fmt.Println("timeout")

	}

	return &User{Name: name}
}

func (u *User) Build(ctx context.Context) *http.Client {
	user := NewUser("kkk", 1)
	if user.Name == "kkk" {

	}
	return http.DefaultClient
}

func GetUserByID(id uint) *User {
	return &User{
		Name: "",
		Age:  0,
	}
}

package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kackerx/go-mall/common/util"
)

func TestFoo(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	go foo(ctx)

	select {
	case <-ctx.Done():
		fmt.Println("parent done")
	case <-time.After(time.Second):
		fmt.Println("one sec")
		cancel()

	case <-time.After(time.Second * 3):
		fmt.Println("timeout")
	}

	time.Sleep(time.Second * 5)
	fmt.Println("end")
}

func foo(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	select {
	case <-ctx.Done():
		fmt.Println("son cancel")
	}
}

func TestUtil(t *testing.T) {
	// fmt.Println(util.MaskEmail("kingvstr@hotmail.com"))
	fmt.Println(util.MaskRealName("王凯是旋"))
}

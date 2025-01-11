package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/go-cmd/cmd"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
)

const (
	HistorySet  = "tool:video:history"
	PendingList = "tool:video:pending"
)

var (
	rdb *redis.Client
	ctx = context.Background()
	// 命令行参数
	inURL   string
	downURL string
	out     bool
)

func main() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	var rootCmd = &cobra.Command{
		Use:   "tool",
		Short: "工具集合",
	}

	var videoCmd = &cobra.Command{
		Use:   "video",
		Short: "视频URL处理工具",
		Args:  cobra.MaximumNArgs(1), // 允许最多一个位置参数
		Run: func(cmd *cobra.Command, args []string) {
			choice := showOptionsAndGetChoice()
			fmt.Printf("你选择了: %s\n", choice)
			// 处理输入URL
			if inURL != "" {
				handleVideoIn(inURL)
				return
			}
			// 处理输出
			if out {
				handleVideoOut()
				return
			}
			if downURL != "" {
				name := args[0]
				handleVideoDown(downURL, name)
				return
			}
			// 如果既没有--in也没有--out，显示使用帮助
			cmd.Help()
		},
	}

	// 添加命令行参数

	videoCmd.Flags().StringVarP(&downURL, "down", "d", "", "下载视频URL")
	videoCmd.Flags().StringVarP(&inURL, "in", "i", "", "添加视频URL")
	videoCmd.Flags().BoolVarP(&out, "out", "o", false, "获取一个待处理的URL")

	rootCmd.AddCommand(videoCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func showOptionsAndGetChoice() string {
	options := []string{"Option 1", "Option 2", "Option 3", "Option 4"}
	fmt.Println("请选择以下选项：")
	for i, option := range options {
		fmt.Printf("%d. %s\n", i+1, option)
	}
	var choice int
	for {
		fmt.Print("输入你的选择 (1-4): ")
		_, err := fmt.Scan(&choice)
		if err == nil && choice >= 1 && choice <= len(options) {
			break
		}
		fmt.Println("无效的选择，请重新输入。")
	}
	return options[choice-1]
}

func handleVideoDown(url, name string) {
	exist, err := rdb.SIsMember(ctx, HistorySet, url).Result()
	if err != nil {
		panic("检查URL时出错:" + err.Error())
	}

	if exist {
		fmt.Println("url已经存在")
		return
	}

	if name == "" {
		name = uuid.NewString()
	}
	// name += "_" + time.Now().Format("2006-01-02T15:04:05")
	name = "/Users/apple/Downloads/" + name + time.Now().Format("2006-01-02-15_04_05") + ".mp4"
	fmt.Printf("下载视频: %s, ------ %s", name, url)
	// 使用ffmpeg -i https://xxx.wujiangwl.xyz/uploads/m3u8/%E4%B8%AD%E6%96%87%E5%AD%97%E5%B9%95/20230701/b1fa1c16242dc0aff52a7d1ddb827e52.m3u8 -c copy f1.mp4

	// 创建命令
	// ffmpeg := cmd.NewCmd("/opt/homebrew/bin/ffmpeg", "-i", url, "-c", "copy", name)
	ffmpeg := cmd.NewCmdOptions(cmd.Options{Streaming: true}, "/opt/homebrew/bin/ffmpeg", "-i", url, "-c", "copy", name)
	// ffmpeg.Start()
	// 打印命令详情
	// fmt.Printf("命令路径: %s\n", ffmpeg.Name)
	// fmt.Printf("命令参数: %v\n", ffmpeg.Args)

	// 监听标准输出
	var wg sync.WaitGroup
	wg.Add(2)
	done := make(chan struct{})

	go func() {
		defer wg.Done()
		for line := range ffmpeg.Stderr {
			fmt.Printf("\rERR: %s\n", line)
		}
		close(done)
	}()
	go func() {
		defer wg.Done()
		for {
			select {
			case line := <-ffmpeg.Start():
				fmt.Printf("\r状态: %s\n", line)
			case <-done:
				fmt.Println("doneeeeeeeeeeeee")
				return
			}
		}
		// status := <-ffmpeg.Start()
		// fmt.Printf("状态: %s\n", status)
	}()
	wg.Wait()

	fmt.Printf("下载完成: %s, ------ %s\n", name, url)
	rdb.SAdd(ctx, HistorySet, url)
}

func handleVideoIn(url string) {
	exists, err := rdb.SIsMember(ctx, HistorySet, url).Result()
	if err != nil {
		fmt.Println("检查URL时出错:", err)
		return
	}

	if exists {
		fmt.Println("URL已经存在:", url)
		return
	}

	pipe := rdb.Pipeline()
	pipe.SAdd(ctx, HistorySet, url)
	pipe.RPush(ctx, PendingList, url)
	_, err = pipe.Exec(ctx)

	if err != nil {
		panic("添加URL时出错:" + err.Error())
	}
	fmt.Println("URL已添加:", url)
}

func handleVideoOut() {
	url, err := rdb.LPop(ctx, PendingList).Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("没有待处理的URL")
		} else {
			fmt.Println("获取URL时出错:", err)
		}
		return
	}
	fmt.Println("获取到URL:", url)
}

package main

import (
	"Zinx/znet"
)

func main() {
	//1.创建server句柄，使用Zinx的api

	s := znet.NewServer("[zinx V0.1]")
	//2.启动server
	s.Serve()
}

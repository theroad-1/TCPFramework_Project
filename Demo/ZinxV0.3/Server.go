package main

import (
	"Zinx/ziface"
	"Zinx/znet"
	"fmt"
)

type PingRouter struct {
	znet.BaseRouter
}

func (r *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping.."))
	if err != nil {
		fmt.Println("call back before ping error")
	}
}

//在处理conn业务的主方法hook
func (r *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle..")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping.."))
	if err != nil {
		fmt.Println("call back ping error")
	}
}

//在处理conn业务之后的钩子方法hook
func (r *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Route PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping"))
	if err != nil {
		fmt.Println("call back after ping error")
	}
}
func main() {

	//1.创建server句柄，使用Zinx的api

	s := znet.NewServer("[zinx V0.3]")
	//2.给当前zinx框架添加一个自定义的router
	s.AddRouter(&PingRouter{})
	//2.启动server
	s.Serve()
}

package main

import (
	"Zinx/ziface"
	"Zinx/znet"
	"fmt"
)

type PingRouter struct {
	znet.BaseRouter
}

//在处理conn业务的主方法hook
func (r *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle..")
	//先读取客户端数据，再回写ping
	fmt.Println("recv from client :msgID =", request.GetMsgID(),
		",data =", string(request.GetData()))
	err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping..."))
	if err != nil {
		fmt.Println(err)
	}
}

//hello zinx test自定义路由
type HelloZinxRouter struct {
	znet.BaseRouter
}

//在处理conn业务的主方法hook
func (this *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle..")
	//先读取客户端数据，再回写ping
	fmt.Println("recv from client :msgID =", request.GetMsgID(),
		",data =", string(request.GetData()))
	err := request.GetConnection().SendMsg(201, []byte("Hello welcome to zinx"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {

	//1.创建server句柄，使用Zinx的api
	s := znet.NewServer()
	//2.给当前zinx框架添加一个自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})
	//2.启动server
	s.Serve()
}

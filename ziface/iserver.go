//抽象层
package ziface

//定义一个服务器接口
type IServer interface {
	//启动
	Start()
	//停止
	Stop()
	//运行
	Serve()

	//路由功能：给当前的服务注册一个路由方法，供客户端的连接处理使用
	AddRouter(router IRouter)
}

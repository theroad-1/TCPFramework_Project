package znet

import (
	"Zinx/ziface"
	"fmt"
	"net"
)

//IServer的接口实现,定义一个Server的服务器模块
type Server struct {
	//服务器的名称
	Name string
	//服务器绑定的ip版本
	IPVersion string
	//服务器监听的ip
	IP string
	//服务器监听的端口
	Port int
	//当前的Server添加一个router，server注册的连接对应的处理业务
	Router ziface.IRouter
}

func (s *Server) Start() {
	fmt.Printf("[Start]Server Listenner at IP : %s,Port:%d, is starting\n", s.IP, s.Port)

	go func() {
		//1.获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp add error:", err)
			return
		}
		//2.监听服务器的地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, "err:", err)
			return
		}
		fmt.Println("start Zinx server success...", s.Name, "success, Listenning...")
		var cid uint32
		cid = 0
		//3.阻塞的等待客户端链接，处理客户端链接业务
		for {

			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("accept err", err)
				continue
			}
			//将处理新连接的业务方法和conn进行绑定，得到我们的连接模块
			dealConn := NewConnection(conn, cid, s.Router)
			cid++
			go dealConn.Start()
		}
	}()

}
func (s *Server) Stop() {

}
func (s *Server) Serve() {
	s.Start()
	//做一些启动服务器之后的额外业务
	select {}
}

//路由功能：给当前的服务注册一个路由方法，供客户端的连接处理使用
func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("Add Router success!!")
}

//初始化Serve模块的方法
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
		Router:    nil,
	}
	return s
}

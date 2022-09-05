package znet

import (
	"Zinx/utils"
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
	//当前server的消息管理模块，用来绑定MsgID和对应的处理业务API关系
	MsgHandler ziface.IMsgHandle
	//该server的连接管理器
	ConnMgr ziface.IConnManager
	//该server创建连接之后自动调用hook函数
	OnConnStart func(conn ziface.IConnection)
	//该server销毁连接之前启动调用hook函数
	OnConnStop func(conn ziface.IConnection)
}

func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name : %s,listenner at IP : %s,Port %d is starting\r",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version %s,MaxConn:%d,MaxPackageSize:%d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)
	//fmt.Printf("[Start]Server Listenner at IP : %s,Port:%d, is starting\n", s.IP, s.Port)

	go func() {
		//0 开启消息队列及worker工作池
		s.MsgHandler.StartWorkerPool()

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

			//连接管理模块：设置最大连接个数的判断，如果超过最大连接，那么关闭此新的连接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				//TODO 给客户端响应一个超出最大连接的错误包
				fmt.Println("=====>Too many Connections MaxConn = ", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}

			//将处理新连接的业务方法和conn进行绑定，得到我们的连接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			go dealConn.Start()
		}
	}()

}
func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server name ", s.Name)
	s.ConnMgr.ClearConn()
}
func (s *Server) Serve() {
	s.Start()
	//做一些启动服务器之后的额外业务

	//阻塞
	select {}
}

//路由功能：给当前的服务注册一个路由方法，供客户端的连接处理使用
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router success!!")
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

//初始化Serve模块的方法
func NewServer() ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

//注册OnConnStart钩子函数的方法
func (s *Server) SetOnConnStart(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

//注册OnConnStop钩子函数的方法
func (s *Server) SetOnConnStop(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

//调用OnConnStart钩子函数的方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("in znet.server --->Call OnCallStart()...")
		s.OnConnStart(conn)
	}
}

//调用OnConnStop钩子函数的方法
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("in znet.server ---->Call OnCallStop()...")
		s.OnConnStop(conn)
	}
}

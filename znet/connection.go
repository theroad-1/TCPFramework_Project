package znet

import (
	"Zinx/ziface"
	"fmt"

	"net"
)

/*
链接模块
*/
type Connection struct {
	//当前链接的socket：TCP套接字
	Conn *net.TCPConn
	//链接的ID
	ConnID uint32
	//当前的链接状态
	isClosed bool

	//告知当前路径已经退出/停止 的channel
	ExitChan chan bool
	//该链接处理的方法Router
	Router ziface.IRouter
}

//初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		Router:   router,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}
	return c
}
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("ConnID = ", c.ConnID, "Reader is exit,remote addr=", c.RemoteAddr().String())
	defer c.Stop()
	for {
		//读取客户端的数据到buf中，最大512字节
		buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("receive buf err :", err)
			continue
		}

		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			data: buf,
		}
		//从路由找到注册绑定的Conn对应的router调用
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)

		}(&req)

	}
}

//启动链接 让当前的链接准备开始工作
func (c *Connection) Start() {
	fmt.Println("conn start()..ConnID=", c.ConnID)
	//启动从当前链接的读数据业务
	go c.StartReader()
}

//停止链接 结束当前链接的工作
func (c *Connection) Stop() {
	fmt.Println("conn stop()..ConnID = ", c.ConnID)
	//判断是否已经关闭
	if c.isClosed == true {
		return
	}

	c.isClosed = true
	//关闭socket链接
	c.Conn.Close()
	//回收资源
	close(c.ExitChan)
}

//获取当前链接的绑定socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前连接模块的链接id
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端的tcp状态 IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//发送数据，将数据发送给远程的客户端
func (c *Connection) Send(data []byte) error {
	return nil
}

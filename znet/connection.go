package znet

import (
	"Zinx/utils"
	"Zinx/ziface"
	"errors"
	"fmt"
	"io"

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

	//无缓冲的管道，用于读写goroutine之间的消息通道
	msgChan chan []byte

	//消息的管理msgID和对应的处理业务API关系
	MsgHandler ziface.IMsgHandle
}

//初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		isClosed:   false,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
	}
	return c
}
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("[Reader is exit]ConnID = ", c.ConnID, "remote addr=", c.RemoteAddr().String())
	defer c.Stop()
	for {
		////读取客户端的数据到buf中，最大512字节
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("receive buf err :", err)
		//	continue
		//}

		//创建一个拆包解包的对象
		dp := NewDataPack()

		//读取客户端的Msg Head 二进制流8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head err:", err)
			break
		}

		//拆包，得到msgID和msgDatalen放在msg信息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack err", err)
			break
		}
		//根据dataLen 再次读取Data，放在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error:", err)
				break
			}
		}
		msg.SetData(data)

		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			//已经开启了工作池机制，将信息发送给worker工作池处理即可
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			//根据绑定好的MsgID 找到对应的处理api业务 执行
			go c.MsgHandler.DoMsgHandler(&req)
		}

	}
}

/*
写消息goroutine 专门发送给客户端的模块
*/
func (c *Connection) StartWriter() {
	fmt.Println("[Writer goroutine is running]")
	defer fmt.Println("[conn writer exit]", c.RemoteAddr().String())
	//不断阻塞等待channel消息，写给客户端
	for true {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("send data err:", err)
				return
			}
		case <-c.ExitChan:
			return

		}
	}
}

//启动链接 让当前的链接准备开始工作
func (c *Connection) Start() {
	fmt.Println("conn start()..ConnID=", c.ConnID)
	//启动从当前链接的读数据业务
	go c.StartReader()
	go c.StartWriter()
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
	c.ExitChan <- true
	//回收资源
	close(c.ExitChan)
	close(c.msgChan)
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

//提供一个sendMsg方法 将我们要发送给客户端的数据，先进行封包，再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection closed when send msg")
	}
	//将data进行封包 MsgDataLen MsgID Data
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack err msg id :", msgId)
		return errors.New("Pack error msg")
	}

	//将数据发送给客户端
	c.msgChan <- binaryMsg

	return nil
}

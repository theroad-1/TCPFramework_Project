package utils

import (
	"Zinx/ziface"
	"encoding/json"
	"io/ioutil"
)

/*
储存一切有关Zinx框架的全局参数，供其他模块使用
一些参数可以通过zinx.json由用户进行配置

*/
type GlobalObj struct {
	/*
		Server
	*/
	TcpServer ziface.IServer //当前Zinx全局的Server对象
	Host      string         //服务器主机监听的IP
	TcpPort   int            //服务器主机监听的端口号
	Name      string         //服务器名称

	/*
		Zinx
	*/
	Version          string //Zinx版本号
	MaxConn          int    //服务器主机允许的最大连接数
	MaxPackageSize   uint32 //当前Zinx框架数据包的最大值
	WorkerPoolSize   uint32 //当前业务工作worker池的goroutine数量
	MaxWorkerTaskLen uint32 //zinx框架允许用户最多开辟多少个worker（限定条件）
}

/*
定义一个全局的对外Global对象
*/
var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {

	data, err := ioutil.ReadFile("E:\\GoProject\\Zinx_Project\\MMO_Game_Zinx\\conf\\zinx.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, GlobalObject)
	if err != nil {
		panic(err)
	}
}

/*
提供一个init方法，初始化globalobj
*/
func init() {
	//先配置使用默认的值
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "V0.9",
		TcpPort:          8999,
		Host:             "0.0.0.0", //当前主机所有网卡的地址，单个网卡可以写127.0.0.1
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,   //worker工作池的队列的个数
		MaxWorkerTaskLen: 1024, //每个worker对应的消息队列的任务的数量最大值

	}

	//从conf/zinx.json去加载一些用户自定义的参数
	GlobalObject.Reload()
}

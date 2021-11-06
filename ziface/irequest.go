package ziface

/*
irequest接口
实际时是把客户端请求的链接信息和和请求数据，包装到一个Request中
*/
type IRequest interface {
	//得到当前链接
	GetConnection() IConnection
	//得到请求的消息数据
	GetData() []byte
	//
}

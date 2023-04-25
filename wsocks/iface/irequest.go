package iface

/*
	IRequest 接口：
	实际上是把客户端请求的链接信息 和 请求的数据 包装到了 Request 里
*/

type IRequest interface {
	//获取请求连接信息
	GetConnection() IConnection

	// 获取请求消息数据
	GetData() []byte
}

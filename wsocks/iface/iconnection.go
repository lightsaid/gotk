package iface

import "net"

// 定义链接接口
type IConnection interface {
	// 启动链接，当 accept 后，将 conn 交给 Connection struct 处理业务
	Start()

	// 停止/关闭链接, 结束当前链接
	Stop()

	// 获取原始的tcp链接 (net.Conn)
	GetTCPConnection() *net.TCPConn

	// 获取链接ID
	GetConnID() uint32
}

// 定义一个统一处理链接业务的接口
// cc 原始 tcp 链接对象
// data 客户端请求数据
// length 数据长度
type HandleConnFunc func(cc *net.TCPConn, data []byte, length int) error

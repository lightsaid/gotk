package wsocks

import "github.com/lightsaid/gotk/wsocks/iface"

// Request 客户端请求对象封装结构(请求体)
type Request struct {
	// 已经和客户端建立好的链接
	conn iface.IConnection
	// 客户端请求的数据
	data []byte
}

// GetConnection 获取请求连接信息
func (r *Request) GetConnection() iface.IConnection {
	return r.conn
}

// GetData 获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.data
}

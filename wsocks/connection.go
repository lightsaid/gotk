package wsocks

import (
	"net"

	"github.com/lightsaid/gotk/wsocks/iface"
)

type Connection struct {
	// 链接对象
	Conn *net.TCPConn

	// 链接ID
	ConnID uint32

	// 当前链接是否关闭
	IsClosed bool

	// 处理该链接handler
	handleFunc iface.HandleConnFunc

	// 告知该链接已经退出/停止channel
	Exit chan bool
}

func NewConnection(conn *net.TCPConn, connID uint32, handle iface.HandleConnFunc) iface.IConnection {
	return &Connection{
		Conn:       conn,
		ConnID:     connID,
		handleFunc: handle,
		IsClosed:   false,
		Exit:       make(chan bool, 1),
	}
}

// 启动链接，当 accept 后，将 conn 交给 Connection struct 处理业务
func (c *Connection) Start() {

}

// 停止/关闭链接, 结束当前链接
func (c *Connection) Stop() {

}

// 获取原始的tcp链接 (net.Conn)
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取链接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

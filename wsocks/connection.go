package wsocks

import (
	"fmt"
	"net"

	"github.com/lightsaid/gotk/wsocks/iface"
)

// Connection 链接对象结构
type Connection struct {
	// 链接对象
	Conn *net.TCPConn

	// 链接ID
	ConnID uint32

	// 当前链接是否关闭
	IsClosed bool

	// 处理该链接handler
	// handleFunc iface.HandleConnFunc  交由 Router 处理业务

	// 该连接的处理方法router
	Router iface.IRouter

	// 告知该链接已经退出/停止channel
	Exit chan struct{}
}

// NewConnection 创建一个链接对象
func NewConnection(conn *net.TCPConn, connID uint32, router iface.IRouter) iface.IConnection {
	return &Connection{
		Conn:   conn,
		ConnID: connID,
		// handleFunc: handle,
		Router:   router,
		IsClosed: false,
		Exit:     make(chan struct{}),
	}
}

// Start 启动链接，当 accept 后，将 conn 交给 Connection struct 处理业务
func (c *Connection) Start() {
	// 开启处理该链接读取到客户端数据之后的请求业务
	go c.StartReader()

	// 阻塞，直到接受到退出通知
	<-c.Exit

	// for {
	// 	select {
	// 	case <-c.Exit:
	// 		接受到退出通知
	// 		return
	// 	}
	// }
}

// Stop 停止/关闭链接, 结束当前链接
func (c *Connection) Stop() {
	if c.IsClosed {
		return
	}

	//TODO Connection Stop() 如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用

	c.IsClosed = true

	// 关闭socket链接
	c.Conn.Close()

	//通知从缓冲队列读数据的业务，该链接已经关闭
	// c.Exit <- struct{}{}

	close(c.Exit)

}

// GetTCPConnection 获取原始的tcp链接 (net.Conn)
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取链接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// StartReader 启动读模块，应该使用一个 goroutine 启动该函数
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running.")
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit.")
	defer c.Stop()

	for {
		var buf = make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			c.Exit <- struct{}{}
			continue
		}

		// // 调用当前链接绑定处理函数
		// if err := c.handleFunc(c.Conn, buf, n); err != nil {
		// 	fmt.Printf("[%d] handle error %v\n", c.ConnID, err)
		// 	c.Exit <- struct{}{}
		// 	return
		// }

		req := Request{
			conn: c,
			data: buf,
		}

		go func(req iface.IRequest) {
			// 执行注册路由方法
			c.Router.PreHandle(req)
			c.Router.Handle(req)
			c.Router.PostHandle(req)
		}(&req)
	}
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

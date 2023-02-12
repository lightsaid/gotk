package usocket

import (
	"fmt"
	"io"
	"net"
)

type IConnection interface {
	// 启动链接，让当前链接准备开始工作
	Start()

	// 停止链接，结束当前链接的工作
	Stop()

	// 获取当前链接绑定的socket conn
	GetTCPConnection() *net.TCPConn

	// 获取当前链接的ID
	GetConnID() uint32

	// 获取远程客户端的 TCP IP Port
	RemoteAddr() net.Addr

	// 发送数据
	Send(data []byte) error
}

// 定义一个处理链接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error

// Connection 链接模块
type Connection struct {
	// 当前链接的socket TCP 套接字
	Conn      *net.TCPConn
	ConnID    uint32
	isClosed  bool
	handleAPI HandleFunc
	ExistChan chan bool
}

func NewConnection(conn *net.TCPConn, ConnID uint32, callback_api HandleFunc) IConnection {
	c := &Connection{
		Conn:      conn,
		ConnID:    ConnID,
		handleAPI: callback_api,
		isClosed:  false,
		ExistChan: make(chan bool, 1),
	}

	return c
}

// 链接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running.")
	defer fmt.Println("connID=", c.ConnID, " Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	var count = 0

	for {
		// 读取客户端的数据到buf中，最大1024字节
		buf := make([]byte, 1024)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			if cnt == 0 && err == io.EOF {
				break
			}
			fmt.Println("recv buf error: ", err.Error())
			if count > 5 {
				break
			}
			count++
			continue
		}

		// 调用当前链接所绑定的HandlerAPI
		if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Println("Callback error: ", err)
			break
		}
	}
}

// 启动链接，让当前链接准备开始工作
func (c *Connection) Start() {
	// 启动当前链接的读数据的业务
	go c.StartReader()

	// TODO: 启动当前链接的写数据的业务

}

// 停止链接 结束当前链接工作
func (c *Connection) Stop() {
	fmt.Println("conn stop: ", c.ConnID)
	if c.isClosed {
		return
	}

	c.isClosed = true

	// 回收资源
	c.Conn.Close()
	close(c.ExistChan)
}

// 获取当前链接绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前链接的ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取远程客户端的 TCP IP Port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据
func (c *Connection) Send(data []byte) error {
	return nil
}

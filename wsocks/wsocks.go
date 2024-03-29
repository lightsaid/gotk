package wsocks

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/lightsaid/gotk/wsocks/iface"
)

type Server struct {
	// 服务器名称
	Name string
	// ip 版本, 对应的是 net.ListenTCP(network string, laddr *TCPAddr) network: "tcp", "tcp4", "tcp6"
	IPVersion string
	// 服务器主机
	Host string
	// 服务端口
	Port int

	// 当前Server由用户绑定的回调router,也就是Server注册的链接对应的业务处理
	Router iface.IRouter
}

func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	//回显业务
	fmt.Println("[Conn Handle] CallBackToClient ... ")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err ", err)
		return errors.New("CallBackToClient error")
	}
	return nil
}

func NewServer(name string) iface.IServer {
	srv := &Server{
		Name:      name,
		IPVersion: "tcp4",
		Host:      "0.0.0.0",
		Port:      9000,
		Router:    nil,
	}

	return srv
}

func (s *Server) Start() {
	log.Printf("Start %s websocket server...", s.Name)

	go func() {
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.Host, s.Port))
		if err != nil {
			log.Fatal(err)
		}

		lis, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			log.Fatal(err)
		}

		var cid uint32 = 1

		for {
			conn, err := lis.AcceptTCP()
			if err != nil {
				log.Println("Accept error ", err)
				continue
			}

			connection := NewConnection(conn, cid, s.Router)
			cid++

			go connection.Start()

			// go func() {
			// 	// 不断从循环从客户端读取数据
			// 	for {
			// 		// 暂定 512 字节
			// 		buf := make([]byte, 512)
			// 		n, err := conn.Read(buf)
			// 		if err != nil {
			// 			log.Println("read data error ", err)
			// 			continue
			// 		}
			// 		// 读取多少字节，就回响多少字节
			// 		if _, err := conn.Write(buf[:n]); err != nil {
			// 			log.Println("write data error ", err)
			// 			continue
			// 		}
			// 	}
			// }()
		}

	}()

}
func (s *Server) Stop() {
	log.Printf("Stop %s server.", s.Name)
}
func (s *Server) Serve() {
	s.Start()

	fmt.Println("键入任意字符停止....")
	var input string
	fmt.Scanln(&input)
}

// AddRouter 路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (s *Server) AddRouter(router iface.IRouter) {
	s.Router = router

	fmt.Println("Add Router succ! ")
}

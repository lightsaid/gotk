package usocket

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/lightsaid/gotk/random"
)

// ServerOpts 默认值
const DefaultServerOptsIP = "0.0.0.0"
const DefaultServerOptsPort = 5205
const DefaultServerNetwork = "tcp"

// IServer 定义 TCP 服务接口
type IServer interface {
	// 启动服务
	Start()

	// 停止服务
	Stop()

	// 运行服务
	Serve()
}

// ServerOpts TCP 服务基础参数/选项
type ServerOpts struct {
	// 服务器名称
	Name string

	// NetWork "tcp", "tcp4", "tcp6"...
	Network string

	// ip 地址
	IP string

	// 端口
	Port int

	// 监听IP地址 IP:Port
	listenAddr string
}

// Server 服务结构体，实现服务接口（IServer）
type Server struct {
	ServerOpts
}

// NewServer 实例化一个TCP服务
func NewServer(opts ServerOpts) IServer {
	opts = defaultServerOpts(opts)
	srv := &Server{
		opts,
	}

	return srv
}

func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	log.Println("CallBack~")
	if _, err := conn.Write(data); err != nil {
		return err
	}
	return nil
}

// Start 启动服务
func (s *Server) Start() {
	go func() {
		// 获取一个TCP地址
		addr, err := net.ResolveTCPAddr(s.Network, s.listenAddr)
		if err != nil {
			log.Println("resolve tcp addr error: ", err)
			return
		}

		// 监听 TCP 服务
		ln, err := net.ListenTCP(s.Network, addr)
		if err != nil {
			log.Println("listen tcp error: ", err)
			return
		}

		log.Printf("Start %s success. listen addr on %s.", s.Name, s.listenAddr)

		for {
			conn, err := ln.AcceptTCP()
			if err != nil {
				log.Println("Accept tcp error: ", err)
				continue
			}

			var cid uint32 = 100
			dealConn := NewConnection(conn, cid, CallBackToClient)
			cid++

			// 启动当前Start
			go dealConn.Start()

			// go func() {
			// 	defer conn.Close()
			// 	for {
			// 		buf := make([]byte, 1024)
			// 		// 阻塞，等待读取数据；如果client close关闭，则返回 err = io.EOF
			// 		n, err := conn.Read(buf)
			// 		fmt.Println("Reading...")
			// 		if err != nil {
			// 			if err == io.EOF {
			// 				log.Println("client close...")
			// 				break
			// 			}
			// 			log.Println("conn read error: ", err)
			// 			continue
			// 		}

			// 		message := fmt.Sprintf("%s: %s", s.Name, string(buf[:n]))
			// 		if _, err := conn.Write([]byte(message)); err != nil {
			// 			log.Println("conn write error: ", err)
			// 			continue
			// 		}
			// 	}
			// }()
		}
	}()
}

// Stop 停止服务
func (s *Server) Stop() {}

// Serve 运行服务
func (s *Server) Serve() {
	s.Start()

	select {}
}

// defaultServerOpts 对零值的服务选项初始化默认值
func defaultServerOpts(opts ServerOpts) ServerOpts {
	if opts.Name == "" {
		opts.Name = "TCP_SERVER_" + strings.ToUpper(random.RandomString(6))
	}

	if opts.IP == "" {
		opts.IP = DefaultServerOptsIP
	}

	if opts.Port == 0 {
		opts.Port = DefaultServerOptsPort
	}

	if opts.Network == "" {
		opts.Network = DefaultServerNetwork
	}

	opts.listenAddr = fmt.Sprintf("%s:%d", opts.IP, opts.Port)

	return opts
}

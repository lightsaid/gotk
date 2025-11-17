package utils

import (
	"encoding/json"
	"os"

	"github.com/lightsaid/gotk/wsocks/iface"
)

var EnvConfig *Config

type Config struct {
	TCPServer     iface.IServer `json:"-"`               // 当前wsocks全局server对象
	Host          string        `json:"host"`            // 主机
	TCPPort       int           `json:"tcp_port"`        // 端口
	Name          string        `json:"name"`            // 服务名称
	Version       string        `json:"version"`         // 当前版本
	MaxPacketSize uint32        `json:"max_packet_size"` // 数据包最大值
	MaxConns      int           `json:"max_conns"`       // 允许最大链接数
}

func (c *Config) Reload() {
	data, err := os.ReadFile("configs/server.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &EnvConfig)
	if err != nil {
		panic(err)
	}
}

func init() {
	EnvConfig = &Config{
		Name:          "WsocksServerApp",
		Version:       "V0.4",
		TCPPort:       9999,
		Host:          "127.0.0.1",
		MaxConns:      10000,
		MaxPacketSize: 2048,
	}

	EnvConfig.Reload()
}

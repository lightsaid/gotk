package iface

// 这种接口与实现分离模式，分工明确，逻辑清晰，方便测试

// 定义服务器接口
type IServer interface {
	// 启动服务
	Start()

	// 停止服务
	Stop()

	// 开启业务服务方法
	Serve()
}

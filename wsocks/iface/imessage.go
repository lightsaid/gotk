package iface

/*
将请求的一个消息封装到message中，定义抽象层接口
*/
type IMessage interface {
	GetDataLen() uint32 // 获取消息数据段长度
	GetMsgId() uint32   // 获取消息
	GetData() []byte    // 获取消息内容

	SetMsgId(uint32)   // 设计消息ID
	SetData([]byte)    // 设计消息内容
	SetDataLen(uint32) // 设置消息段长度
}

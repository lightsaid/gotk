package wsocks

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/lightsaid/gotk/wsocks/iface"
	"github.com/lightsaid/gotk/wsocks/utils"
)

// DataPack 封装包类实例
type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

// GetHeadLen 获取包头长度方法
func (db *DataPack) GetHeadLen() uint32 {
	// 消息Id（uint32） + 消息长度(uint32) = uint32(4字节) + uint32(4字节) = 8
	return 8
}

// Pack 封包
func (dp *DataPack) Pack(msg iface.IMessage) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})
	// 解决tcp粘包问题：
	// 消息由 Heade + Body 组成，其中 Head 包括 dataLen + msgID，Body 仅有数据data
	// 写入顺序：dataLen -> msgId -> data

	// 写入 dataLen
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	// 写入msgId
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	// 写入data
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

func (db *DataPack) Unpack(binaryData []byte) (iface.IMessage, error) {
	// 创建一个输入二进制的 ioReader
	dataBuff := bytes.NewReader(binaryData)

	// 解压head部分，获取dataLen和msgId
	msg := &Message{}
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	// 读取msgId
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	if utils.EnvConfig.MaxPacketSize > 0 && msg.DataLen > utils.EnvConfig.MaxPacketSize {
		return nil, errors.New("too large msg data recieved")
	}

	// 这里只需要把head的数据拆包出来就可以了，然后再通过head的长度，再从conn读取一次数据
	return msg, nil
}

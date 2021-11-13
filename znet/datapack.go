package znet

import (
	"Zinx/utils"
	"Zinx/ziface"
	"bytes"
	"encoding/binary"
	"errors"
)

//封包、拆包的具体模块
type DataPack struct{}

//封包、拆包实例的一个初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
	//返回头部的8字节，Datalen、Id
	return 8
}

func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	//创建一个存放byte字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//将dataLen写进databuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	//将MsgId写进databuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	//将data数据写进databuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

//拆包方法 （将包的Head信息读出来） 之后再根据head信息里的data的长度，再进行一次读
func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	//创建一个从二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	//只解压head信息，得到datalen和id信息
	msg := &Message{}

	//读datalen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	//读msgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	//判断datalen是否已经超出最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large msg data recieve")
	}
	return msg, nil
}

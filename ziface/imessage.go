package ziface

type IMessage interface {
	//get&set
	GetMsgId() uint32
	GetMsgLen() uint32
	GetData() []byte

	SetMsgId(uint32)
	SetMsgLen(uint32)
	SetData([]byte)
}

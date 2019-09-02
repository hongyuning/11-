package iface

type IMessage interface {
	GetData() []byte
	GetDataLen() uint32
	GetMsgID() uint32

	SetData([]byte)
	SetDataLen(uint32)
	SetMsgID(uint32)
}

package net1

import "zinx/V1-basic-server/zinx/ZinxframeWork/iface"

type Message struct {
	data []byte
	len uint32
	msgid uint32//描述消息的字段，id
}

//创建message的方法
func NewMessage(data []byte,len,msgid uint32) iface.IMessage {
	return &Message{
		data:  data,
		len:   len,
		msgid: msgid,
	}
}
//方法
func (msg *Message) GetData() []byte {
	return msg.data
}
func (msg *Message) GetDataLen() uint32 {
	return msg.len
}
func (msg *Message) GetMsgID() uint32 {
	return msg.msgid
}

//++++++++++++++
func (msg *Message) SetData(data []byte) {
	msg.data = data
}
func (msg *Message) SetDataLen(len uint32) {
	msg.len = len
}
func (msg *Message) SetMsgID(msgid uint32) {
	msg.msgid = msgid
}

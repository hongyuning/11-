package net1

import "zinx/V1-basic-server/zinx/ZinxframeWork/iface"

type Request struct {
	//链接
	conn iface.IConnection
//使用message结构封装具体的数据和长度
   message iface.IMessage
}
func NewRequest (conn iface.IConnection,message iface.IMessage)iface.IRquest{
	return &Request{
		conn: conn,
		message:message,
	}
}
func (r *Request)GetConnection()iface.IConnection{
return r.conn
}
func (r *Request)GetMessage ()iface.IMessage{
	return r.message
}
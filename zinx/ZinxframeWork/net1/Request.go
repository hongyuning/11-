package net1

import "zinx/V1-basic-server/zinx/ZinxframeWork/iface"

type Request struct {
	//链接
	conn iface.IConnection
	//数据
	data []byte
	//数据长度
	len uint32
}
func NewRequest (conn iface.IConnection,data []byte,len uint32)iface.IRquest{
	return &Request{
		conn: conn,
		data: data,
		len:  len,
	}
}
func (r *Request)GetConnection()iface.IConnection{
return r.conn
}
func (r *Request)GetData ()[]byte{
	return r.data
}
func (r *Request)GetDatalen()uint32{
	return  r.len
}
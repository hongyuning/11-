package net1

import (
	"fmt"
	"net"
	"zinx/V1-basic-server/zinx/ZinxframeWork/iface"
)

type Connection struct {
	conn     *net.TCPConn
	connID   uint32
	isClosed bool
	//callBack iface.CallBackFunc //用户的注册业务处理函数
	router   iface.IRouter
}

func NewConnection (conn *net.TCPConn,cid uint32,router iface.IRouter) iface.IConnection {
	return &Connection{
		conn:     conn,
		connID:   cid,
		isClosed: false,
		//callBack: callback,
     router:router,
	}
}
func (c *Connection)Start(){
	for {
		buf:=make([]byte,4096)
		readcount,err:=c.conn.Read(buf)
		if err!=nil{
			fmt.Println("readcount  err",err)
			return
		}
		data:=string(buf[:readcount])
		fmt.Println("Server <=====client",readcount,data)
		req:=NewRequest(c,buf[:readcount],uint32(readcount))
		//c.callBack(req)
		c.router.PreHandle(req)
		c.router.Handle(req)
		c.router.PostHandle(req)
	}
}
func (c *Connection)Stop(){
fmt.Println("链接中断。。。。")
if c.isClosed{
	return
}
_=c.conn.Close()
}
func (c *Connection)Send(data []byte)(int,error){
	writebackcount,err:=c.conn.Write(data)
	if err!=nil{
		fmt.Println("writeCount  err",err)
		return 0,err
	}
	return writebackcount,nil
}
func (c *Connection)GetConnID()uint32{
return 0
}
func (c *Connection)GetTCPConn()*net.TCPConn{
return c.GetTCPConn()
}
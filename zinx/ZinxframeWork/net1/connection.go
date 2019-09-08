package net1

import (
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/V1-basic-server/zinx/ZinxframeWork/Configdecribe"
	"zinx/V1-basic-server/zinx/ZinxframeWork/iface"
)

type Connection struct {
	conn     *net.TCPConn
	connID   uint32
	isClosed bool
	//callBack iface.CallBackFunc //用户的注册业务处理函数
	//router iface.IRouter
	msghandle *MsgHandle
	msgChan   chan []byte
	myserver  iface.IServer

	propertysynx sync.RWMutex
	property     map[string]interface{}
}

func NewConnection(conn *net.TCPConn, cid uint32, msghandle *MsgHandle, server iface.IServer) iface.IConnection {
	return &Connection{
		conn:     conn,
		connID:   cid,
		isClosed: false,
		//callBack: callback,
		//router: router,
		msghandle: msghandle,
		msgChan:   make(chan []byte),
		myserver:  server,
		property:  make(map[string]interface{}),
	}
}
func (c *Connection) StartRead() {
	fmt.Println("Start  Read.....")
	defer fmt.Println("Stop  Read.....")
	defer c.Stop()
	for {
		//拆包
		dp := NewDataPack()
		buffer := make([]byte, dp.GetDataPackHeadLen())
		readcount, err := io.ReadFull(c.conn, buffer)
		if err != nil {
			fmt.Println("readcount err", err)
			return
		}
		fmt.Printf("读到数据头，%d\n", readcount)

		//只得到message的长度和类型
		message, err := dp.Unpack(buffer)
		if err != nil {
			fmt.Println("message err", err)
		}
		//校验数据是否是有效数据
		dataLen := message.GetDataLen()
		if dataLen == 0 {
			fmt.Printf("数据长度为0，无效")
			continue
		}
		//存放消息体
		dataBuf := make([]byte, dataLen)
		datacount, err := io.ReadFull(c.conn, dataBuf)
		if err != nil {
			if err != nil {
				fmt.Println("datacount err", err)
			}
		}
		fmt.Printf("Server <===== Client, data:%s,cnt:%d, msgid:%d\n", dataBuf, datacount, message.GetMsgID())
		message.SetData(dataBuf)
		req := NewRequest(c, message)
		//c.callBack(req)
		if Configdecribe.GlobalConfig.WorkSize > 0 {
			go c.msghandle.SendMsgToQue(req)
		} else {
			go c.msghandle.DoMsghandler(req)
		}
	}
}
func (c *Connection) Start() {
	go c.StartRead()
	go c.StarWriter()
	c.myserver.CallStartHookFunc(c)

}
func (c *Connection) StarWriter() {
	fmt.Println("start Write ......")
	defer fmt.Println("Stop Writer.....")
	for {
		datapack, ok := <-c.msgChan
		if !ok {
			return
		}
		_, err := c.conn.Write(datapack)
		if err != nil {
			fmt.Println("writeCount  err", err)
			return
		}

	}

}
func (c *Connection) Stop() {
	fmt.Println("链接中断。。。。")
	if c.isClosed {
		return
	}
	c.myserver.GetMsgMaganer().RemoveConn(int(c.connID))
	c.myserver.CallStopHookFunc(c)
	c.isClosed = true
	close(c.msgChan)

	err := c.conn.Close()
	if err != nil {
		fmt.Println("链接关闭失败......")
	}
}
func (c *Connection) Send(data []byte, msgid uint32) (int, error) {
	//defer  c.Stop()
	//封包
	dp := NewDataPack()
	datapack, err := dp.Pack(NewMessage(data, uint32(len(data)), msgid))
	if err != nil {
		fmt.Println("datapack   err", err)
		return -1, err
	}
	fmt.Println("开始写入.....")
	c.msgChan <- datapack

	return -1, nil
}
func (c *Connection) GetConnID() uint32 {
	return c.connID
}
func (c *Connection) GetTCPConn() *net.TCPConn {
	return c.GetTCPConn()
}
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertysynx.Lock()
	fmt.Println(key, value)
	c.property[key] = value
	c.propertysynx.Unlock()
}
func (c *Connection) GetProperty(key string) interface{} {
	c.propertysynx.RLock()
	value := c.property[key]
	c.propertysynx.RUnlock()
	return value
}

package net1

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

func TestDataPack(t *testing.T) {
	go func() {
		listen,err:=net.Listen("tcp","0.0.0.0:8888")
		if err!=nil{
			t.Error("listen err",err)
		}
		conn,err:=listen.Accept()
		if err!=nil{
			t.Error("conn err",err)
		}
		for  {
			buffer:=make([]byte,8)
			readcount,err:=io.ReadFull(conn,buffer)
			if err!=nil{
				t.Error("readcount err",err)
			}
			fmt.Printf("读到数据头，%d\n",readcount)


		//拆包
		dp:=NewDataPack()
		message,err:=dp.Unpack(buffer)
			if err!=nil{
				t.Error("message err",err)
			}
		//校验数据是否是有效数据
		dataLen:=message.GetDataLen()
		if dataLen==0{
			fmt.Printf("数据长度为0，无效")
			continue
		}
		//存放消息体
dataBuf:=make([]byte,dataLen)
datacount,err:=io.ReadFull(conn,dataBuf)
			if err!=nil{
				t.Error("datacount err",err)
			}
			fmt.Printf("Server <===== Client, data:%s,cnt:%d, msgid:%d\n", dataBuf, datacount, message.GetMsgID())
	}
	}()

	go func() {
		//封包，发送
		//把多个包黏在一起，一起发送
		//1. 准备数据（封包）
		data1 := []byte("你好")
		data2 := []byte("Hello world")
		data3 := []byte("国庆节即将到来!")

		//a. 创建message
		msg0 := NewMessage([]byte{}, 0, 0)
		msg1 := NewMessage(data1, uint32(len(data1)), 0)
		msg2 := NewMessage(data2, uint32(len(data2)), 1)
		msg3 := NewMessage(data3, uint32(len(data3)), 2)

		//b. 对message进行封包
		dp := NewDataPack()
		info0, _ := dp.Pack(msg0)
		info1, _ := dp.Pack(msg1)
		info2, _ := dp.Pack(msg2)
		info3, _ := dp.Pack(msg3)

		//将三个消息的字节流拼接到一起，一次性发送给服务器
		infoSend := append(info0, info1...)
		infoSend = append(infoSend, info2...)
		infoSend = append(infoSend, info3...)

		//2. 发送
		conn, err := net.Dial("tcp", "127.0.0.1:8888")
		if err != nil {
			t.Errorf("client dial err:%v\n", err)
			return
		}

		cnt, err := conn.Write(infoSend)
		if err != nil {
			t.Errorf("client send err:%v\n", err)
			return
		}

		fmt.Println("Client ====> Server cnt:", cnt)
	}()

	//select {}
	time.Sleep(2 * time.Second)


}
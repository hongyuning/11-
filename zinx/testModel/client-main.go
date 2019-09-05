package  main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/V1-basic-server/zinx/ZinxframeWork/net1"
)

func main(){

		//封包，发送
		//把多个包黏在一起，一起发送
		//1. 准备数据（封包）
		//2. 发送
		conn, err := net.Dial("tcp4", "127.0.0.1:8848")
		if err != nil {
			fmt.Println("conn  err", err)
			return
		}
		go func() {
			for  {
				time.Sleep(2*time.Second)
				data1 := []byte("你好")
				data2 := []byte("Hello world")
				data3 := []byte("国庆节即将到来!")

				//a. 创建message
				msg0 := net1.NewMessage([]byte{}, 0, 0)
				msg1 := net1.NewMessage(data1, uint32(len(data1)), 0)
				msg2 := net1.NewMessage(data2, uint32(len(data2)), 1)
				msg3 := net1.NewMessage(data3, uint32(len(data3)), 2)

				//b. 对message进行封包
				dp := net1.NewDataPack()
				info0, _ := dp.Pack(msg0)
				info1, _ := dp.Pack(msg1)
				info2, _ := dp.Pack(msg2)
				info3, _ := dp.Pack(msg3)

				//将三个消息的字节流拼接到一起，一次性发送给服务器
				infoSend := append(info0, info1...)
				infoSend = append(infoSend, info2...)
				infoSend = append(infoSend, info3...)

				cnt, err := conn.Write(infoSend)
				if err != nil {
					fmt.Println("conn  err", err)
					return
				}
				fmt.Println("Client ====> Server cnt:", cnt)

			}
		}()

		//++++++++++++++++
		 go func() {
			 for  {
				 time.Sleep(2*time.Second)
				 dp1:=net1.NewDataPack()
				 buffer:=make([]byte,dp1.GetDataPackHeadLen())
				 readcount,err:=io.ReadFull(conn,buffer)
				 if err!=nil{
					 fmt.Println("readcount  err", err)
					 return
				 }
				 fmt.Printf("读到数据头，%d\n",readcount)


				 //拆包
				 message,err:=dp1.Unpack(buffer)
				 if err!=nil{
					 fmt.Println("conn  err", err)
					 return
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
				 if err!=nil {
					 fmt.Println("datacount  err", err)
					 return
				 }
				 fmt.Printf("Server =====>> Client, data:%s,cnt:%d, msgid:%d\n", dataBuf, datacount, message.GetMsgID())

			 }
		}()
	select {

		}
		//+++++++++++++++++++
	}

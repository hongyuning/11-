package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"math/rand"
	"net"
	"time"
	message "zinx/V1-basic-server/momoGame/PB-protobuff"
)

type Client struct {
	Pid    int
	X      int
	Y      int
	Z      int
	V      int
	Conn   net.Conn
	online chan bool
}

func NewClient(ip string, port int) *Client {
	adrres := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.Dial("tcp", adrres)
	if err != nil {
		panic(err)
	}
	return &Client{
		Pid:    0,
		X:      0,
		Y:      0,
		Z:      0,
		V:      0,
		Conn:   conn,
		online: make(chan bool),
	}
}
func (ct *Client) Start() {
	//接收数据
	go func() {
		for {
			fmt.Println("客户端处理业务...")
			time.Sleep(2 * time.Second)
			headdata := make([]byte, 8)
			_, err := io.ReadFull(ct.Conn, headdata)
			if err != nil {
				fmt.Println("headdata read err:", err)
			}
			//拆包
			dp := NewDataPack()
			message1, err := dp.UnDataPack(headdata)
			if err != nil {
				fmt.Println("messga get err:", err)
			}
			if message1.len== 0 {
				fmt.Println("无具体数据，无需读取!")
				continue
			}
			data := make([]byte, message1.len)
			_, err = io.ReadFull(ct.Conn, data)
			if err != nil {
				fmt.Println("data read err:", err)
			}
			message1.data = data
			ct.HanleMsg(message1)
		}
	}()
	//发送数据
	<-ct.online
	go func() {
		for   {
         ct.robotAction()
         time.Sleep(2*time.Second)
		}
	}()
	select {
	}
}

func (ct *Client) HanleMsg(msg *Message) {
	fmt.Println("同步玩家消息.....", msg.msgid)
	var syncplayid message.SyncPid
	if msg.msgid == 1 {
		fmt.Println("开始同步玩家消息......")
		err := proto.Unmarshal(msg.data, &syncplayid)
		if err != nil {
			fmt.Println("syncplayid err:", err)
			return
		}
		ct.Pid = int(syncplayid.Pid)
	} else if msg.msgid == 200 {
		fmt.Println("获取广播逻辑......")
        var BroadCast  message.BroadCast
		err:=proto.Unmarshal(msg.data,&BroadCast)
		if err != nil {
			fmt.Println("broadcast err:", err)
			return
		}
		//判断具体的业务类型：1-聊天  2-位置 4-玩家移动
		if BroadCast.Tp == 2 && BroadCast.Pid == int32(ct.Pid) {
			ct.X = int(BroadCast.GetP().X)
			ct.Y = int(BroadCast.GetP().Y)
			ct.Z = int(BroadCast.GetP().Z)
			ct.V = int(BroadCast.GetP().V)
			fmt.Printf("玩家 id: %d 已经上线成功, 坐标: X :%d Y : %d Z: %d V:%d\n",
				ct.Pid, ct.X, ct.Y, ct.Z, ct.V)
			ct.online<-true
		}else if BroadCast.Tp == 1 && BroadCast.Pid == int32(ct.Pid) {
			fmt.Println("世界聊天：玩家:%d 说的话是:%s", ct.Pid, BroadCast.GetContent())
		}
	}
}

func (ct *Client)robotAction(){
	rander:=rand.Intn(2)
	takeInfo :=fmt.Sprintf(  "大家好,欠马云的钱要还给张勇了！！！,我是玩家%d",ct.Pid)
	if rander==0{
		//聊天
		protosenddata:=message.Talk{
			Content:     takeInfo       ,
		}
		ct.SenMessage(2,&protosenddata)
	}else if rander==1{
		//自动移动
		x := ct.X
		z := ct.Z

		randPos := rand.Intn(2)
		if randPos == 0 {
			//0, x,z加上一个数据
			x += rand.Intn(10)
			z += rand.Intn(10)
		} else {
			//1, z,z减去一个数据
			x -= rand.Intn(10)
			z -= rand.Intn(10)
		}

		//纠正坐标
		if x > 410 {
			x = 410
		} else if x < 85 {
			x = 85
		}

		if z > 400 {
			z = 400
		} else if z < 75 {
			z = 75
		}

		//面部朝向
		randv := rand.Intn(2)
		v := ct.V
		if randv == 0 {
			v = 25
		} else {
			v = 350
		}
		protopostion:=message.Position{
			X:                   float32(x) ,
			Y:                    float32(ct.Y),
			Z:                    float32(z),
			V:                    float32(v),
		}
		ct.SenMessage(3,&protopostion)
	}
}
func (ct *Client)SenMessage(msgid uint32,senddata proto.Message){
binaryInfo,err:=proto.Marshal(senddata)
if err!=nil{
	fmt.Println("binaryInfo err",err)
	return
}
msg:=NewMessage(msgid,binaryInfo)
dp:=NewDataPack()
SendPack,err:=dp.Pack(msg)
	if err!=nil{
		fmt.Println("sendpack err",err)
		return
	}
cnt,_:=ct.Conn.Write(SendPack)
fmt.Println("Client ====> Server , cnt:", cnt)
}
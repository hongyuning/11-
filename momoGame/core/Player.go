package core

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"sync"
	message "zinx/V1-basic-server/momoGame/PB-protobuff"
	"zinx/V1-basic-server/zinx/ZinxframeWork/iface"
)

type  Player struct {
	Pid int
	Conn iface.IConnection
	X int //横向坐标
	Y int  //高度，已经制定好
	Z  int  //纵向坐标
	V  int  //面部朝向
}
var PidGen  int=1
var lock sync.Mutex  //防止两个客户端登陆
func NewPlayer (conn iface.IConnection)*Player{
	lock.Lock()

	p:=&Player{
		Pid:  PidGen,
		Conn: conn,
		X:    160+rand.Intn(10),
		Y:    0,
		Z:    160+rand.Intn(10),
		V:    0,
	}
	PidGen++
	lock.Unlock()
	return  p
}
//同步玩家id
func (p *Player)SyncPid (){
	fmt.Println("开始同步pid给客户端.....",p.Pid)
	protodata:=message.SyncPid{
		Pid:                  int32(p.Pid),
	}
    p.SendMsg(1,&protodata)
	fmt.Println("发送成功.......")

}
//编码并发送
func (p *Player)SendMsg(msgid int,protoData proto.Message){
	protoInfo,err:=proto.Marshal(protoData)
	if err!=nil{
		fmt.Println("proto编码失败......",err)
		return
	}
	//发送给客户端
	_,err=p.Conn.Send(protoInfo,uint32(msgid))
	if err!=nil{
		fmt.Println("发送数据失败......",err)
		return
	}
}
//初始化玩家位置
func (p  *Player)SyncPosition() {
	protodata := message.BroadCast{
		Pid: int32(p.Pid),
		Tp:  2, //代表位置
		Data: &message.BroadCast_P{P: &message.Position{
			X: float32(p.X),
			Y: float32(p.Y),
			Z: float32(p.Z),
			V: float32(p.V),
		}},
	}
	p.SendMsg(200,&protodata)

}
func (p *Player)SendTalkMessageToAllPlayers(talkmessage string){
talkMesage:=message.BroadCast{
	Pid:                  int32(p.Pid),
	Tp:                   1,
	Data:                 &message.BroadCast_Content{Content:talkmessage},

}
players:=WorldMgrGlobal.GetAllPlayer()
for  _,player:=range players{
	player.SendMsg(200,&talkMesage)
}
}
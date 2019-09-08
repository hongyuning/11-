package router
import (
	"fmt"
	"github.com/golang/protobuf/proto"
	message "zinx/V1-basic-server/momoGame/PB-protobuff"
	"zinx/V1-basic-server/momoGame/core"
	"zinx/V1-basic-server/zinx/ZinxframeWork/iface"
	"zinx/V1-basic-server/zinx/ZinxframeWork/net1"
)

type WorldChat struct {
	net1.Router
}

func (router *WorldChat)Handle(req  iface.IRquest){
	fmt.Println("聊天的路由被调用.....")
	talkContent:=req.GetMessage().GetData()
	var talk_proto  message.Talk
	err:=proto.Unmarshal(talkContent,&talk_proto)
	if err!=nil{
		fmt.Println("proto解码失败.....")
		return
	}
	//得到真正的聊天内容
	talkmessage:=talk_proto.Content
	pidinterdace:=req.GetConnection().GetProperty("pid")
	pid,ok:=pidinterdace.(int)
	if !ok{
		fmt.Printf("%T,%v\n",pid,pid)
		fmt.Println("类型断言失败....")
		return
	}
	player:=core.WorldMgrGlobal.GetPlayByPid(pid)
	player.SendTalkMessageToAllPlayers(talkmessage)
}
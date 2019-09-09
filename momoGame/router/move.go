package router

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	message "zinx/V1-basic-server/momoGame/PB-protobuff"
	"zinx/V1-basic-server/momoGame/core"
	"zinx/V1-basic-server/zinx/ZinxframeWork/iface"
	"zinx/V1-basic-server/zinx/ZinxframeWork/net1"
)

type Move struct {
	net1.Router
}

//处理玩家移动胡路由
func (router *Move)Handle(req iface.IRquest){
	fmt.Println("调用玩家移动的路由成功......")
	protoinfo:=req.GetMessage().GetData()
	var  PlayPosition  message.Position
	err:=proto.Unmarshal(protoinfo,&PlayPosition)
	if err!=nil{
		fmt.Println("移动解码错误.....",err)
		return
	}
	pidinterdace:=req.GetConnection().GetProperty("pid")
	pid,ok:=pidinterdace.(int)
	if !ok{
		fmt.Println("类型断言失败.....")
		return
	}
	player:=core.WorldMgrGlobal.GetPlayByPid(pid)
   player.UpdataPostion(PlayPosition.X,PlayPosition.Y,PlayPosition.Z,PlayPosition.V)
}
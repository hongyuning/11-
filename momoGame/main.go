package main

import (
	"fmt"
	"zinx/V1-basic-server/momoGame/core"
	"zinx/V1-basic-server/momoGame/router"
	"zinx/V1-basic-server/zinx/ZinxframeWork/iface"
	"zinx/V1-basic-server/zinx/ZinxframeWork/net1"
)

func main() {
	server := net1.NewServer("MMOGameServer")
	//
	server.RegistStartHookFunc(OnConnStartHookFunc)
	server.RegistStopHookFunc(OnConnStopHookFunc)
	server.AddRouter(2,&router.WorldChat{})
	server.AddRouter(3,&router.Move{})
	server.Serve()

}

//钩子函数，做业务之前的准备工作；路由做具体业务的处理；框架做的是数据的传输
func OnConnStartHookFunc(conn iface.IConnection) {
	//创建玩家
	player := core.NewPlayer(conn)
	player.SyncPid()
	//player.SyncPosition()
	conn.SetProperty("pid",player.Pid)

	//添加世界管理器
	core.WorldMgrGlobal.AddPlayer(player)
	totalplayers := len(core.WorldMgrGlobal.Players)
	fmt.Println("新登录玩家为：", player.Pid, "当前玩家总数：", totalplayers)
	player.SybcSurroundPlayersPostion()
}
func OnConnStopHookFunc(conn iface.IConnection) {
	//conn.Stop()
	pidInterface := conn.GetProperty("pid")

	//断言
	pid, ok := pidInterface.(int)
	if !ok {
		fmt.Println("pid 断言失败!")
		return
	}

	p := core.WorldMgrGlobal.GetPlayByPid(pid)

	p.OffLine()
}

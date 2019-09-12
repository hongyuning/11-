package core

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"math/rand"

	"sync"
	message "zinx/V1-basic-server/momoGame/PB-protobuff"
	"zinx/V1-basic-server/zinx/ZinxframeWork/iface"
)

type Player struct {
	Pid  int
	Conn iface.IConnection
	X    int //横向坐标
	Y    int //高度，已经制定好
	Z    int //纵向坐标
	V    int //面部朝向
}

var PidGen int = 1
var lock sync.Mutex //防止两个客户端登陆
func NewPlayer(conn iface.IConnection) *Player {
	lock.Lock()

	p := &Player{
		Pid:  PidGen,
		Conn: conn,
		X:    160 + rand.Intn(10),
		Y:    0,
		Z:    160 + rand.Intn(10),
		V:    0,
	}
	PidGen++
	lock.Unlock()
	return p
}

//同步玩家id
func (p *Player) SyncPid() {
	fmt.Println("开始同步pid给客户端.....", p.Pid)
	protodata := message.SyncPid{
		Pid: int32(p.Pid),
	}
	p.SendMsg(1, &protodata)
	fmt.Println("发送成功.......")

}

//编码并发送
func (p *Player) SendMsg(msgid int, protoData proto.Message) {
	protoInfo, err := proto.Marshal(protoData)
	if err != nil {
		fmt.Println("proto编码失败......", err)
		return
	}
	//发送给客户端
	_, err = p.Conn.Send(protoInfo, uint32(msgid))
	if err != nil {
		fmt.Println("发送数据失败......", err)
		return
	}
}

//初始化玩家位置
func (p *Player) SyncPosition() {
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
	p.SendMsg(200, &protodata)

}
func (p *Player) SendTalkMessageToAllPlayers(talkmessage string) {
	talkMesage := message.BroadCast{
		Pid:  int32(p.Pid),
		Tp:   1,
		Data: &message.BroadCast_Content{Content: talkmessage},
	}
	players := WorldMgrGlobal.GetAllPlayer()
	for _, player := range players {
		player.SendMsg(200, &talkMesage)
	}
}

//玩家上线后同步位置
func (p *Player) SybcSurroundPlayersPostion() {
	BroadData := message.BroadCast{
		Pid: int32(p.Pid),
		Tp:  2,
		Data: &message.BroadCast_P{P: &message.Position{
			X: float32(p.X),
			Y: float32(p.Y),
			Z: float32(p.Z),
			V: float32(p.V),
		}},
	}
	//通过位置得到周围所有的玩家id
	pids := WorldMgrGlobal.gridmgr.GetSurroundPidsByPos(p.X, p.Z)
	//通过玩家id找到所有的玩家
	players := make([]*Player, 0, len(pids))
	fmt.Println(len(pids), "-------------------")
	for _, playerid := range pids {
		player := WorldMgrGlobal.GetPlayByPid(playerid)
		player.SendMsg(200, &BroadData)
		players = append(players, player)
	}

	//获取周边所有人的信息
	protoplayers1 := make([]*message.Player, 0, len(players))
	for _, playerr := range players {
		proto_player := message.Player{
			Pid: int32(playerr.Pid),
			P: &message.Position{
				X: float32(playerr.X),
				Y: float32(playerr.Y),
				Z: float32(playerr.Z),
				V: float32(playerr.V),
			},
		}
		protoplayers1 = append(protoplayers1, &proto_player)
	}
	//拼接结构
	syncPlayersProto := message.SyncPlayerProto{
		Players: protoplayers1,
	}
	p.SendMsg(202, &syncPlayersProto)
}
func (p *Player) getSurroundPlayersByPos() []*Player {
	//通过位置得到周围所有的玩家id
	pids := WorldMgrGlobal.gridmgr.GetSurroundPidsByPos(p.X, p.Z)
	//通过玩家id找到所有的玩家
	players := make([]*Player, 0, len(pids))
	fmt.Println(len(pids), "-------------------")
	for _, playerid := range pids {
		player := WorldMgrGlobal.GetPlayByPid(playerid)
		players = append(players, player)
	}
	return players
}

//玩家下线
func (p *Player) OffLine() {
	fmt.Println("playerid:", p.Pid, ", 即将下线!")
	//1. 通知周边的人，player下线，从这人的视野中消失
	surroundPlayers := p.getSurroundPlayersByPos()

	//2.拼接SyncPid结构
	proto_struct := message.SyncPid{Pid: int32(p.Pid)}

	//发送给周边的人，告知当前玩家已经下线，消息类型：201
	for _, player := range surroundPlayers {
		player.SendMsg(201, &proto_struct)
	}

	//3. 从世界管理器中将player删除
	WorldMgrGlobal.RemovePlayer(p)

	//4. 从格子管理器中，将player位置信息删掉
	WorldMgrGlobal.gridmgr.RemovePlayerIdFromGridByPos(p.Pid, p.X, p.Z)
}

//同步玩家位置
func (p *Player) UpdataPostion(X, Y, Z, V float32) {

	//处理跨格子显示问题
	//1. 计算当前的玩家是否已经跨越格子了
	//a. 旧的格子id：oldgid：   13
	oldgrid := WorldMgrGlobal.gridmgr.GetGidFromPos(p.X, p.Z)
	//b. 新的格子id：newgid    9
	newgrid := WorldMgrGlobal.gridmgr.GetGidFromPos(int(X), int(Z))
	//c . 判断，如果格子id发生变化
	if oldgrid != newgrid {
		//触发格子切换的逻辑
		//将player添加到新的grid中
		WorldMgrGlobal.gridmgr.AddPlayerIdToGrid(p.Pid, newgrid) //todo 为什么只是添加了玩家的id
		//从旧的grid删除player
		WorldMgrGlobal.gridmgr.RemovePlayIdFromGrid(p.Pid, oldgrid)
		//触发个其他玩家显示的逻辑 <<====
		p.OnExchangeGrid(oldgrid, newgrid)
	}
	//更新玩家位置
	p.X = int(X)
	p.Y = int(Y)
	p.Z = int(Z)
	p.V = int(V)
	protodata := message.BroadCast{
		Pid: int32(p.Pid),
		Tp:  4,
		Data: &message.BroadCast_P{P: &message.Position{
			X: X,
			Y: Y,
			Z: Z,
			V: V,
		}},
	}
	players := p.getSurroundPlayersByPos()
	for _, player := range players {
		player.SendMsg(200, &protodata)
	}
}

//玩家的显示逻辑
func (p *Player) OnExchangeGrid(oldgridid, newgridid int) {
	fmt.Println("玩家格子显示逻辑被调用.....")
	oldSurroundGrids := WorldMgrGlobal.gridmgr.GetSurroudingGridsByGrid(oldgridid)
	newSurroundGrids := WorldMgrGlobal.gridmgr.GetSurroudingGridsByGrid(newgridid)
	oldgridsGroup := make(map[int]bool, len(oldSurroundGrids))
	for _, grid := range oldSurroundGrids {
		oldgridsGroup[grid.gridId] = true
	}
	newgridsGroup := make(map[int]bool, len(newSurroundGrids))
	for _, gridd := range newSurroundGrids {
		newgridsGroup[gridd.gridId] = true
	}
	LeaveGrids := make([]*Grid, 0)
	for _, oldgridsid := range oldSurroundGrids {
		if !newgridsGroup[oldgridsid.gridId] {
			//获取到需要离开视野的格子id集合
			LeaveGrids = append(LeaveGrids, oldgridsid)
		}
	}
	protodata := message.SyncPid{
		Pid: int32(p.Pid),
	}

	//得到格子
	for _, grid := range LeaveGrids {
		players := WorldMgrGlobal.GetPlayersByGrid(grid.gridId)
		for _, playerr := range players {
			//向每个玩家播报自己，要求在格子里存储
			playerr.SendMsg(201, &protodata)
			//向自己播报每个玩家的位置
			Another_prodata := message.SyncPid{
				Pid: int32(playerr.Pid),
			}
			//发送给自己其他玩家
			p.SendMsg(201, &Another_prodata)
		}
	}
     mine_protoPostion:=message.BroadCast{
		 Pid:                  int32(p.Pid),
		 Tp:                   2,
		 Data:                 &message.BroadCast_P{P:&message.Position{
			 X:                    float32(p.X),
			 Y:                    float32(p.Y),
			 Z:                    float32(p.Z),
			 V:                    float32(p.V),
		 }},
	 }


	//视野显示处理
	//1.找到需要显示的格子
	for _, grid := range newSurroundGrids {
		//找到旧格子里没有的格子
        if !oldgridsGroup[grid.gridId]{
         newplayers:= WorldMgrGlobal.GetPlayersByGrid(grid.gridId)
        	for _,player:=range newplayers{
        		//1.向所有玩家广播自己的位置
        		player.SendMsg(200,&mine_protoPostion)
        		//2.向自己广播所有玩家的位置
        		Another_Player_Postion:=message.BroadCast{
					Pid:                  int32(player.Pid),
					Tp:                   2,
					Data:                 &message.BroadCast_P{P:&message.Position{
						X:                    float32(player.X),
						Y:                    float32(player.Y),
						Z:                    float32(player.Z),
						V:                    float32(player.V),
					}},
				}
				p.SendMsg(200,&Another_Player_Postion)
			}//for

		}
	}

}

package core

import (
	"sync"
	"zinx/V1-basic-server/momoGame/momoGameutils"
)

type WorldManager struct {
	//画布
	gridmgr *GridManager
	//玩家集合
	Players map[int]*Player
	//读写锁
	PlayerLock sync.RWMutex
}
func NewWorldManager ()*WorldManager{
	wm:=&WorldManager{
		gridmgr:    NewGridManager(momoGameutils.GM_MIN_X, momoGameutils.GM_MAX_X, momoGameutils.GM_CNTS_X, momoGameutils.GM_MIN_Y, momoGameutils.GM_MAX_Y, momoGameutils.GM_CNTS_Y),
		Players:    make(map[int]*Player),
		PlayerLock: sync.RWMutex{},
	}
	return wm
}
var WorldMgrGlobal *WorldManager
//提前初始化好一个世界管理器
func init(){
	WorldMgrGlobal=NewWorldManager()
}
//增加玩家
func (wm *WorldManager)AddPlayer (player *Player){
	wm.PlayerLock.Lock()
	wm.Players[player.Pid]=player
	wm.gridmgr.AddPlayerIdToGridByPos(player.Pid,player.X,player.Z)
	wm.PlayerLock.Unlock()
}
//删除玩家
func (wm *WorldManager)RemovePlayer (player *Player){
	wm.PlayerLock.Lock()
	delete(wm.Players,player.Pid)
	wm.gridmgr.RemovePlayerIdFromGridByPos(player.Pid,player.X,player.Z)
	wm.PlayerLock.Unlock()
}
//获取指定玩家
func (wm *WorldManager)GetPlayByPid (pid int)*Player{
  wm.PlayerLock.RLock()
 player:= wm.Players[pid]
  wm.PlayerLock.RUnlock()
 return  player
}
//获取所有玩家
func (wm *WorldManager)GetAllPlayer()[]*Player{
wm.PlayerLock.RLock()
   players:=make([]*Player,0,len(wm.Players))
	for _,player:=range wm.Players{
  players=append(players,player)
	}
wm.PlayerLock.RUnlock()
return  players
}
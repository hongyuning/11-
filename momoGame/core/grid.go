package core

import (
	"fmt"
	"sync"
)

type Grid struct {
	gridId int//格子id
	minX,maxX int //格子横坐标的最小值和最大值
	minY,maxY int
	players map[int]interface{}  //玩家集合
	playerlock sync.RWMutex
}

func NewGrid (gridid ,minx,maxx,miny,maxy  int)*Grid{
	return &Grid{
		gridId:   gridid  ,
		minX:       minx,
		maxX:       maxx,
		minY:       miny,
		maxY:       maxy,
		players:    make(map[int]interface{}),
		playerlock: sync.RWMutex{},
	}
}
//增加玩家
func (g *Grid)AddPlayer(playerid int,player interface{}){
fmt.Println("添加玩家",playerid)
g.playerlock.Lock()
g.players[playerid]=player
g.playerlock.Unlock()
}
//删除玩家
func (g *Grid) RemovePlayer(playerid int) {
	fmt.Println("删除玩家:", playerid)

	//操作map，要上锁
	g.playerlock.Lock()
	defer g.playerlock.Unlock()

	delete(g.players, playerid)
}

//获取当前格子内所有玩家id
func (g *Grid) GetAllPlayerIds() (playerids []int) {
	fmt.Println("获取所有玩家id")

	//操作map，要上锁
	g.playerlock.RLock()
	defer g.playerlock.RUnlock()

	for id := range g.players {
		playerids = append(playerids, id)
	}
	return
}
//实现打印的格式化输出
func (g *Grid)String ()string{
	return  fmt.Sprintf("gid : %d, minX: %d, maxX: %d, minY: %d, maxY:%d, playerCount: %d\n",
		g.gridId, g.minX, g.maxX, g.minY, g.maxY, len(g.players))
}

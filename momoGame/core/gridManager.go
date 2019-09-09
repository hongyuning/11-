package core

import (
	"fmt"
)

//格子管理器
type GridManager struct {
	minX, maxX, cntX int
	minY, maxY, cntY int
	grids            map[int]*Grid //小格子
}

func NewGridManager(minX, maxX, cntX, minY, maxY, cntY int) *GridManager {
	gridmanager := GridManager{
		minX:  minX,
		maxX:  maxX,
		cntX:  cntX,
		minY:  minY,
		maxY:  maxY,
		cntY:  cntY,
		grids: make(map[int]*Grid),
	}
	for y := 0; y < cntY; y++ {
		for x := 0; x < cntX; x++ {
			//获取格子的id
			gridid := y*cntX + x
			gridminxx := minX + x*gridmanager.GetWidth()
			gridmaxx := gridminxx + gridmanager.GetWidth()
			gridminy := minY + y*gridmanager.GetHigh()
			gridmaxy := gridminy + gridmanager.GetHigh()
			gridmanager.grids[gridid] = NewGrid(gridid, gridminxx, gridmaxx, gridminy, gridmaxy)
		}
	}
	return &gridmanager
}

//格式化打印
func (gm *GridManager) String() string {
	s := fmt.Sprintf("gridManager:\n minx:%d,maxx:%d,cntx:%d,miny:%d,maxy:%d,cnty;%d\n", gm.minX, gm.maxX, gm.cntX, gm.minY, gm.maxY, gm.cntY)
	for _, grid := range gm.grids {
		s += fmt.Sprint(grid)
	}
	return s
}

//获取单个格子宽度
func (gm *GridManager) GetWidth() int {
	return (gm.maxX - gm.minX) / gm.cntX
}

//获取单个格子的高度
func (gm *GridManager) GetHigh() int {
	return (gm.maxY - gm.minY) / gm.cntY
}

//添加玩家id到格子中
func (gm *GridManager) AddPlayerIdToGrid(pid, gridid int) {
	gm.grids[gridid].AddPlayer(pid, nil)
}

//从格子中删除playerid
func (gm *GridManager) RemovePlayIdFromGrid(pid, gridid int) {
	gm.grids[gridid].RemovePlayer(pid)
}

//返回所有玩家的id
func (gm *GridManager) GetAllPlayerId(gridid int) []int {
	return gm.grids[gridid].GetAllPlayerIds()
}

//把玩家的位置定位到每个小格子中
func (gm *GridManager) GetGidFromPos(x, y int) int {
	idx := (x - gm.minX) / gm.GetWidth()
	idy := (y - gm.minY) / gm.GetHigh()
	return idy*gm.cntX + idx
}

//1. 添加playerid到格子   ===> AddPlayerIdToGridByPos（pid, x，y)
func (gm *GridManager) AddPlayerIdToGridByPos(pid, x, y int) {
	gridId := gm.GetGidFromPos(x, y)
	gm.AddPlayerIdToGrid(pid, gridId)
}

//2. 从格子中删除playerid  ===》 RemovePlayerIdFromGridByPos（pid，x，y）
func (gm *GridManager) RemovePlayerIdFromGridByPos(pid, x, y int) {
	gridId := gm.GetGidFromPos(x, y)
	gm.RemovePlayIdFromGrid(pid, gridId)
}

//根据当前的格子，返回周边所有的格子
func (gm *GridManager) GetSurroudingGridsByGrid(gid int) (grids []*Grid) {
	grid, ok := gm.grids[gid]
	if !ok {
		return
	}
	grids = append(grids, grid)
	//判断左边是否有格子，如果有放到切片中，x>0,则认为有格子
	if grid.minX >= gm.minX+gm.GetWidth() {
		grids = append(grids, gm.grids[gid-1])
	}
	if grid.maxX <= gm.maxX-gm.GetWidth() {
		grids = append(grids, gm.grids[gid+1])
	}
	for _, gridd := range grids {
		//判断上下是否有格子
		if gridd.minY >= gm.minY+gm.GetHigh() {
			grids = append(grids, gm.grids[gridd.gridId-gm.cntX])
		}
		if gridd.maxY <= gm.maxY-gm.GetHigh() {
			grids = append(grids, gm.grids[gridd.gridId+gm.cntX])
		}
	}
	return
}

//通过位置得到周围格子
func (gm *GridManager) GetSurroudingGridsByPos(x, y int) (grids []*Grid) {
	gid := gm.GetGidFromPos(x, y)
	return gm.GetSurroudingGridsByGrid(gid)
}
//获取当前格子里的全部玩家id
func (gm *GridManager) GetPidsByGid(gid int) []int {
	playerids := gm.grids[gid].GetAllPlayerIds()
	return playerids
}
//通过位置获取周围玩家的id
func (gm *GridManager) GetSurroundPidsByPos(x, y int) (surroundPlayerId []int) {
	gid := gm.GetGidFromPos(x, y)
	surroundgids := gm.GetSurroudingGridsByGrid(gid)
	for _, gid := range surroundgids {
		surroundPlayerId = append(surroundPlayerId, gm.GetPidsByGid(gid.gridId)...)
	}
	fmt.Println(len(surroundgids),"**********")
	return
}

package test

import (
	"fmt"
	"testing"
	"zinx/V1-basic-server/momoGame/core"
)

//func TestNewGridManager(t *testing.T) {
//	gm := NewGridManager(0, 5, 5, 0, 100, 5)
//	fmt.Print(gm)
//}
func TestGridManager_GetSurroudingGridsByGrid(t *testing.T) {
	gm:= core.NewGridManager(0,25,5,0,25,5)
	for gid,_:=range gm.grids{
		surroundgrids:=gm.GetSurroudingGridsByGrid(gid)
		a:=[]int{}
		for _,gridss:=range surroundgrids{
			a=append(a,gridss.gridId)
		}
		fmt.Println("gid:.......",a)
	}
	fmt.Println(".......................")
}

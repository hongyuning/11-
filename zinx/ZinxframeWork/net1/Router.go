package net1

import (
	"fmt"
	"zinx/V1-basic-server/zinx/ZinxframeWork/iface"
)

type Router struct {

}

//处理业务前的数据处理
func (r *Router)PreHandle(req iface.IRquest){
	fmt.Println("  prehandle  callde   ...")
}
//真正的业务处理
func (r *Router)Handle (req iface.IRquest){
	fmt.Println("Handle called!")
}
//处理业务之后的清理工作
func (r *Router)PostHandle(req iface.IRquest){
	fmt.Println("PostHandle called!")
}
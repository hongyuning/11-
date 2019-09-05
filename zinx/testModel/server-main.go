package  main

import (
	"fmt"
	"strings"
	"zinx/V1-basic-server/zinx/ZinxframeWork/iface"
	net "zinx/V1-basic-server/zinx/ZinxframeWork/net1"
)
type TestRouter struct {
	net.Router
}
//处理业务前的数据处理
func (r *TestRouter)PreHandle(req iface.IRquest){
	fmt.Println("  prehandle  callde   ...")
}
//真正的业务处理
func (r *TestRouter)Handle (req iface.IRquest){
	data:=req.GetMessage().GetData()
	conn:=req.GetConnection()
	writeBackdata:=strings.ToUpper(string(data))
	writedata:=[]byte(writeBackdata)
	sendnum,err:=conn.Send(writedata[:len(writedata)],200)
	if err!=nil{
		fmt.Println(err)
		return
	}
	fmt.Println("Server====>Client  ",sendnum,writeBackdata,200)
	fmt.Println("Hanile called!")
}
//处理业务之后的清理工作
func (r *TestRouter)PostHandle(req iface.IRquest){
	fmt.Println("PostHandle called!")
}
type MoveRouter struct {
	net.Router
}
func (mr *MoveRouter)Handle (request iface.IRquest){
	fmt.Println("处理移动请求的路由逻辑")
}
type AttackRouter struct {
	net.Router
}

func (router *AttackRouter) Handle(req iface.IRquest) {
	fmt.Println("处理攻击请求的路由逻辑")
}



func main() {
	a:= net.NewServer("5555")
	a.AddRouter(0,&TestRouter{})
	a.AddRouter(1,&MoveRouter{})
	a.AddRouter(2,&AttackRouter{})
	a.Serve()
}
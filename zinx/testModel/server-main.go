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
	data:=req.GetData()
	conn:=req.GetConnection()
	writeBackdata:=strings.ToUpper(string(data))
	writedata:=[]byte(writeBackdata)
	sendnum,err:=conn.Send(writedata[:len(writedata)])
	if err!=nil{
		fmt.Println(err)
		return
	}
	fmt.Println("Server====>Client  ",sendnum,writeBackdata)
	fmt.Println("Hanile called!")
}
//处理业务之后的清理工作
func (r *TestRouter)PostHandle(req iface.IRquest){
	fmt.Println("PostHandle called!")
}


func main() {
	a:= net.NewServer("5555")
	a.AddRouter(&TestRouter{})
	a.Serve()
}
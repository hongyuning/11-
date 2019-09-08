package net1

import (
	"fmt"
	"zinx/V1-basic-server/zinx/ZinxframeWork/Configdecribe"
	"zinx/V1-basic-server/zinx/ZinxframeWork/iface"
)

type  MsgHandle struct {
	msghanler map[uint32]iface.IRouter
	//工作的消息队列
	worksize int
	taskQueue []chan iface.IRquest
}

func NewMsgHandle ()*MsgHandle{
	return &MsgHandle{
		msghanler:make(map[uint32]iface.IRouter),
		worksize: Configdecribe.GlobalConfig.WorkSize,
		taskQueue:make([]chan  iface.IRquest,Configdecribe.GlobalConfig.WorkSize),
	}
}

func (mh *MsgHandle)AddRouter(msgid uint32,router  iface.IRouter){
	_,ok:=mh.msghanler[msgid]
	if ok{
		fmt.Println("路由已经存在，不需要添加......")
		return
	}
	mh.msghanler[msgid]=router
fmt.Println("路由添加成功.....",msgid)
}
//执行路由的hanle函数
func(mh *MsgHandle) DoMsghandler(request iface.IRquest){
	msgid:=request.GetMessage().GetMsgID()
	router,ok:=mh.msghanler[msgid]
	if !ok{
		fmt.Println("路由不存在......")
		return
	}
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}

//启动work池，给每一个消息队列分配空间
func (mh *MsgHandle)StartWork (){
	for i:=0;i<mh.worksize;i++{
		fmt.Println("启动worker.....")
		mh.taskQueue[i]=make(chan iface.IRquest,Configdecribe.GlobalConfig.TashqueSize)
		go func(i int) {
			for{
				req:=<-mh.taskQueue[i]
				mh.DoMsghandler(req)
			}
		}(i)
	}
}
//分配队列
func (mh *MsgHandle)SendMsgToQue(request iface.IRquest){
connID:=request.GetConnection().GetConnID()
quequeid:=int(connID)%mh.worksize
mh.taskQueue[quequeid]<-request

}
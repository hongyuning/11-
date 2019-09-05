package net1

import (
	"fmt"
	"zinx/V1-basic-server/zinx/ZinxframeWork/iface"
)

type  MsgHandle struct {
	msghanler map[uint32]iface.IRouter
}

func NewMsgHandle ()*MsgHandle{
	return &MsgHandle{msghanler:make(map[uint32]iface.IRouter)}
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
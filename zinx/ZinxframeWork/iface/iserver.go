package iface

type IServer interface {
	Start()
	Stop ()
	Serve()
	AddRouter(uint32,IRouter)
	GetMsgMaganer() IConnManager
    //两个钩子函数的注册
    RegistStartHookFunc (func(connection IConnection))
	RegistStopHookFunc (func(connection IConnection))
	//两个钩子函数的调用
	CallStartHookFunc(connection IConnection)
	CallStopHookFunc(connection IConnection)
}
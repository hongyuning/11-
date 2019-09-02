package iface

type IServer interface {
	Start()
	Stop ()
	Serve()
	AddRouter(IRouter)
}
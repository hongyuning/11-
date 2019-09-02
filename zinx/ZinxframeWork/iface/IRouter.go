
package iface

type IRouter interface {
	PreHandle(IRquest)
	Handle(IRquest)
	PostHandle(IRquest)
}


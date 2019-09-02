package iface

import "net"

type IConnection interface {
	Start()
	Stop()
	Send([]byte)(int,error)
	GetConnID()uint32
	GetTCPConn()*net.TCPConn
}

type CallBackFunc func(rquest IRquest)
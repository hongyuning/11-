package net1

import (
	"fmt"
	"net"
	"zinx/V1-basic-server/zinx/ZinxframeWork/Configdecribe"
	"zinx/V1-basic-server/zinx/ZinxframeWork/iface"
)

type Server struct {
	IP      string
	Port    uint32
	Name    string
	Version string
	//Router  iface.IRouter
	msghandle *MsgHandle
}

func NewServer(name string) iface.IServer {
	return &Server{
		IP:      Configdecribe.GlobalConfig.Ip,
		Port:    Configdecribe.GlobalConfig.Port,
		Name:    Configdecribe.GlobalConfig.Name,
		Version: Configdecribe.GlobalConfig.Version,
		msghandle:NewMsgHandle(),
	}
}

func (s *Server) Start() {
	fmt.Println("start 被调用")
	address := fmt.Sprintf("%s:%d", s.IP, s.Port)
	//address:=s.IP+":"+string(s.Port)
	tcpaddr, err := net.ResolveTCPAddr(s.Version, address)
	if err != nil {
		fmt.Println("tcp address  err", err)
		return
	}
	TCPListener, err := net.ListenTCP(s.Version, tcpaddr)
	if err != nil {
		fmt.Println("TCPListener err", err)
		return
	}
	var cid uint32
	go func() {
		for {
			conn, err := TCPListener.AcceptTCP()
			if err != nil {
				fmt.Println("conn  err", err)
				return
			}
			myconnection := NewConnection(conn, cid,s.msghandle)
			cid++
			go myconnection.Start()

		}
	}()
}
func (s *Server) Stop() {
	fmt.Println("stop 被调用")
}
func (s *Server) Serve() {
	fmt.Println("server 被调用")
	s.Start()
	fmt.Println(11)
	for {
	}
}
func (s *Server) AddRouter(msg uint32,router iface.IRouter) {
	fmt.Println("AddRouter被调用。。。。。")
	s.msghandle.AddRouter(msg,router)
}

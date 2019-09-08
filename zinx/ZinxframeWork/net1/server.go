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
	msghandle   *MsgHandle
	connmanager iface.IConnManager

	StartHookFunc func(connection iface.IConnection)
	StopHookFunc  func(connection iface.IConnection)
}

func NewServer(name string) iface.IServer {
	return &Server{
		IP:          Configdecribe.GlobalConfig.Ip,
		Port:        Configdecribe.GlobalConfig.Port,
		Name:        Configdecribe.GlobalConfig.Name,
		Version:     Configdecribe.GlobalConfig.Version,
		msghandle:   NewMsgHandle(),
		connmanager: NewConnManager(), //链接管理
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

	//启动消息队列
	s.msghandle.StartWork()
	var cid uint32
	go func() {
		for {
			conn, err := TCPListener.AcceptTCP()
			if err != nil {
				fmt.Println("conn  err", err)
				return
			}
			if s.connmanager.GetConnCount() >= Configdecribe.GlobalConfig.MustConnCount {
				fmt.Println("超过最大链接数量.....链接被拒绝")
				_ = conn.Close()
				continue
			}
			myconnection := NewConnection(conn, cid, s.msghandle, s)
			s.connmanager.AddConn(myconnection)
			cid++

			go myconnection.Start()

		}
	}()
}
func (s *Server) Stop() {
	fmt.Println("stop 被调用")
	s.connmanager.ClearConn()
}
func (s *Server) Serve() {
	fmt.Println("server 被调用")
	s.Start()
	fmt.Println(11)
	select {}
}
func (s *Server) AddRouter(msg uint32, router iface.IRouter) {
	fmt.Println("AddRouter被调用。。。。。")
	s.msghandle.AddRouter(msg, router)
}
func (s *Server) GetMsgMaganer() iface.IConnManager {
	return s.connmanager
}
func (s *Server) RegistStartHookFunc(startfunc func(connection iface.IConnection)) {
	s.StartHookFunc = startfunc
}
func (s *Server) RegistStopHookFunc(stopfunc func(connection iface.IConnection)) {
	s.StopHookFunc = stopfunc
}

//两个钩子函数的调用
func (s *Server) CallStartHookFunc(connection iface.IConnection) {
	if s.StartHookFunc == nil {
		return
	}
	s.StartHookFunc(connection)
}
func (s *Server) CallStopHookFunc(connection iface.IConnection) {
	if s.StopHookFunc == nil {
		return
	}
	s.StopHookFunc(connection)
}

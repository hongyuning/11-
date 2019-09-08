package iface

type IConnManager interface {
AddConn (IConnection)
RemoveConn (int)
GetConn (int)IConnection
GetConnCount()int
ClearConn()
}

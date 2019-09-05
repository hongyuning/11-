package iface
type IRquest interface {
	GetConnection()IConnection
	GetMessage ()IMessage
}

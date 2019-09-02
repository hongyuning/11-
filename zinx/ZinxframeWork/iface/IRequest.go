package iface
type IRquest interface {
	GetConnection()IConnection
	GetData ()[]byte
	GetDatalen()uint32
}

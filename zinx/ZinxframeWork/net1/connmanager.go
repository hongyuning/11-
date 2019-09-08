package net1

import (
	"fmt"
	"sync"
	"zinx/V1-basic-server/zinx/ZinxframeWork/iface"
)

type ConnManager struct {
	conns    map[int]iface.IConnection
	connLock sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		conns: make(map[int]iface.IConnection),
	}
}

func (cm *ConnManager) AddConn(conn iface.IConnection) {
	cid := int(conn.GetConnID())
	cm.connLock.Lock()
	if _, ok := cm.conns[cid]; ok {
		fmt.Println("链接已经存在，无需添加....", cid)
		return
	}
	cm.conns[cid] = conn
	cm.connLock.Unlock()
	fmt.Printf("链接%d，添加成功", conn.GetConnID())
}
func (cm *ConnManager) RemoveConn(connid int) {
	cm.connLock.Lock()
	//cm.conns[connid]=nil
	delete(cm.conns, connid)
	cm.connLock.Unlock()
	fmt.Println("删除链接 ", connid)
}
func (cm *ConnManager) GetConn(connid int) iface.IConnection {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()
	conn := cm.conns[connid]
	return conn
}
func (cm *ConnManager) GetConnCount() int {
	count := len(cm.conns)
	return count
}

//清除所有链接
func (cm *ConnManager) ClearConn() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	for cid, conn := range cm.conns {
		conn.Stop()
		cm.conns[cid] = nil
	}
}

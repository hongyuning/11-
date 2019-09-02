package net1

import (
	"bytes"
	"encoding/binary"
	"zinx/V1-basic-server/zinx/ZinxframeWork/iface"
)

type DataPack struct {
}

//封包函数
func (dp *DataPack)Pack(msg iface.IMessage)([]byte,error){
data:=msg.GetData()
dataLen:=msg.GetDataLen()
msgid:=msg.GetMsgID()
var buff  bytes.Buffer
//写消息头，长度
err:=binary.Write(&buff,binary.LittleEndian,dataLen)
if err!=nil{
	return  nil,err
}
//读消息头，写msgid
err=binary.Write(&buff,binary.LittleEndian,msgid)
	if err!=nil{
		return  nil,err
	}
//3.写消息体
err=binary.Write(&buff,binary.LittleEndian,data)
	if err!=nil{
		return  nil,err
	}
return buff.Bytes(),nil

}
//解包数据  进来的data就是8字节
func (dp *DataPack)Unpack(data []byte)(iface.IMessage,error){
reader:=bytes.NewReader(data)
var message Message
	err:=binary.Read(reader,binary.LittleEndian,&message.len)
	if err!=nil{
		return  nil,err
	}
	err=binary.Read(reader,binary.LittleEndian,&message.msgid)
	if err!=nil{
		return  nil,err
	}
	return &message,nil

}
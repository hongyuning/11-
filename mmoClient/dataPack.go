package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Message struct {
	len   uint32
	msgid uint32
	data  []byte
}
func NewMessage (msgid uint32,data []byte)*Message{
	return &Message{
		len:   uint32(len(data)),
		msgid: msgid,
		data:  data,
	}
}

type datapack struct {
}

func NewDataPack() *datapack {
	return &datapack{}
}

func (dp *datapack) UnDataPack(data []byte) (*Message, error) {
   fmt.Println("开始解包出数据头.....")
   var message  Message
   reader:=bytes.NewReader(data)
   err:=binary.Read(reader,binary.LittleEndian,&message.len)
   if err!=nil{
   	fmt.Println("message  length err",err)
   	return nil,err
   }
   err=binary.Read(reader,binary.LittleEndian,&message.msgid)
	if err!=nil{
		fmt.Println("message  msgid err",err)
		return nil,err
	}
   return  &message,nil
}

func (dp *datapack) Pack(message *Message) ([]byte, error) {
	var buff bytes.Buffer
	//1. 写入数据长度
	if err := binary.Write(&buff, binary.LittleEndian, &message.len); err != nil {
		fmt.Println("binary write len err:", err)
		return nil, err
	}
	//2. 写入数据类型
	if err := binary.Write(&buff, binary.LittleEndian, &message.msgid); err != nil {
		fmt.Println("binary write msgid err:", err)
		return nil, err
	}
	//3. 写入数据
	if err := binary.Write(&buff, binary.LittleEndian, &message.data); err != nil {
		fmt.Println("binary write data err:", err)
		return nil, err
	}

	return buff.Bytes(), nil
}

package  main

import (
	"fmt"
	"net"
	"time"
)

func main(){
conn,err:=net.Dial("tcp4",":8848")
if err!=nil{
	fmt.Println("conn err",err )
}
data:=[]byte("hello ")
	for  {
		writecount,err:=conn.Write(data)
		if err!=nil{
			fmt.Println("writecount err",err )
		}
		fmt.Println(writecount)
		buf:=make([]byte,4096)
		readcount,err:=conn.Read(buf)
		if err!=nil{
			fmt.Println("readcount err",err )
		}
		fmt.Println("Server====》client",readcount,string(buf[:readcount]))
		time.Sleep(1*time.Second)
	}

}
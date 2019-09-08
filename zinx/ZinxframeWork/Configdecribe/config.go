package Configdecribe

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	Ip          string
	Port        uint32
	Name        string
	Version     string
	WorkSize    int  //工作池的数量
	TashqueSize int//接收请求的数量
	MustConnCount int
}

func init() {
	err := LoadConfig()
	if err != nil {
		fmt.Println("加载配置文件失败", err)
		os.Exit(-1)
	}
	fmt.Printf("配置文件信息为：%v\n", GlobalConfig)
}

var GlobalConfig Config

func LoadConfig() error {
	configInfo, err := ioutil.ReadFile("./config/config.json")
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = json.Unmarshal(configInfo, &GlobalConfig)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("读取文件成功")
	return nil
}

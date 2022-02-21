package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Info struct {
	Data []MpInfo `json:"mpInfo"`
}

type MpInfo struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Icon    string `json:"icon"`
}

func main() {
	// 打开文件
	file, _ := os.Open("/Users/nateyang/Documents/GitHub/go-pioneer/conf.json")

	// 关闭文件
	defer file.Close()

	//NewDecoder创建一个从file读取并解码json对象的*Decoder，解码器有自己的缓冲，并可能超前读取部分json数据。
	decoder := json.NewDecoder(file)

	conf := Info{}
	//Decode从输入流读取下一个json编码值并保存在v指向的值里
	err := decoder.Decode(&conf)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("name:" + conf.Data[0].Name)
}

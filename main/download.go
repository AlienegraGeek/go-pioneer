package main

import (
	"AlienegraGeek/go-pioneer/util"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
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

	// 设置路由，如果访问/，则调用index方法
	http.HandleFunc("/index", mock)
	http.HandleFunc("/download", fileDownload)
	http.HandleFunc("/confDownload", confDownload)
	http.HandleFunc("/getMpInfo", getMiniProgramInfo)

	// 启动web服务，监听9090端口
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func fileDownload(w http.ResponseWriter, r *http.Request) {
	//filename := get_filename_from_request(r)
	//filepath := "/Users/nateyang/Documents/GitHub/go-pioneer/script/app.js"
	filepath := "/Users/nateyang/Documents/GitHub/go-pioneer/script/icon_coin_fil.png"

	file, _ := os.Open(filepath)
	defer file.Close()

	//fileHeader := make([]byte, 512)
	//file.Read(fileHeader)

	fileStat, _ := file.Stat()

	w.Header().Set("Content-Disposition", "attachment; filepath="+filepath)
	//w.Header().Set("Content-Type", http.DetectContentType(fileHeader))
	w.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))

	_, _ = file.Seek(0, 0)
	i, err := io.Copy(w, file)
	fmt.Println("io.Copy:", i, err)

	return
}

func confDownload(w http.ResponseWriter, r *http.Request) {
	//filename := get_filename_from_request(r)
	//filepath := "/Users/nateyang/Documents/GitHub/go-pioneer/script/app.js"
	filepath := "/Users/nateyang/Documents/GitHub/go-pioneer/conf.json"

	file, _ := os.Open(filepath)
	defer file.Close()

	//fileHeader := make([]byte, 512)
	//file.Read(fileHeader)

	fileStat, _ := file.Stat()

	w.Header().Set("Content-Disposition", "attachment; filepath="+filepath)
	//w.Header().Set("Content-Type", http.DetectContentType(fileHeader))
	w.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))

	_, _ = file.Seek(0, 0)
	i, err := io.Copy(w, file)
	fmt.Println("io.Copy:", i, err)

	return
}

func getMiniProgramInfo(w http.ResponseWriter, r *http.Request) {
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
	resp, _ := json.Marshal(conf)
	w.Write(resp)
}

func mock(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "http mock")
}

func NewStatusOKJson(data interface{}) util.JSONResponse {
	result := fmt.Sprintf("%+v", data)
	fmt.Printf("\x1b[%dm[api-return] return : %v \x1b[0m\n", 32, result)
	return util.JSONResponse{
		Code: http.StatusOK,
		JSON: data,
	}
}

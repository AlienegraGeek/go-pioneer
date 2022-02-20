package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {

	// 设置路由，如果访问/，则调用index方法
	http.HandleFunc("/index", mock)
	http.HandleFunc("/download", fileDownload)

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

func mock(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "http mock")
}

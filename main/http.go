package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	//单独写回调函数
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/post", postHandler)
	// addr：监听的地址
	// handler：回调函数
	http.ListenAndServe("192.168.2.142:2040", nil)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data := r.URL.Query()
	fmt.Println(data.Get("name"))
	fmt.Println(data.Get("age"))
	answer := `{"status": "ok"}`
	w.Write([]byte(answer))
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	//1. 请求类型是aplication/x-www-form-urlencode时解析form数据
	fmt.Println(r.Body)
	r.ParseForm()
	fmt.Println(r.PostForm) //打印form数数据
	fmt.Println(r.PostForm.Get("name"), r.PostForm.Get("age"))
	//2. 请求类型是application/json时从r.Body读取数据
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("read request.Body failed, err", err)
		return
	}
	fmt.Println(string(b))
	answer := `{"status":"ok"}`
	w.Write([]byte(answer))
}

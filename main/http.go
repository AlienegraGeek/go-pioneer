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
	http.HandleFunc("/web/login", loginHandler)
	http.HandleFunc("/zee/test", zeeTestHandler)
	// addr：监听的地址
	// handler：回调函数
	http.ListenAndServe("192.168.2.142:9000", nil)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data := r.URL.Query()
	fmt.Println(data.Get("accNo"))
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

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	defer r.Body.Close()
	//1. 请求类型是aplication/x-www-form-urlencode时解析form数据
	fmt.Println(r.Body)
	fmt.Println(r)
	r.ParseForm()
	fmt.Println(r.PostForm) //打印form数数据
	fmt.Println(r.PostForm.Get("accNo"), r.PostForm.Get("password"))
	//2. 请求类型是application/json时从r.Body读取数据
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("read request.Body failed, err", err)
		return
	}
	fmt.Println(string(b))
	answer := `{"code":"0","token":"123456789"}`
	w.Write([]byte(answer))
}

func zeeTestHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	//1. 请求类型是aplication/x-www-form-urlencode时解析form数据
	fmt.Println(r.Body)
	//r.ParseForm()
	//fmt.Println(r.PostForm) //打印form数数据
	//fmt.Println(r.PostForm.Get("deviceSerial"))
	//2. 请求类型是application/json时从r.Body读取数据
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("read request.Body failed, err", err)
		return
	}
	fmt.Println(string(b))
	//answer := `{"data":{"code":"0","msg":"success"}}`
	answer := `{"code":"0","token":"123456789"}`
	//	answer := `{
	//    "data": {
	//        "appId": "c672899cc9874cc6982f50ecb94a812c"
	//    },
	//    "code": "0",
	//    "msg": "操作成功"
	//}`
	w.Write([]byte(answer))
}

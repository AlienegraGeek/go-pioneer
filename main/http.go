package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"time"
)

var CHARS = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}

const (
	TOKEN = "YOUR_TOKEN" // 替换成您在公众平台后台设置的Token值
)

func main() {
	//open.InitGPT()
	//单独写回调函数
	http.HandleFunc("/", handler)
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/post", postHandler)
	http.HandleFunc("/web/login", loginHandler)
	http.HandleFunc("/zee/test", zeeTestHandler)
	http.HandleFunc("/mage/test", mageTestHandler)
	http.HandleFunc("/wx/test", wxTestHandler)
	// addr：监听的地址
	// handler：回调函数
	http.ListenAndServe(":2040", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
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

type Res struct {
	Code string `json:"code"`
	Res  string `json:"res"`
}

func mageTestHandler(w http.ResponseWriter, r *http.Request) {
	randStr := GetRoundName(1)
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
	//answer, _ := json.Marshal(rs)
	answer, _ := json.Marshal(Res{Code: "0", Res: randStr})
	w.Write(answer)
}

func wxTestHandler(w http.ResponseWriter, r *http.Request) {
	sign := r.URL.Query().Get("signature")
	ts := r.URL.Query().Get("timestamp")
	nonce := r.URL.Query().Get("nonce")
	echostr := r.URL.Query().Get("echostr")
	if checkSignature(sign, ts, nonce) {
		// 如果请求来自微信服务器，则原样返回echostr参数值
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprintf(w, "%s", echostr)
		if err != nil {
			fmt.Println("write error:", err)
			return
		}
		return
	}
	// 如果请求不是来自微信服务器，则返回错误信息
	http.Error(w, "Bad Request", http.StatusBadRequest)
}

// 检查请求是否来自微信服务器
func checkSignature(signature, timestamp, nonce string) bool {
	s := []string{TOKEN, timestamp, nonce}
	sort.Strings(s)
	sha1String := sha1.New()
	sha1String.Write([]byte(strings.Join(s, "")))
	hash := fmt.Sprintf("%x", sha1String.Sum(nil))
	return signature == hash
}

func GetRoundName(size int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < size; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

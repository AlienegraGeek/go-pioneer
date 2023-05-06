package main

import (
	"AlienegraGeek/go-pioneer/open"
	"AlienegraGeek/go-pioneer/util"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

// TextMessage 代表文本消息
type TextMessage struct {
	Content string `json:"content"`
}

// ChatMessage 代表群聊消息
type ChatMessage struct {
	ChatID  string `json:"chatid"`
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

// WechatRequest 代表微信请求
type WechatRequest struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   int64  `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Content      string `xml:"Content"`
	Event        string `xml:"Event"`
	EventKey     string `xml:"EventKey"`
	MsgID        int64  `xml:"MsgId"`
}

// WechatResponse 代表微信响应
type WechatResponse struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   int64  `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Content      string `xml:"Content,omitempty"`
}

const (
	TOKEN     = "EntySquare666"                    // 替换成您在公众平台后台设置的Token值
	APPSECRET = "5e4ac95e08aea15bafab3cce7a5598ae" // 替换为您的微信公众号AppID
	APPID     = "wx12ab0d88a0fa7f22"               // 替换为您的微信公众号AppSecret
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data := r.URL.Query()
	fmt.Println(data.Get("accNo"))
	fmt.Println(data.Get("age"))
	answer := `{"status": "ok"}`
	w.Write([]byte(answer))
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
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

func MageTestHandler(w http.ResponseWriter, r *http.Request) {
	randStr := util.GetRandomName(1)
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
	answer, _ := json.Marshal(randStr)
	w.Write(answer)
}

func WxTestHandler(w http.ResponseWriter, r *http.Request) {
	sign := r.URL.Query().Get("signature")
	ts := r.URL.Query().Get("timestamp")
	nonce := r.URL.Query().Get("nonce")
	echostr := r.URL.Query().Get("echostr")
	if checkSignatureTest(sign, ts, nonce) {
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

// 处理微信请求
func HandleWechat(w http.ResponseWriter, r *http.Request) {
	if !checkSignature(w, r) {
		fmt.Println("signature fail")
		http.Error(w, "signature fail", http.StatusInternalServerError)
		return
	}
	// 解析微信请求的 XML 数据
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Failed to read request body:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	var req WechatRequest
	err = xml.Unmarshal(reqBody, &req)
	if err != nil {
		fmt.Println("Failed to unmarshal request body:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// 根据请求内容进行相应的处理
	if req.MsgType == "text" && req.Content == "hello" {
		// 如果请求是文本消息，并且内容为"hello"，则回复一个相同的文本消息
		resp := WechatResponse{
			ToUserName:   req.FromUserName,
			FromUserName: req.ToUserName,
			CreateTime:   req.CreateTime,
			MsgType:      "text",
			Content:      "hello",
		}
		fmt.Println("user req content hello:", resp.Content)
		writeResponse(w, resp)
	} else if req.MsgType == "text" {
		done := make(chan bool)
		//open.InitGPT("一句话简单介绍一些golang")
		var gRes = "waiting..."
		go func() {
			gRes, err = open.FetchGPT(req.Content)
			if err != nil {
				fmt.Println("gpt res error:", err)
				http.Error(w, "gpt res error", http.StatusInternalServerError)
				return
			}
			done <- true
		}()
		<-done
		// 如果请求是文本消息，并且不是"hello"，则回复一个相同的文本消息
		resp := WechatResponse{
			ToUserName:   req.FromUserName,
			FromUserName: req.ToUserName,
			CreateTime:   req.CreateTime,
			MsgType:      "text",
			Content:      gRes,
		}
		fmt.Println("gpt response done:", gRes)
		writeResponse(w, resp)
	} else {
		respXML, _ := xml.Marshal("success")
		w.Write(respXML)
	}
}

// writeResponse 用于将响应写入 http.ResponseWriter 中
func writeResponse(w http.ResponseWriter, resp WechatResponse) {
	//respXML, err := xml.MarshalIndent(resp, "", "  ")
	respXML, err := xml.Marshal(resp)
	if err != nil {
		fmt.Println("Failed to marshal response body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Write(respXML)
}

// 检查请求是否来自微信服务器
func checkSignatureTest(signature, timestamp, nonce string) bool {
	s := []string{TOKEN, timestamp, nonce}
	sort.Strings(s)
	sha1String := sha1.New()
	sha1String.Write([]byte(strings.Join(s, "")))
	hash := fmt.Sprintf("%x", sha1String.Sum(nil))
	return signature == hash
}

// CheckSignature 用于校验微信服务器发送的请求签名
func checkSignature(w http.ResponseWriter, r *http.Request) bool {
	signature := r.URL.Query().Get("signature")
	timestamp := r.URL.Query().Get("timestamp")
	nonce := r.URL.Query().Get("nonce")
	echostr := r.URL.Query().Get("echostr")

	params := []string{TOKEN, timestamp, nonce}
	sort.Strings(params)
	hash := sha1.New()
	_, _ = io.WriteString(hash, strings.Join(params, ""))
	if hex.EncodeToString(hash.Sum(nil)) != signature {
		fmt.Fprintln(w, "Invalid signature")
		return false
	}
	fmt.Fprintln(w, echostr)
	return true
}

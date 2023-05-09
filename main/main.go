package main

import (
	"fmt"
	"net/http"
)

func main() {

	//const WECHATYTOKEN = "puppet_paimon_68a6bcfa-551d-4929-9a0a-70379d65aaa9"
	//
	//var bot = wechaty.NewWechaty(wechaty.WithPuppetOption(wp.Option{
	//	Token: WECHATYTOKEN,
	//}))
	//
	//bot.OnScan(chaty.OnScan).OnLogin(chaty.OnLogin).OnMessage(chaty.OnMessage).OnLogout(chaty.OnLogout)
	//
	//bot.DaemonStart()

	http.HandleFunc("/get", GetHandler)
	http.HandleFunc("/post", PostHandler)
	http.HandleFunc("/mage/test", MageTestHandler)
	http.HandleFunc("/wx/test", WxTestHandler)
	http.HandleFunc("/wxChat", HandleWechat)

	err := http.ListenAndServe(":2040", nil)
	if err != nil {
		fmt.Print("http listen error", err)
	}
}

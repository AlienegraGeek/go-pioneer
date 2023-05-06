package main

import (
	"AlienegraGeek/go-pioneer/chaty"
	"fmt"
	"github.com/wechaty/go-wechaty/wechaty"
	wp "github.com/wechaty/go-wechaty/wechaty-puppet"
	"github.com/wechaty/go-wechaty/wechaty/user"
	"net/http"
)

func main() {

	var bot = wechaty.NewWechaty(wechaty.WithPuppetOption(wp.Option{
		Token: "puppet_paimon_68a6bcfa-551d-4929-9a0a-70379d65aaa9",
	}))

	bot.OnScan(chaty.OnScan).OnLogin(func(ctx *wechaty.Context, user *user.ContactSelf) {
		fmt.Printf("User %s logined\n", user.Name())
	}).OnMessage(chaty.OnMessage).OnLogout(func(ctx *wechaty.Context, user *user.ContactSelf, reason string) {
		fmt.Printf("User %s logouted: %s\n", user, reason)
	})

	bot.DaemonStart()

	http.HandleFunc("/get", GetHandler)
	http.HandleFunc("/post", PostHandler)
	http.HandleFunc("/mage/test", MageTestHandler)
	http.HandleFunc("/wx/test", WxTestHandler)
	http.HandleFunc("/wxChat", HandleWechat)

	http.ListenAndServe(":2040", nil)

}

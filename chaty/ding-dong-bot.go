package chaty

import (
	"AlienegraGeek/go-pioneer/open"
	"fmt"
	"github.com/mdp/qrterminal/v3"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/wechaty/go-wechaty/wechaty"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

func we() {
	var bot = wechaty.NewWechaty()

	bot.OnScan(OnScan).OnLogin(func(ctx *wechaty.Context, user *user.ContactSelf) {
		fmt.Printf("User %s logined\n", user.Name())
	}).OnMessage(OnMessage).OnLogout(func(ctx *wechaty.Context, user *user.ContactSelf, reason string) {
		fmt.Printf("User %s logouted: %s\n", user, reason)
	})

	bot.DaemonStart()
}

func OnMessage(ctx *wechaty.Context, message *user.Message) {
	log.Println(message)

	if message.Self() {
		log.Println("Message discarded because its outgoing")
	}

	if message.Age() > 2*60*time.Second {
		log.Println("Message discarded because its TOO OLD(than 2 minutes)")
	}

	if message.Type() != schemas.MessageTypeText || (!strings.HasPrefix(message.Text(), "@Entiny") && !strings.HasPrefix(message.Text(), "@Malaka")) {
		log.Println("Message discarded because it does not match '@Entiny'")
		return
	}
	reqStr := message.Text()
	// 将字符串转换成 rune 数组
	rs := []rune(reqStr)
	// 截取第 6 个字符到最后一个字符
	n := utf8.RuneCountInString(reqStr)
	problem := string(rs[8:n])
	//	fmt.Println("problem:", problem[1])
	gRes, err := open.FetchGPT(problem)
	if err != nil {
		fmt.Println("gpt res error:", err)
		return
	}
	// 1. reply 'dong'
	_, err = message.Say(gRes)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("response:", gRes)

	// 2. reply image(qrcode image)
	//fileBox, _ := file_box.FromUrl("https://wechaty.github.io/wechaty/images/bot-qr-code.png", "", nil)
	//_, err = message.Say(fileBox)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//log.Printf("REPLY: %s\n", fileBox)
}

func OnScan(ctx *wechaty.Context, qrCode string, status schemas.ScanStatus, data string) {
	if status == schemas.ScanStatusWaiting || status == schemas.ScanStatusTimeout {
		qrterminal.GenerateHalfBlock(qrCode, qrterminal.L, os.Stdout)

		qrcodeImageUrl := fmt.Sprintf("https://wechaty.js.org/qrcode/%s", url.QueryEscape(qrCode))
		fmt.Printf("onScan: %s - %s\n", status, qrcodeImageUrl)
		return
	}
	fmt.Printf("onScan: %s\n", status)
}

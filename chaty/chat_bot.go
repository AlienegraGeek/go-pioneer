package chaty

import (
	"AlienegraGeek/go-pioneer/open"
	"fmt"
	"github.com/mdp/qrterminal/v3"
	"github.com/wechaty/go-wechaty/wechaty"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
	"log"
	"net/url"
	"os"
	"time"
)

func OnMessage(ctx *wechaty.Context, message *user.Message) {
	log.Println(message)

	if message.Self() {
		log.Println("Message discarded because its outgoing")
	}

	if message.Age() > 2*60*time.Second {
		log.Println("Message discarded because its TOO OLD(than 2 minutes)")
	}

	log.Println("from id:", message.Talker().ID())
	log.Println("from name:", message.Talker().Name())
	talker := message.Talker().Name()
	log.Println("mention text:", message.MentionText())
	log.Println("mention state:", message.MentionSelf())
	message.MentionList()

	//if message.Type() != schemas.MessageTypeText || (!strings.HasPrefix(message.Text(), "@Entiny") && !strings.HasPrefix(message.Text(), "@Malaka")) {
	//	log.Println("Message discarded because it does not match '@Entiny' & '@Malaka'")
	//	return
	//}

	if message.Type() != schemas.MessageTypeText || !message.MentionSelf() {
		log.Println("Message discarded because it does not mentioned")
		return
	}
	//reqStr := message.Text()
	// 将字符串转换成 rune 数组
	//rs := []rune(reqStr)
	// 截取第 6 个字符到最后一个字符
	//n := utf8.RuneCountInString(reqStr)
	//problem := string(rs[8:n])
	problem := message.MentionText()
	//fmt.Println("problem:", problem[1])
	var gRes = "haha"
	gRes, err := open.FetchContextGPT(talker, problem)
	if err != nil {
		fmt.Println("gpt res error:", err)
		return
	}
	// 1. reply 'dong'
	_, err = message.Say("@" + talker + "\n" + gRes)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("response:", "@"+talker+" "+gRes)

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

func OnLogin(ctx *wechaty.Context, user *user.ContactSelf) {
	fmt.Printf("User %s logined\n", user.Name())
}

func OnLogout(ctx *wechaty.Context, user *user.ContactSelf, reason string) {
	fmt.Printf("User %s logouted: %s\n", user, reason)
}

package test

import (
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestStringSplit(t *testing.T) {
	s := "@Entiny:哈哈哈"
	ss := strings.Split(s, ":")
	fmt.Println("ss:", ss[1])
}

func TestStringSub(t *testing.T) {
	s := "@Entiny 哈哈哈"
	// 将字符串转换成 rune 数组
	rs := []rune(s)
	// 截取第 6 个字符到最后一个字符
	n := utf8.RuneCountInString(s)
	s2 := string(rs[8:n])
	fmt.Println(s2) // 输出：, 你好！
}

//func TestBot(t *testing.T) {
//	en := entiny.NewEntiny()
//	bot, err := en.InitWeChatyBot("puppet_paimon_68a6bcfa-551d-4929-9a0a-70379d65aaa9")
//	if err != nil {
//		return
//	}
//	//bot.DaemonStart()
//	err = en.SendMsg(bot, "Mage", "pooluo", "panic")
//	if err != nil {
//		return
//	}
//}

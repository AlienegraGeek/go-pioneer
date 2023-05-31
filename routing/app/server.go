package app

import (
	"AlienegraGeek/go-pioneer/config"
	"AlienegraGeek/go-pioneer/routing/types"
	"AlienegraGeek/go-pioneer/util"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
	"time"
)

func HandleTest(c *fiber.Ctx) error {
	testReq := types.TestParam{}
	err := c.BodyParser(&testReq)
	if err != nil {
		return c.JSON(util.MessageResponse(config.MESSAGE_FAIL, "can not transfer request to struct", "请求参数错误"))
	}
	fmt.Println("phone = ", testReq.Phone) // 手机号
	// 模拟长轮询等待过程
	//time.Sleep(5 * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	timer := time.NewTimer(time.Duration(5) * time.Second)
	counter := 1

	//userData := types.TestRes{
	//	Name:    "Kobe",
	//	Phone:   "13666666666",
	//	Balance: 666,
	//	Card:    "1234567812345678",
	//	Date:    "2023/05/31",
	//}

	userData := types.TestRes{}
	content, err := os.ReadFile("conf.json")
	if err != nil {
		fmt.Printf("err", err)
	}
	err = json.Unmarshal(content, &userData)
	if err != nil {
		fmt.Println("解析 JSON 失败：", err)
	}
	fmt.Println(userData)
	for {
		select {
		case <-ticker.C:
			fmt.Printf("当前是第 %d 秒\n", counter)
			counter++

		case <-timer.C:
			fmt.Println("定时器到达，结束打印")
			// 检查是否有新数据，根据实际情况进行判断
			if hasNewData() {
				// 返回数据给客户端
				//randStr := util.GetRandomName(6)
				return c.JSON(util.SuccessResponse(userData))
			}
		}
	}
}

func hasNewData() bool {
	// 在这里实现检查是否有新数据的逻辑，返回true或false
	// 根据实际需求和业务逻辑来判断是否有新数据
	return true
}

package app

import (
	"github.com/gofiber/fiber/v2"
	"go-pioneer/config"
	"go-pioneer/min"
	"go-pioneer/util"
)

type MinParam struct {
	Bucket     string `json:"bucket"`
	ObjectName string `json:"object_name"`
}

func GetPreSignedUrl(c *fiber.Ctx) error {
	//testReq := MinParam{}
	//err := c.BodyParser(&testReq)
	//if err != nil {
	//	return c.JSON(util.MessageResponse(config.MESSAGE_FAIL, "can not transfer request to struct", "请求参数错误"))
	//}
	//bucketName := "my-bucket"
	bucketName := c.Query("bucket", "")
	objectName := c.Query("object_name", "")
	client := min.GetInstance()
	preSignedURL, err := min.UploadPreSigned(client, bucketName, objectName)
	if err != nil {
		return c.JSON(util.MessageResponse(config.MESSAGE_FAIL, "", "上传错误"))
	}
	return c.JSON(fiber.Map{
		"url": preSignedURL,
	})
}

func DownloadPreSignedUrl(c *fiber.Ctx) error {
	//testReq := MinParam{}
	//err := c.BodyParser(&testReq)
	//if err != nil {
	//	return c.JSON(util.MessageResponse(config.MESSAGE_FAIL, "can not transfer request to struct", "请求参数错误"))
	//}
	//bucketName := "my-bucket"
	bucketName := c.Query("bucket", "")
	objectName := c.Query("object_name", "")
	client := min.GetInstance()
	preSignedURL, err := min.DownloadPreSigned(client, bucketName, objectName)
	if err != nil {
		return c.JSON(util.MessageResponse(config.MESSAGE_FAIL, "", "上传错误"))
	}
	return c.JSON(fiber.Map{
		"url": preSignedURL,
	})
}

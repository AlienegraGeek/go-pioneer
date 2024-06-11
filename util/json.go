package util

import (
	"github.com/sirupsen/logrus"
	"go-pioneer/log"
)

const (
	CodeOk       = 0  // 成功
	CodeErr      = -1 // 失败
	CodeErrToken = -2 // token相关的异常
	//CodeReject   = "2" // 拒绝
	//CodeTimeout  = "3" // 超时
)

// DataResponse represents an HTTP response which contains a JSON body.
type DataResponse struct {
	// HTTP status code.
	Code int `json:"code"`
	// JSON represents the JSON that should be serialized and sent to the client
	Data interface{} `json:"data"`
}

func SuccessResponse(data interface{}) DataResponse {
	return DataResponse{
		Code: 0,
		Data: data,
	}
}

// MessageResponse returns a JSONResponse with a 'message' key containing the given text.
func MessageResponse(code int, msg, msgZh string) DataResponse {
	log.Log.WithFields(logrus.Fields{
		"code":   code,
		"msg_zh": msgZh,
	}).Warnf(msg)
	return DataResponse{
		Code: code,
		Data: struct {
			Message   string `json:"message"`
			MessageZh string `json:"message_zh"`
		}{msg, msgZh},
	}
}

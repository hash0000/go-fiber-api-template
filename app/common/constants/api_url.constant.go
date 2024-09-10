package constants

import (
	"os"
)

type urlType struct {
	ApiUrl                 string
	UploadTaleTg           string
	UploadTgAnalytic       string
	UploadErrorTg          string
	FabulaCreateText2Image string
	FabulaGetOrder         string
	ChatGptRequest         string
}

var Url urlType

func Load() {
	Url = urlType{
		ApiUrl:                 os.Getenv("API_URL"),
		UploadTaleTg:           os.Getenv("TG_BOT_UPLOAD_URL"),
		UploadTgAnalytic:       os.Getenv("TG_BOT_UPLOAD_ANALYTIC"),
		UploadErrorTg:          os.Getenv("TG_BOT_ERROR_URL"),
		FabulaCreateText2Image: "https://integration-api.fabula-app.com/order/createTextToImageOrder",
		FabulaGetOrder:         "https://integration-api.fabula-app.com/order/getUserImages?order_id=%d",
		ChatGptRequest:         "https://api.openai.com/v1/chat/completions",
	}
}

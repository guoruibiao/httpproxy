package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guoruibiao/gorequests"
	"io"
	"net/http"
)


func Post(ctx *gin.Context) {
	// Method 代理与源请求保持一致
	method := ctx.Request.Method
	url    := ctx.PostForm("url")

	// header 透传
	headers := make(map[string]string)
	for key, value := range ctx.Request.Header {
		headers[key] = value[0]
	}

	// 如果是 post 请求，body 也要透传
	payload := make(map[string]string)
	for key, value := range ctx.Request.PostForm {
		payload[key] = value[0]
	}
	fmt.Println(payload["appid"])
	fmt.Println(payload["token"])


	response, err := gorequests.NewRequest(method, url).Headers(headers).Form(payload).DoRequest()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}else {
		content, _ := response.Content()
		fmt.Println(content)
		io.WriteString(ctx.Writer, content)
	}
}

func main() {
	app := gin.Default()

	app.POST("/proxy", Post)


	// 启动服务进行监听
	_ = app.Run(":80")
}
package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guoruibiao/gorequests"
	"io"
	"net/http"
	"os"
	"strconv"
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


	var port = flag.Int("port", -1, "the port of the service")
	flag.Parse()
	if *port <= 0 {
		fmt.Println("WRONG PORT WITH ", *port)
		os.Exit(-1)
	}

	address := ":" + strconv.Itoa(*port)
	// 启动服务进行监听
	_ = app.Run(address)
}
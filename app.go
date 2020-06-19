package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/guoruibiao/gorequests"
	"github.com/guoruibiao/httpproxy/utils"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)


func Post(ctx *gin.Context) {
	// Method 代理与源请求保持一致
	method := strings.ToUpper(ctx.PostForm("method"))
	url    := ctx.PostForm("url")

	// header 透传
	headers := make(map[string]string)
	for key, value := range ctx.Request.Header {
		headers[key] = value[0]
	}
	/*
	 * IP 相关处理
	 * 1. 没有使用代理服务器
	 *    REMOTE_ADDR            真实 IP
	 *    HTTP_VIA               没有数值或者不显示
	 *    HTTP_X_FORWARDED_FOR   没数值或者不显示
	 *
	 * 2. 使用透明代理 √
	 *    REMOTE_ADDR            最后一个代理服务器的 IP
	 *    HTTP_VIA               代理服务器 IP
	 *    HTTP_X_FORWARDED_FOR   真实 IP，经过多个代理服务器时，尾部追加，无法"隐身"
	 *
	 * 3. 使用普通匿名代理
	 *    REMOTE_ADDR            最后一个代理服务器 IP
	 *    HTTP_VIA               代理服务器 IP
	 *    HTTP_X_FORWARDED_FOR   代理 IP，经过多个代理服务器时，尾部追加，无法"隐身"
	 *
	 * 4. 使用欺骗性代理
	 *    REMOTE_ADDR            代理服务器 IP
	 *    HTTP_VIA               代理服务器 IP
	 *    HTTP_X_FORWARDED_FOR   随机 IP，经过多个代理服务器时，尾部追加，可以"隐身"
	 *
	 * 5. 使用高匿名代理
	 *    REMOTE_ADDR            代理服务器 IP
	 *    HTTP_VIA               没数值或者不显示
	 *    HTTP_X_FORWARDED_FOR   没数值或者不显示
	 *
	 */
	remoteAddr, _ := utils.GetExternalIP()
	headers["REMOTE_ADDR"]          = remoteAddr.String()
	headers["HTTP_VIA"]             = remoteAddr.String()
	if xForward, exist := headers["HTTP_X_FORWARDED_FOR"]; exist && xForward != ""{
		prefix := ctx.Request.Header["HTTP_X_FORWARDED_FOR"][0]
		headers["HTTP_X_FORWARDED_FOR"] = fmt.Sprintf("%s,%s", prefix, ctx.Request.RemoteAddr)
	}else{
		// 本机测试 [::1]:12345
		if strings.HasPrefix(ctx.Request.RemoteAddr, "[::1]") {
			headers["HTTP_X_FORWARDED_FOR"] = remoteAddr.String()
		}else {
			headers["HTTP_X_FORWARDED_FOR"] = strings.Split(ctx.Request.RemoteAddr, ":")[0]
		}
	}

	// 如果是 post 请求，body 也要透传
	formData := make(map[string]string)
	for key, value := range ctx.Request.PostForm {
		formData[key] = value[0]
	}

	// 如果 post 的数据是一个嵌套层次复杂的 json 串，响应的 content-type 需要转为 application/json
	// TODO 待支持 暂时好像也用不到
	/*
	var payload map[string]interface{}
	if bytes, err := ioutil.ReadAll(ctx.Request.Body); err == nil {
		if err = json.Unmarshal(bytes, &payload); err == nil {
			for key, value := range payload {
				payload[key] = value
			}
		}
	}
	 */

	response, err := gorequests.NewRequest(method, url).Headers(headers).Form(formData).DoRequest()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}else {
		content, _ := response.Content()
		io.WriteString(ctx.Writer, content)
	}

	// 打日志
	logline := fmt.Sprintf("parameters=%+v\n---------------------------\n, response=%v \n", formData, response)
	fmt.Println(logline)
}

func main() {
	app := gin.Default()

	app.POST("/proxy", Post)


	var port = flag.Int("port", -1, "the port of the service")
	flag.Parse()
	// TODO 上线删除
	*port = 80
	if *port <= 0 {
		fmt.Println("WRONG PORT WITH ", *port)
		os.Exit(-1)
	}

	address := ":" + strconv.Itoa(*port)
	// 启动服务进行监听
	_ = app.Run(address)
}
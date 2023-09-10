package netcall

import (
	"bytes"
	"encoding/json"
	"file-push/log"
	"fmt"
	"io"
	"net/http"
)

func HttpPostJson(params string, url string) (string, error) {

	fmt.Println(params)

	reader := bytes.NewReader([]byte(params))

	request, err := http.NewRequest("POST", url, reader)
	defer request.Body.Close() //程序在使用完回复后必须关闭回复的主体
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	//必须设定该参数,POST参数才能正常提交，意思是以json串提交数据

	client := http.Client{}
	log.Infof("do quest url %s,params %s", url, params)
	resp, err := client.Do(request) //Do 方法发送请求，返回 HTTP 回复
	if err != nil {
		log.Errorf("response error %v", err)
		return "", err
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	//byte数组直接转成string，优化内存
	str := string(respBytes)
	return str, err
}
func FtpReqCallbackIfNecessary(message ResponseVo, callbackUrl string) (string, error) {
	if callbackUrl != "" {
		params, _ := json.Marshal(message)
		postJson, err := HttpPostJson(string(params), callbackUrl)
		return postJson, err
	}
	log.Infof("callback url is null, do not need callback...")
	return "", nil
}

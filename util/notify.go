package util

import (
	"fmt"
	"net/http"
	"io/ioutil"
	url2 "net/url"
)

const(
	FTQQ_SEND_KEY = "2979-476b0312e6118c0e6d04d52ddfd4ec97"
	FTQQ_PUSH_BEAR_API = "https://pushbear.ftqq.com/sub"

	BARK_PUSH_URL = "http://202.120.117.179"
	BARK_PUSH_KEY = "PFtWy2EWyejHT5vmsUFNCJ"

	TOPIC = "AneMQ"
)

const (
	WECHAT_CHANNEL int = iota
	EMAIL_CHANNEL
	BARK_CHANNEL
)

var queueChan = make(chan string,100)
var taskQueue = make(chan func(),1000)

func init(){
	// 初始化处理队列
	go NotifyHandlerQueue()
}

func WechatNotify(notify string) bool{
	taskQueue <- func() {
		sendWechat(notify)
	}
	return true
}


func MailNotify(notify string) bool {
	taskQueue <- func() {
		sendMail(notify)
	}
	return true
}

func BarkNotify(notify string) bool{
	taskQueue <- func() {
		sendBark(notify)
	}
	return true
}

func NotifyHandlerQueue() {

	for sendFunc := range taskQueue {
		sendFunc()
	}
}

func sendWechat(notify string) bool{
	desp := notify
	buildUrl := fmt.Sprintf("%s?sendkey=%s&text=%s&desp=%s",FTQQ_PUSH_BEAR_API,FTQQ_SEND_KEY,TOPIC,desp)
	return httpGet(buildUrl)
}

func sendMail(notify string) bool{
	return true
}


func sendBark(notify string) bool{
	buildUrl := fmt.Sprintf("%s/%s/%s/%s",BARK_PUSH_URL,BARK_PUSH_KEY,TOPIC,notify)
	return httpGet(buildUrl)
}

/**
	工具函数
 */
func httpGet(buildUrl string) bool{
	url, _:= url2.Parse(buildUrl)
	url.RawQuery = url.Query().Encode()
	fmt.Println(url.String())
	resp, err := http.Get(url.String())

	if err!=nil {
		// error
		return false
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err!=nil {
		//error
		return false
	}
	return true
}

func TestSendBark(notifyInfo string){
	sendBark(notifyInfo)
}

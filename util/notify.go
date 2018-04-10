package util

import (
	"fmt"
	"net/http"
	"io/ioutil"
	url2 "net/url"
)

const(
	SEND_KEY = "2979-476b0312e6118c0e6d04d52ddfd4ec97"
	PUSH_BEAR_API = "https://pushbear.ftqq.com/sub"
	TOPIC = "AneMQ"

	WECHAT_CHANNEL int = iota
	EMAIL_CHANNEL
)

var queueChan = make(chan string,100)

var taskQueue = make(chan func(),1000)
func WechatNotify(notify string) bool{

	//queueKey := fmt.Sprintf("%s:%s",strconv.Itoa(WECHAT_CHANNEL),notify)
	//queueChan <- queueKey

	taskQueue <- func() {
		sendWechat(notify)
	}
	return true
}


func MailNotify(notify string) bool {
	//queueKey := fmt.Sprintf("%s:%s",strconv.Itoa(EMAIL_CHANNEL),notify)
	//queueChan <- queueKey

	taskQueue <- func() {
		sendMail(notify)
	}
	return true
}

func NotifyHandlerQueue() {
	//for key := range queueChan {
	//	if index := strings.LastIndex(key,":"); index!=-1 {
	//		tag := key[0:index]
	//		notify := key[index:len(key)]
	//		v,_ :=strconv.Atoi(tag);
	//		if v == WECHAT_CHANNEL {
	//			sendWechat(notify)
	//			continue
	//		}
	//
	//		if v==EMAIL_CHANNEL {
	//			sendMail(notify)
	//		}
	//	}
	//}

	for sendFunc := range taskQueue {
		sendFunc()
	}
}

func sendWechat(notify string) bool{
	desp := notify
	buildUrl := fmt.Sprintf("%s?sendkey=%s&text=%s&desp=%s",PUSH_BEAR_API,SEND_KEY,TOPIC,desp)
	url, _:= url2.Parse(buildUrl)
	url.RawQuery = url.Query().Encode()
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
	//fmt.Println(string(body))
	return true
}

func sendMail(notify string) bool{
	return true
}
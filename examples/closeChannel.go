package main

import (
	"fmt"
	"time"
)

var channel chan string

var done chan bool

var close chan bool

func init(){
	channel = make(chan string,1000)
	done = make(chan bool,1)
	close = make(chan bool,1)

}
func main(){
	go send()
	go handler()
	<- close
}

func handler() {
	for {
		select {
		case info := <-channel:
			fmt.Println(info)
		case <-done :
			close <- false
			return
		}
	}
}

func send() {
	for i:=1;i<10 ;i++  {
		channel <- fmt.Sprintf("%s-%d","hello",i)
		time.Sleep(time.Second)
	}
	done <- false
}

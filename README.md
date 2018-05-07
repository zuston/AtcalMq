# AtcalMq
the message queue for ane trace big data, which serves for the maching learning prediction

# Introduction
此为ane消息队列处理框架，所以主要分为consumer和producer两个守护进程，当然两者是分开执行的。编译而成的主要为 main 和 push 二进制程序。

# Compile
编译linux下的可执行程序执行如下命令：  
### 消费者  
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go  
### 生产者  
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build push.go  
### 队列控制台
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build console.go

# Install
1. __如果对代码进行修改，需要配置go的开发环境__
2. __也可以下载编译好的运行文件，直接运行二进制程序。__

# Usage
1. main 二进制文件为 ane 消息队列的消费者，主要处理接收的消息的存储和备份操作。执行 ./main -h 可以查看到多种参数配置。主要有两个参数需要注意。   
__执行 ./main -show true 时候，只是处理给 console 程序提供 rpc 接口，不进行消息的消费操作.__  
__执行 ./main -queueName "{queueName}" 的时候，是给消费者预定一定消费队列名，默认是全选。__  
__执行 ./main 的时候，是一切默认，全部消费队列注册，开始消费。__
2. 队列控制台 ./console 程序只是提供给服务器端的消息队列的积压情况的可视化展示，展示分钟级别的处理能力。必须在 main 程序执行之后执行，因为需要接受 main 程序的 rpc 消息。
3. push 程序为推送程序，是一个定时执行的守护程序。根据参数，可以设定定时的间隔时间和model映射的文件位置，还有时间参数（用以重发在某个时间点之后的消息）

# Config
本框架主要有两个地方的配置文件，且自带默认地址。
1. 关于消息队列的配置问题，生产环境位置为 /opt/mq.cfg, 测试环境位置即为当前代码的根目录。在此代码中不含有，已被 .gitignore 了。具体配置请参见服务器的 /opt 下的位置。
2. 关于 push 信息的 map 配置问题。环境位置均为 /opt 下，默认名字均可在push代码中自定义。在此代码中可见，格式类似于 optional.model文件中。

# Architecture
`golang`新手，对`golang`的编程范式还不够。所以设计方面还需要提升。主要讲一个`consumer`这一侧的`rabbitmq`的消费调用。  

```go
func main(){
  // 开启一个rabbitmq的消费者实例
  cf, err := rabbitmq.NewConsumerFactory(mq_uri,exchange,exchange_type)
  // 向实例注册某个消息队列，且写入处理方法的回调
  cf.Register("ane_test",testHandler)
  // 注册多个消息队列，且用同一个回调方法
  cf.RegisterAll([]string{"ane_1","ane_2"},testHandler)
  // 开启处理模式
  go cf.Handle()
  // 开始监视器rpc,主要为 console 程序提供可视化可以不加
  rabbitmq.NewWatcher()
  
  func testHandler(queueName string, msgChan <-chan amqp.Delivery){
    for msg := range msgChan{}
  }
}
```
__就是这么简单和清晰__  
```go
// 生产者
func main(){
  pf, _ := rabbitmq.NewProducerFactory(mq_uri,exchange,exchange_type,false)
	 go pf.Handle()
  pf.Publish("ane_push","hello world")
}
```

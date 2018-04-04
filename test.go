package main

import (
	"net/rpc"
	"net/http"
	"log"
	"net"
	"time"
)



type Arith int

func (t *Arith) Multiply(args int, reply *([]string)) error {
	*reply = append(*reply, "test")
	return nil
}

func main() {
	arith := new(Arith)

	rpc.Register(arith)
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
	time.Sleep(5 * time.Second)

	client, err := rpc.DialHTTP("tcp", "127.0.0.1" + ":1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}


	reply := make([]string, 10)
	err = client.Call("Arith.Multiply", 1, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	log.Println(reply)
}
package main

import (
	"github.com/tsuna/gohbase"
	"github.com/tsuna/gohbase/hrpc"
	"context"
	"time"
	"fmt"
	"math/rand"
	"github.com/ivpusic/grpool"
)

var hbaseClient gohbase.Client

func init()  {
	hbaseClient = gohbase.NewClient("slave4,slave2,slave3")
}

func main(){

	pool := grpool.NewPool(10, 1000)
	defer pool.Release()

	go func() {
		for i:=0;i<500 ;i++  {
			pool.JobQueue <- func() {
				basicInfoCfMapper := map[string]map[string][]byte{"basic":{"uid":[]byte("123")}}
				rowkey := fmt.Sprintf("%d%d%d",rand.Intn(100),time.Now().Unix(),rand.Intn(100))
				biPutRequest, err := hrpc.NewPutStr(context.Background(),"testTable",rowkey,basicInfoCfMapper)
				if err!=nil {
					fmt.Println("put error, %s",err)
				}

				_, err = hbaseClient.Put(biPutRequest)
				if err!=nil {
					fmt.Println("final error, %s",err)
				}
			}
		}
	}()


	//go func() {
	//	for i:=0;i<500 ;i++  {
	//		pool.JobQueue <- func() {
	//			basicInfoCfMapper := map[string]map[string][]byte{"basic":{"uid":[]byte("123")}}
	//			rowkey := fmt.Sprintf("%d%d",time.Now().Unix(),rand.Intn(10000))
	//			biPutRequest, err := hrpc.NewPutStr(context.Background(),"testTable",rowkey,basicInfoCfMapper)
	//			if err!=nil {
	//				fmt.Println("put error, %s",err)
	//			}
	//
	//			_, err = hbaseClient.Put(biPutRequest)
	//			if err!=nil {
	//				fmt.Println("final error, %s",err)
	//			}
	//		}
	//	}
	//}()

	select {

	}

}
package main

import (
	"github.com/zuston/AtcalMq/util"
	"time"
)

func main()  {
	logger := util.NewLog4Go()
	logger.SetType(1)
	logger.SetRollingDaily("/temp/","log4go")
	logger.Println("what the fuck world")

	logger.Println("what the fuck world")
	time.Sleep(time.Minute*10)
}

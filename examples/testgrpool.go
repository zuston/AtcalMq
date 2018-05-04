package main

import (
"fmt"
"time"

"github.com/ivpusic/grpool"
)

func main() {

	// number of workers, and size of job queue
	pool := grpool.NewPool(100, 50)

	// release resources used by pool
	defer pool.Release()

	// submit one or more jobs to pool
	for i := 0; i < 500; i++ {
		count := i

		pool.JobQueue <- func() {
			fmt.Printf("I am worker! Number %d\n", count)
			time.Sleep(2*time.Second)
		}
	}

	// dummy wait until jobs are finished
	time.Sleep(10 * time.Second)
}

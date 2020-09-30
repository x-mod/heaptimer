package main

import (
	"context"
	"log"
	"time"

	"github.com/x-mod/heaptimer"
)

func main() {
	log.Println("timer start ...")
	timer := heaptimer.New(
		heaptimer.Duration(time.Millisecond * 10),
	)
	go func() {
		if err := timer.Serve(context.TODO()); err != nil {
			log.Println("timer serving failed: ", err)
		}
	}()
	timer.Push(1, time.Now())
	timer.Push(2, time.Now())
	timer.Push(3, time.Now())
	timer.Push(4, time.Now())

	timer.PushWithDuration(5, time.Second*3)
	timer.PushWithDuration(6, time.Second*5)
	timer.PushWithDuration(7, time.Second*7)

	// <-timer.Close()
	go func() {
		time.Sleep(time.Millisecond * 5100)
		<-timer.Close()
	}()
	// method 1
	for i := range timer.C {
		log.Println("pop:", i)
	}
	// method 2
	// for {
	// 	if i, ok := timer.Pop(); ok {
	// 		log.Println("pop:", i)
	// 	} else {
	// 		break
	// 	}
	// }
	// drain
	for {
		if v := timer.Drain(); v != nil {
			log.Println("drain:", v)
		} else {
			break
		}
	}

	log.Println("timer closed")
}

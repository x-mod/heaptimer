package main

import (
	"context"
	"log"
	"time"

	"github.com/x-mod/heaptimer"
)

func main() {
	log.Println("timer start ...")
	timer := heaptimer.New()
	go func() {
		if err := timer.Serve(context.TODO()); err != nil {
			log.Println("timer serving failed: ", err)
		}
	}()

	go func() {
		timer.Push(1, time.Now())
		timer.Push(2, time.Now())
		timer.Push(3, time.Now())
		time.Sleep(time.Second)
		timer.Push(4, time.Now())

		timer.PushWithDuration(5, time.Second)
		timer.PushWithDuration(6, time.Second*3)
		timer.PushWithDuration(7, time.Second*5)
	}()

	go func() {
		<-timer.Serving()
		log.Println("timer is serving ....")
		time.Sleep(time.Second * 15)
		log.Println("timer is closing ....")
		<-timer.Close()
	}()
	// method 1
	for i := range timer.V {
		log.Println("pop:", i)
	}
	log.Println("timer closed")
}

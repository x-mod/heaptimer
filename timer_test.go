package heaptimer

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestTimer_Push(t *testing.T) {
	timer := New(time.Microsecond * 100)
	go timer.Serve(context.TODO())

	timer.Push(1, time.Second)
	timer.Push(2, 0)
	timer.Push(3, 0)
	timer.Push(4, 0)

	timer.Push(5, time.Second*3)
	timer.Push(6, time.Second*5)
	timer.Push(7, time.Second*7)

	go func() {
		time.Sleep(time.Second * 3)
		<-timer.Close()
	}()
	//method 1
	for i := range timer.C {
		log.Println("pop:", i)
	}
	//method 2
	// for {
	// 	if i, ok := timer.Pop(); ok {
	// 		log.Println("pop:", i)
	// 	} else {
	// 		break
	// 	}
	// }
	log.Println("timer closed")
}

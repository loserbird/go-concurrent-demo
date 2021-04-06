package goroutine

import (
	"fmt"
	"sync"
	"time"
)

func condDemo() {
	c := sync.NewCond(&sync.Mutex{})
	queue := make([]interface{}, 0, 10)
	removeFromQueue := func(delay time.Duration) {
		time.Sleep(delay)
		c.L.Lock()
		queue = queue[1:]
		fmt.Printf("removing from queue")
		c.L.Unlock()
		c.Signal()
	}

	for i := 0; i < 10; i++ {
		c.L.Lock()
		for len(queue) == 2 {
			c.Wait()
		}
		fmt.Printf("adding to queue")
		queue = append(queue, struct{}{})
		go removeFromQueue(1 * time.Second)
		c.L.Unlock()
	}
}

type Button struct {
	Clicked *sync.Cond
}

func condDemo2() {
	button := Button{
		Clicked: sync.NewCond(&sync.Mutex{}),
	}
	suscribe := func(c *sync.Cond, fn func()) {
		var goroutineRunning sync.WaitGroup
		goroutineRunning.Add(1)
		go func() {
			goroutineRunning.Done()
			c.L.Lock()
			defer c.L.Unlock()
			c.Wait()
			fn()
		}()
	}

	var clickRegistered sync.WaitGroup
	clickRegistered.Add(3)
	suscribe(button.Clicked, func() {
		fmt.Println("Maximizing window")
		clickRegistered.Done()
	})
	suscribe(button.Clicked, func() {
		fmt.Println("display dialog box")
		clickRegistered.Done()
	})
	suscribe(button.Clicked, func() {
		fmt.Println("mouse clicked")
		clickRegistered.Done()
	})
	button.Clicked.Broadcast()
	clickRegistered.Wait()
}

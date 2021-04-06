package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"sync"
	"time"
)

func main() {
	nums := []int{2, 5, 4, 9, 11, 3}
	sort.Slice(nums, func(i, j int) bool {
		return nums[i] < nums[j]
	})
	for _, v := range nums {
		fmt.Println(v)
	}
}

func chanDemo1() {
	stringStream := make(chan string)
	go func() {
		stringStream <- "hello"
	}()
	salutation, ok := <-stringStream
	fmt.Printf("(%v): %v", ok, salutation)
}

func chanDemo2() {
	intStream := make(chan int)
	go func() {
		defer close(intStream)
		for i := 1; i <= 5; i++ {
			intStream <- i
		}
	}()

	for i := range intStream {
		fmt.Printf("%v ", i)
	}
}

func chanDemo3() {
	var stdoutBuff bytes.Buffer
	defer stdoutBuff.WriteTo(os.Stdout)

	intStream := make(chan int, 4)
	go func() {
		defer close(intStream)
		defer fmt.Fprintln(&stdoutBuff, "producer done")
		for i := 0; i < 5; i++ {
			fmt.Fprintf(&stdoutBuff, "Sending: %d\n", i)
			intStream <- i
		}
	}()

	for i := range intStream {
		fmt.Fprintf(&stdoutBuff, "Receive %v.\n ", i)
	}
}

// select case 语句伪随机，随机选择一个就绪语句
func chanDemo4() {
	c1 := make(chan interface{})
	close(c1)
	c2 := make(chan interface{})
	close(c2)

	var c1Count, c2Count int
	for i := 1000; i >= 0; i-- {
		select {
		case <-c1:
			c1Count++
		case <-c2:
			c2Count++
		}
	}
	fmt.Printf("c1Count: %d\nc2Count: %d", c1Count, c2Count)
}

// select语句如果没有就绪的case,会阻塞，可以加个超时或者加个default case
func chanDemo5() {
	var c <-chan int
	select {
	case <-c:
	case <-time.After(1 * time.Second):
		fmt.Println("time out")
	}
}

func chanDemo6() {
	done := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()

	workCount := 0
loop:
	for {
		select {
		case <-done:
			break loop
		default:

		}
		workCount++
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("workcount: %d", workCount)
}

func chanDemo7() {
	printData := func(wg *sync.WaitGroup, data []byte) {
		defer wg.Done()

		var buff bytes.Buffer
		for _, b := range data {
			fmt.Fprintf(&buff, "%c", b)
		}
		fmt.Println(buff.String())
	}

	var wg sync.WaitGroup
	wg.Add(2)
	data := []byte("golang")
	go printData(&wg, data[:3])
	go printData(&wg, data[3:])
	wg.Wait()
}

func chanDemo8() {
	newRandomStream := func(done <-chan interface{}) <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("closure existed")
			defer close(randStream)
			for {
				select {
				case randStream <- rand.Int():
				case <-done:
					return
				}
			}
		}()
		return randStream
	}
	done := make(chan interface{})
	randStream := newRandomStream(done)
	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
	close(done)
	time.Sleep(1 * time.Second)
}

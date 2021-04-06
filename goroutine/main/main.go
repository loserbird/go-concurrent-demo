package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {

}

// 输出三个hi
func demo01() {
	var wg sync.WaitGroup
	for _, salutation := range []string{"hello", "world", "hi"} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(salutation)
		}()
	}
	wg.Wait()
}

func demo02() {
	var wg sync.WaitGroup
	for _, salutation := range []string{"hello", "world", "hi"} {
		wg.Add(1)
		go func(salutation string) {
			defer wg.Done()
			fmt.Println(salutation)
		}(salutation)
	}
	wg.Wait()
}

// 计算一个goroutine占用的内存大小， 几kb，很轻量级
func demo03() {
	memConsumed := func() uint64 {
		runtime.GC()
		var s runtime.MemStats
		runtime.ReadMemStats(&s)
		return s.Sys
	}
	var c <-chan interface{}
	var wg sync.WaitGroup
	noop := func() {
		wg.Done()
		<-c
	}

	const numGoroutines = 1e4
	wg.Add(numGoroutines)
	before := memConsumed()
	for i := numGoroutines; i > 0; i-- {
		go noop()
	}
	wg.Wait()
	after := memConsumed()
	fmt.Printf("%.3fkb", float64(after-before)/numGoroutines/1000)
}

func demo04() {
	var count int
	increment := func() {
		count++
	}
	var once sync.Once
	var increments sync.WaitGroup
	increments.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer increments.Done()
			once.Do(increment)
		}()
	}
	increments.Wait()
	fmt.Printf("Count is %d\n", count)
}

func demo05() {
	var count int
	increment := func() {
		count++
	}
	decrement := func() {
		count--
	}
	var once sync.Once
	once.Do(increment)
	once.Do(decrement)
	fmt.Printf("Count is %d\n", count) // 1
}

package main

import "fmt"

func main() {
	orDone := func(done, c <-chan interface{}) <-chan interface{} {
		valStream := make(chan interface{})
		go func() {
			defer close(valStream)
			for {
				select {
				case <-done:
					return
				case val, ok := <-c:
					if ok == false {
						return
					}
					select {
					case valStream <- val:
					case <-done:
					}
				}
			}
		}()
		return valStream
	}
	myChan := make(chan interface{})
	done := make(chan interface{})
	defer close(done)
	defer close(myChan)
	for val := range orDone(done, myChan) {
		fmt.Println(val)
	}
}

package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	//goroutines()

	channels()

	fmt.Println("Hello from main goroutine")

	time.Sleep(time.Second)

	fmt.Println("main returns")
}

func goroutines() {
	// A goroutine is an execution context that is managed by the Go runtime
	// (as opposed to a thread that is managed by the operating system)

	// The go keyword starts the given function in a new goroutine.
	f := func() {
		fmt.Println("Hello from goroutine with func f")
	}
	go f()

	// A goroutine usually has a much smaller startup overhead than an operating system thread.

	// The function running as a goroutine can take parameters, but it cannot return a value.
	g := func(i, j int) {
		fmt.Printf("Hello from goroutine with func g: %d %d\n", i, j)
	}
	i, j := 1, 2
	go g(i, j)

	// A goroutine starts with a small stack that grows as needed.
	go func() {
		fmt.Println("Hello from goroutine with anonymous func")
	}()

	// Creating new goroutines is faster and cheaper than creating operation system threads.

	// The parameters of the goroutine function are evaluated before the goroutine starts
	// and passed to the function once the goroutine starts running.
	go func(i, j int) {
		fmt.Printf("Hello from goroutine with another anonymous func: %d %d\n", i, j)
	}(1, 2)

	// The Go scheduler assigns operating system threads to run goroutines.

	// The number of operating system threads used by the Go runtime is equal to the number
	// of processors/cores on the platform (unless you change this by setting the GOMAXPROCS environment variable
	// or by calling the runtime.GOMAXPROCS function).
	go1 := runtime.NumGoroutine()
	fmt.Printf("The number of goroutines that currently exist: %d\n", go1)

	proc1 := runtime.GOMAXPROCS(4)
	fmt.Printf("The number of processors/cores after updating: %d\n", proc1)

	go2 := runtime.NumGoroutine()
	fmt.Printf("The number of goroutines that currently exist: %d\n", go2)

	// Every run of this code is likely to print out a, b, and c in random order.
	h := func(s string) {
		fmt.Printf("Goroutine %s\n", s)
	}
	for _, s1 := range []string{"a", "b", "c"} {
		go h(s1)
	}

	// Before Go v1.22: data race
	for _, s2 := range []string{"x", "y", "z"} {
		s2 := s2 // fixed for Go v1.21 and earlier

		go func() {
			fmt.Printf("Goroutine %s\n", s2)
		}()
	}

	// After Go v1.22: data race
	var s3 string // here s3 in func goroutines scope
	for _, s3 = range []string{"k", "l", "m"} {
		s3 := s3 // fixed as for Go v1.21 and earlier, here s3 in for loop scope

		go func() {
			fmt.Printf("Goroutine %s\n", s3)
		}()
	}

	// Data race solution
	var wg sync.WaitGroup

	for _, s4 := range []string{"11", "22", "33"} {
		wg.Add(1)

		go func(s string) {
			defer wg.Done()
			fmt.Println(s)
		}(s4)
	}

	wg.Wait()
}

func channels() {
	// Channels allow goroutines to share memory by communicating, as opposed to communicating by sharing memory.
	// You can declare a channel by specifying its type and its capacity.
	ch := make(chan int, 2)

	// A channel is a first-in, first-out (FIFO) conduit.
	// Use the following syntax to send to or receive from channels.
	ch <- 1
	fmt.Printf("len: %d\n", len(ch))
	<-ch

	ch <- 2
	fmt.Printf("len: %d\n", len(ch))
	x := <-ch
	fmt.Printf("%d\n", x)

	ch <- 3
	fmt.Printf("len: %d\n", len(ch))
	// Don't do this in real code with many goroutines writers/readers!
	if len(ch) > 0 {
		y := <-ch
		fmt.Printf("%d\n", y)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		fmt.Println("Writer starts")
		ch <- 4
		ch <- 5
		// a send operation to a channel will block until the channel is ready to accept a value
		ch <- 6
		fmt.Println("Writer returns")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		fmt.Println("Reader starts")
		fmt.Printf("len: %d\n", len(ch))
		a := <-ch
		fmt.Printf("%d\n", a)
		b := <-ch
		fmt.Printf("%d\n", b)
		c := <-ch
		fmt.Printf("%d\n", c)
		// a receive operation from a channel will block until the channel is ready to provide a value
		//d := <-ch
		//fmt.Printf("%d\n", d)

		fmt.Println("Reader returns")
	}()

	wg.Wait()

	// reading from or writing to a nil channel will block indefinitely
	//var nilChan chan int
	//nilChan <- 1
	//<-nilChan

	// writing to a closed channel will always panic
	//doneChan := make(chan int, 2)
	//close(doneChan)
	//doneChan <- 1

	// for a receiver, it is usually important to know whether the channel was closed when the read happened
	doneChan := make(chan int, 2)
	close(doneChan)
	_, ok := <-doneChan
	if !ok {
		fmt.Println("channel was closed")
	}

	// unbuffered channel acts as a synchronization point between two goroutines
	// there are two possible runs at this point:
	// sender attempts to send before receiver is ready to receive
	// receiver attempts to receive before sender is ready to send
	chUnbuffered := make(chan bool)

	// this is equivalent to yy=xx, with additional synchronization guarantees
	// it does not transfer the ownership of the value

	// sender
	go func() {
		xx := true
		chUnbuffered <- xx
	}()

	// receiver
	go func() {
		yy := <-chUnbuffered
		fmt.Println(yy)
	}()

	// TODO: Directional channels
	// TODO: Worker pool sample
}

// TODO: select samples

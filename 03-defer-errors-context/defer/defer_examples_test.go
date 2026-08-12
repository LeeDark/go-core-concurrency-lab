package deferlab

import "fmt"

func Example_lifo() {
	defer fmt.Println("release A")
	defer fmt.Println("release B")
	defer fmt.Println("release C")

	// Output:
	// release C
	// release B
	// release A
}

func Example_argumentEvaluation() {
	x := 10
	defer fmt.Println("argument:", x)

	x = 20

	// Output:
	// argument: 10
}

func Example_closure() {
	x := 10
	defer func() {
		fmt.Println("closure:", x)
	}()

	x = 20

	// Output:
	// closure: 20
}

func Example_namedResult() {
	result := calculate()
	fmt.Println(result)

	// Output:
	// 11
}

func calculate() (result int) {
	defer func() {
		result++
	}()

	return 10
}

func Example_loopDefer() {
	for i := 1; i <= 3; i++ {
		fmt.Println("open", i)
		defer fmt.Println("close", i)
	}

	fmt.Println("work")

	// Output:
	// open 1
	// open 2
	// open 3
	// work
	// close 3
	// close 2
	// close 1
}

package main

import (
	"strconv"
	"syscall/js"
)

func main() {
	println("GO:: Hello world!")
}

//export add
func add(a, b int) int {
	return a + b
}

//export update
func update() {
	// var val js.Value
	document := js.Global().Get("document")

	aStr := document.Call("getElementById", "a").Get("value").String()
	bStr := document.Call("getElementById", "b").Get("value").String()
	a, _ := strconv.Atoi(aStr)
	b, _ := strconv.Atoi(bStr)
	result := add(a, b)
	document.Call("getElementById", "result").Set("value", result)
}

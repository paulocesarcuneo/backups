package main

import (
	"fmt"
)

type A struct {
	A string
	I int
	M map[string]string
}

func main() {
	a := A{}
	var b A
	fmt.Println(a, b)
}

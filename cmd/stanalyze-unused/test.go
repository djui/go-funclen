package main

import (
	"fmt"

	"github.com/djui/go-stanalyzer/foo/bar"
)

const C = 0

const (
	C1 = 1
	C2 = 2
)

var V = 3

func foo() {
	fmt.Println(C1)
	fmt.Println(baz.Z)
}

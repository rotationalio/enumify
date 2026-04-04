package main

import (
	"fmt"
	"os"
)

func main() {
	fname := os.Getenv("GOFILE")
	pkg := os.Getenv("GOPACKAGE")
	fmt.Println(fname, pkg)
}

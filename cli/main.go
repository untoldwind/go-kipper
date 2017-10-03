package main

import (
	"fmt"
	"os"

	"github.com/untoldwind/go-kipper"
)

func handleError(err error) {
	fmt.Fprintln(os.Stdout, err.Error())
	os.Exit(1)
}

func main() {
	clipboard, err := kipper.NewClipboard("go-kipper")
	if err != nil {
		handleError(err)
	}

	fmt.Println(clipboard.Get(kipper.AtomPrimary))

	fmt.Println(clipboard.Get(kipper.AtomClipboard))
}

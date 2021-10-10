package main

import (
	"fmt"
	"os"

	invasioncmd "github.com/ivanovpetr/invasion/cmd"
)

func main() {
	err := invasioncmd.New().Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

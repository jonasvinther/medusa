package main

import (
	"fmt"
	"github.com/jonasvinther/medusa/cmd"
	"os"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

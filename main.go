package main

import (
	"log"
	"os"

	"github.com/bshore/steggo/cmd"
)

func main() {
	cmd.InitRoot()
	if err := cmd.Execute(os.Args[1:]); err != nil {
		log.Printf("exiting with error: %v\n", err)
		os.Exit(1)
	}
}

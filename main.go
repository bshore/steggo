package main

import (
	"log"
	"lsb_encoder/cmd"
	"os"
)

func main() {
	cmd.InitRoot()
	if err := cmd.Execute(os.Args[1:]); err != nil {
		log.Printf("exiting with error: %v\n", err)
		os.Exit(1)
	}
}

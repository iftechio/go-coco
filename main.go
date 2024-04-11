package main

import (
	"os"

	"github.com/iftechio/go-coco/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

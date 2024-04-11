package main

import (
	"os"

	"github.com/xieziyu/go-coco/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

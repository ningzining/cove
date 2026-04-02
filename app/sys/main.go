package main

import (
	"os"

	"github.com/ningzining/cove/app/sys/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

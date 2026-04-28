package main

import (
	"os"

	"github.com/ningzining/cove/internal/system"
)

func main() {
	if err := system.Execute(); err != nil {
		os.Exit(1)
	}
}

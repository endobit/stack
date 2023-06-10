// Package main is the stackd daemon.
package main

import (
	"os"

	"github.com/endobit/stack/internal/service"
)

var version string

func main() {
	root := service.NewRootCmd(version)
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

package main

import (
	"os"

	"github.com/Blaze757/routeflux/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}

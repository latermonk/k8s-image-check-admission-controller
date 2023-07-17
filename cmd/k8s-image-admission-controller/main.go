package main

import (
	"os"

	"k8s-image-check-admission-controller/pkg/config"
)

func main() {
	cmd := config.NewRootCommand()

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

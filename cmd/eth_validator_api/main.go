package main

import (
	"ethereum-validator-api/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}

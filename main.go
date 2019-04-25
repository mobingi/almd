package main

import (
	"math/rand"
	"time"
	"fmt"
	"os"

	"github.com/mobingilabs/ouchan/services/oceand/server"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	command := server.NewOceandCommand()
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

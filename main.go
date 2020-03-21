package main

import (
	"Bid/server"
	"Bid/utils/memory"
	"log"
)

func main() {
	mem := memory.NewMem()
	err := server.Run(mem)
	if err != nil {
		log.Fatal("cannot start the server", err)
	}
}

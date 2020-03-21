package main

import (
	"Bid/server"
	"Bid/utils/memory"
	"log"
)

func main() {
	txn := memory.NewTxn()
	defer txn.Abort()
	err := server.Run(txn)
	if err != nil {
		log.Fatal("cannot start the server", err)
	}
}

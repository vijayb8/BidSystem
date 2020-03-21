package server

import (
	"Bid/router"
	"Bid/utils/memory"
)

func Run(txn memory.TxnIn) error {
	r := router.Init(txn)
	return r.Run(":8080")
}

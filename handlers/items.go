package handlers

import (
	"Bid/structs"
	"Bid/utils/memory"
	"Bid/utils/responses"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/hashicorp/go-memdb"
	"net/http"
)

var (
	itemTable = "items"
)

func GetItems(txn memory.TxnIn) gin.HandlerFunc {
	return func(c *gin.Context) {
		it, err := txn.Get(itemTable, "item_id")
		if err != nil {
			responses.ResponseWithError(c, http.StatusInternalServerError, fmt.Errorf("cannot get items"))
			return
		}
		responses.ResponseWithData(c, http.StatusOK, getItems(it))
		return
	}
}

func GetItemsByBidId(txn memory.TxnIn) gin.HandlerFunc {
	return func(c *gin.Context) {
		bidId := c.Param("bidId")
		it, err := txn.Get(itemTable, "bid_id", bidId)
		if err != nil {
			responses.ResponseWithError(c, http.StatusInternalServerError, fmt.Errorf("cannot get items for the bid id"))
			return
		}
		responses.ResponseWithData(c, http.StatusOK, getItems(it))
		return
	}
}

func CreateItem(txn memory.TxnIn) gin.HandlerFunc {
	return func(c *gin.Context) {
		item := structs.Item{}
		err := c.MustBindWith(&item, binding.JSON)
		if err != nil {
			responses.ResponseWithError(c, http.StatusBadRequest, fmt.Errorf("bad request"))
			return
		}
		err = txn.Write(itemTable, item)
		if err != nil {
			responses.ResponseWithError(c, http.StatusInternalServerError, fmt.Errorf("cannot create item"))
			return
		}
		responses.ResponseWithData(c, http.StatusOK, "Item created successfully")
		return
	}
}

func getItems(it memdb.ResultIterator) []structs.Item {
	var items []structs.Item
	for obj := it.Next(); obj != nil; obj = it.Next() {
		item := obj.(*structs.Item)
		items = append(items, *item)
	}
	return items
}

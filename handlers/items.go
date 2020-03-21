package handlers

import (
	"Bid/structs"
	"Bid/utils/memory"
	"Bid/utils/responses"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/hashicorp/go-memdb"
	"net/http"
)

var (
	itemTable = "items"
)

func GetItems(txn memory.TxnIn) gin.HandlerFunc {
	return func(c *gin.Context) {
		it, err := txn.Get(itemTable, "id", nil)
		if err != nil {
			responses.ResponseWithError(c, http.StatusInternalServerError, fmt.Errorf("cannot get items"))
			return
		}
		responses.ResponseWithData(c, http.StatusOK, getItems(it))
		return
	}
}

func GetItemById(txn memory.TxnIn) gin.HandlerFunc {
	return func(c *gin.Context) {
		itemId := c.Param("itemId")
		it, err := txn.First(bidTable, "id", &itemId)
		if err != nil {
			responses.ResponseWithError(c, http.StatusInternalServerError, fmt.Errorf("cannot get item"))
			return
		}
		responses.ResponseWithData(c, http.StatusOK, it.(structs.Item))
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
		id, _ := uuid.NewRandom()
		item.ItemId = id.String()
		err = txn.Write(itemTable, item)
		if err != nil {
			responses.ResponseWithError(c, http.StatusInternalServerError, fmt.Errorf("cannot create item"))
			return
		}
		responses.ResponseWithData(c, http.StatusOK, item.ItemId)
		return
	}
}

func getItems(it memdb.ResultIterator) []structs.Item {
	var items []structs.Item
	for obj := it.Next(); obj != nil; obj = it.Next() {
		item := obj.(structs.Item)
		items = append(items, item)
	}
	return items
}

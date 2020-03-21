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
	bidTable = "bids"
)

func GetBidsByUserId(txn memory.TxnIn) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")
		it, err := txn.Get(bidTable, "userId", &userId)
		if err != nil || it == nil {
			responses.ResponseWithError(c, http.StatusInternalServerError, fmt.Errorf("No data or cannot get the data for given user id"))
			return
		}
		responses.ResponseWithData(c, http.StatusOK, getBids(it))
		return
	}
}

func GetBids(txn memory.TxnIn) gin.HandlerFunc {
	return func(c *gin.Context) {
		it, err := txn.Get(bidTable, "id", nil)
		if err != nil {
			responses.ResponseWithError(c, http.StatusInternalServerError, fmt.Errorf("cannot get bids"))
			return
		}
		responses.ResponseWithData(c, http.StatusOK, getBids(it))
		return
	}
}

func GetBidById(txn memory.TxnIn) gin.HandlerFunc {
	return func(c *gin.Context) {
		bidId := c.Param("bidId")
		it, err := txn.First(bidTable, "id", &bidId)
		if err != nil {
			responses.ResponseWithError(c, http.StatusInternalServerError, fmt.Errorf("cannot get bids"))
			return
		}
		responses.ResponseWithData(c, http.StatusOK, it.(structs.Bid))
		return
	}
}

func GetBidsByItemId(txn memory.TxnIn) gin.HandlerFunc {
	return func(c *gin.Context) {
		itemId := c.Param("itemId")
		it, err := txn.Get(bidTable, "itemId", &itemId)
		if err != nil {
			responses.ResponseWithError(c, http.StatusInternalServerError, fmt.Errorf("cannot get items for the bid id"))
			return
		}
		responses.ResponseWithData(c, http.StatusOK, getBids(it))
		return
	}
}

func GetMaxBidByItemId(txn memory.TxnIn) gin.HandlerFunc {
	return func(c *gin.Context) {
		itemId := c.Param("itemId")
		it, err := txn.Get(bidTable, "itemId", &itemId)
		if err != nil {
			responses.ResponseWithError(c, http.StatusInternalServerError, fmt.Errorf("cannot get data by itemId"))
			return
		}
		winningBid := structs.Bid{}
		for _, bid := range getBids(it) {
			if bid.BidAmount > winningBid.BidAmount {
				winningBid = bid
			}
		}
		responses.ResponseWithData(c, http.StatusOK, winningBid)
		return
	}
}

func CreateBid(txn memory.TxnIn) gin.HandlerFunc {
	return func(c *gin.Context) {
		bid := structs.Bid{}
		err := c.MustBindWith(&bid, binding.JSON)
		if err != nil {
			responses.ResponseWithError(c, http.StatusBadRequest, fmt.Errorf("bad request"))
			return
		}
		//Check if the user is present to bid
		it, err := txn.First("users", "id", &bid.UserId)
		if err != nil || it == nil {
			responses.ResponseWithError(c, http.StatusInternalServerError, fmt.Errorf("cannot create bid"))
			return
		}

		// Check if item is available to bid
		it, err = txn.First("items", "id", &bid.ItemId)
		if err != nil || it == nil {
			responses.ResponseWithError(c, http.StatusInternalServerError, fmt.Errorf("cannot create bid"))
			return
		}
		item := it.(structs.Item)
		if !item.IsAvailableToBid {
			responses.ResponseWithError(c, http.StatusInternalServerError, fmt.Errorf("Item is not available"))
			return
		}

		id, _ := uuid.NewRandom()
		bid.BidId = id.String()
		err = txn.Write(bidTable, bid)
		if err != nil {
			responses.ResponseWithError(c, http.StatusInternalServerError, fmt.Errorf("cannot create bid"))
			return
		}
		responses.ResponseWithData(c, http.StatusOK, bid.BidId)
		return
	}
}

func getBids(it memdb.ResultIterator) []structs.Bid {
	var bids []structs.Bid
	for obj := it.Next(); obj != nil; obj = it.Next() {
		bid := obj.(structs.Bid)
		bids = append(bids, bid)
	}
	return bids
}

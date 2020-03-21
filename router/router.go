package router

import (
	handlers "Bid/handlers"
	"Bid/utils/memory"
	"github.com/gin-gonic/gin"
)

func Init(txn memory.TxnIn) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	userGroup := r.Group("users")
	userGroup.GET("/", handlers.GetUsers(txn))
	userGroup.POST("/user", handlers.CreateUser(txn))
	userGroup.GET("/:userId/bids", handlers.GetBidsByUserId(txn))

	itemGroup := r.Group("items")
	itemGroup.GET("/", handlers.GetItems(txn))
	itemGroup.POST("/item", handlers.CreateItem(txn))
	itemGroup.GET("/:itemId/bid/max", handlers.GetMaxBidByItemId(txn))
	itemGroup.GET("/:itemId/bids", handlers.GetItemsByBidId(txn))

	bidGroup := r.Group("bids")
	bidGroup.GET("/", handlers.GetBids(txn))
	bidGroup.GET("/:bidID", handlers.GetBidById(txn))
	bidGroup.POST("/bid", handlers.CreateBid(txn))

	return r
}

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
	userTable = "users"
)

func GetUsers(txn memory.TxnIn) gin.HandlerFunc {
	return func(c *gin.Context) {
		it, err := txn.Get(userTable, "id", nil)
		if err != nil {
			responses.ResponseWithError(c, http.StatusInternalServerError, fmt.Errorf("cannot get items"))
			return
		}
		responses.ResponseWithData(c, http.StatusOK, getUsers(it))
		return
	}
}

func CreateUser(txn memory.TxnIn) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := structs.User{}
		err := c.MustBindWith(&user, binding.JSON)
		if err != nil {
			responses.ResponseWithError(c, http.StatusBadRequest, fmt.Errorf("bad request"))
			return
		}
		id, _ := uuid.NewRandom()
		user.UserId = id.String()
		err = txn.Write(userTable, user)
		if err != nil {
			responses.ResponseWithError(c, http.StatusInternalServerError, fmt.Errorf("cannot create user"))
			return
		}
		responses.ResponseWithData(c, http.StatusOK, user.UserId)
		return
	}
}

func getUsers(it memdb.ResultIterator) []structs.User {
	var users []structs.User
	for obj := it.Next(); obj != nil; obj = it.Next() {
		user := obj.(structs.User)
		users = append(users, user)
	}
	return users
}
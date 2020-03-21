package main

import (
	"Bid/router"
	"Bid/structs"
	"Bid/utils/memory"
	"bytes"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"testing"
	"time"
)

var (
	mem        memory.TxnIn
	httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
)

func BenchmarkMain(m *testing.B) {
	mem = memory.NewMem()

	testServerListener, _ := net.Listen("tcp", ":8181")
	testServerSv := &http.Server{Handler: router.Init(mem)}
	go func() { log.Fatal(testServerSv.Serve(testServerListener)) }()
	log.Println(m.Run("CreateBid", BenchmarkCreateBid))
}

func BenchmarkCreateBid(t *testing.B) {
	item := structs.Item{
		Name:             "item",
		Price:            10,
		IsAvailableToBid: true,
	}
	data, _ := jsoniter.MarshalToString(item)
	jsonStr := []byte(data)
	resp, err := httpClient.Post("http://127.0.0.1:8181/items/item", "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	da, _ := ioutil.ReadAll(resp.Body)
	itemId := jsoniter.Get(da, "data").ToString()

	for i := 0; i < 100; i++ {
		user := structs.User{
			Name:  "vijay" + strconv.Itoa(i),
			Email: "test@email.com",
		}
		data, _ := jsoniter.MarshalToString(user)
		var jsonStr = []byte(data)
		resp, err := httpClient.Post("http://127.0.0.1:8181/users/user", "application/json", bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatal(err)
		}
		da, _ := ioutil.ReadAll(resp.Body)
		userId := jsoniter.Get(da, "data").ToString()

		bid := structs.Bid{
			UserId:    userId,
			ItemId:    itemId,
			BidAmount: float32(20 + i),
		}
		data, _ = jsoniter.MarshalToString(bid)
		jsonStr = []byte(data)
		resp, err = httpClient.Post("http://127.0.0.1:8181/bids/bid", "application/json", bytes.NewBuffer(jsonStr))
		if err != nil {
			t.Fatal(err)
		}
		da, _ = ioutil.ReadAll(resp.Body)
		bidId := jsoniter.Get(da, "data").ToString()
		if bidId == "" {
			log.Fatal("cannot get bidId")
		}

		//get items by UserId
		resp, err = httpClient.Get("http://127.0.0.1:8181/users/" + userId + "/bids")
		if err != nil {
			log.Fatal("cannot get bids")
		}
		da, _ = ioutil.ReadAll(resp.Body)
		sz := jsoniter.Get(da, "data").Size()
		if sz == 0 {
			log.Fatal("No data to read")
		}
	}
	//get max bid by itemId
	//get bids by userId
	resp, err = httpClient.Get("http://127.0.0.1:8181/items/" + itemId + "/bids")
	if err != nil {
		log.Fatal("cannot get bids")
	}
	da, _ = ioutil.ReadAll(resp.Body)
	sz := jsoniter.Get(da, "data").Size()
	if sz == 0 {
		log.Fatal("No data to read")
	}
}

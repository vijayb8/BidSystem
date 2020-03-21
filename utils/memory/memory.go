package memory

import (
	"github.com/hashicorp/go-memdb"
	"log"
	"sync"
)

type TxnIn interface {
	Write(table string, data interface{}) error
	Get(table string, index string, args *string) (memdb.ResultIterator, error)
	First(table string, index string, args *string) (interface{}, error)
}

type mem struct {
	db *memdb.MemDB
	mu sync.Mutex
}

func NewMem() *mem {
	return &mem{
		db: InitDB(),
		mu: sync.Mutex{},
	}
}

func InitDB() *memdb.MemDB {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"users": &memdb.TableSchema{
				Name: "users",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "UserId"},
					},
				},
			},
			"items": &memdb.TableSchema{
				Name: "items",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ItemId"},
					},
				},
			},
			"bids": &memdb.TableSchema{
				Name: "bids",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "BidId"},
					},
					"userId": &memdb.IndexSchema{
						Name:    "userId",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "UserId"},
					},
					"itemId": &memdb.IndexSchema{
						Name:    "itemId",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "ItemId"},
					},
				},
			},
		},
	}

	db, err := memdb.NewMemDB(schema)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func (t *mem) Write(table string, data interface{}) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	txn := t.db.Txn(true)
	defer txn.Abort()
	err := txn.Insert(table, data)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (t *mem) Get(table string, index string, args *string) (memdb.ResultIterator, error) {
	txn := t.db.Txn(false)
	defer txn.Abort()
	if args == nil {
		return txn.Get(table, index)
	}
	return txn.Get(table, index, *args)
}

func (t *mem) First(table string, index string, args *string) (interface{}, error) {
	txn := t.db.Txn(false)
	defer txn.Abort()
	if args == nil {
		return txn.First(table, index)
	}
	return txn.First(table, index, *args)
}

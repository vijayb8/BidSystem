package memory

import (
	"github.com/hashicorp/go-memdb"
	"log"
	"sync"
)

type TxnIn interface {
	InitDB() *memdb.Txn
	Write(table string, data interface{}) error
	Get(table string, index string, args ...interface{}) (memdb.ResultIterator, error)
	First(table string, index string, args ...interface{}) (interface{}, error)
	Abort()
}

type txn struct {
	tx *memdb.Txn
	mu sync.Mutex
}

func NewTxn() *txn {
	return &txn{
		tx: nil,
		mu: sync.Mutex{},
	}
}

func (t *txn) InitDB() *memdb.Txn {
	t.mu.Lock()
	defer t.mu.Unlock()
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"user": &memdb.TableSchema{
				Name: "user",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:   "id",
						Unique: true,
					},
				},
			},
		},
	}

	db, err := memdb.NewMemDB(schema)
	if err != nil {
		log.Fatal(err)
	}

	return db.Txn(true)
}

func (t *txn) Write(table string, data interface{}) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	err := t.tx.Insert(table, data)
	if err != nil {
		return err
	}
	t.tx.Commit()
	return nil
}

func (t *txn) Get(table string, index string, args ...interface{}) (memdb.ResultIterator, error) {
	return t.tx.Get(table, index, args)
}

func (t *txn) First(table string, index string, args ...interface{}) (interface{}, error) {
	return t.tx.First(table, index, args)
}

func (t *txn) Abort() {
	t.tx.Abort()
}
// Memdb provides transaction and MVCC - Multi version concurrency control
// there can be n readers without lock and only 1 writer which can make propgress , and they will be synchronous in nature.

package memdb

import (
	"github.com/hashicorp/go-immutable-radix"
	"sync"
	"sync/atomic"
	"unsafe"
)

// MemDB is a in memory db.
//
type MemDB struct {
	schema *DBSchema
	root unsafe.Pointer
	primary bool

	writer sync.Mutex
}

func NewMemDB(schema *DBSchema) (*MemDB, error) {
	// Validate the schema
	if err := schema.Validate(); err != nil {
		return nil, err
	}

	// Create the MemDB
	db := &MemDB {
		schema: schema,
		root: unsafe.Pointer(iradix.New()),
		primary: true,
	}
	if err := db.initialise(); err != nil {
		return nil, err
	}
	return db, nil
}
// Initalise func is used to setup the DB for use after creation.
func (db *MemDB) initialise() error {
	root := db.getRoot()
	for tName, tableSchema := range db.schema.Tables {
		for iName := range tableSchema.Indexes {
			index := iradix.New()
			path := indexPath(tName, iName)
			root, _, _ = root.Insert(path, index)
		}
	}
	db.root = unsafe.Pointer(root)
	return nil
}

func indexPath(table, index string) []byte {
	return []byte(table + "." + index)
}

func (db *MemDB) getRoot() *iradix.Tree {
	root := (*iradix.Tree)(atomic.LoadPointer(&db.root))
	return root
}

//func (db *MemDB) Txn(write bool) *Txn {
//	if write {
//		db.writer.Lock()
//	}
//	txn := &Txn{
//		db: db,
//		write: write,
//		rootTxn: db.getRoot().Txn(),
//	}
//	return txn
//}
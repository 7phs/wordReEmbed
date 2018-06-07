package storages

import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDBTable struct {
	dbName string

	db *leveldb.DB
}

func NewLevelDB(dbName string) *LevelDBTable {
	return &LevelDBTable{
		dbName: dbName,
	}
}

func (o *LevelDBTable) Init() (*LevelDBTable, error) {
	var (
		err error
	)

	o.db, err = leveldb.OpenFile(o.dbName, nil)
	if err != nil {
		return nil, errors.New("failed to open db file '" + o.dbName + "': " + err.Error())
	}

	return o, nil
}

func (o *LevelDBTable) WeightInc(wordIndex, wordIndex2 int, weight float64) {
	key := o.key(wordIndex, wordIndex2)

	data, err := o.db.Get(key, nil)
	if err == nil {
		weight += o.toValue(data)
	} else if err != leveldb.ErrNotFound {
		fmt.Println("LevelDB: WeightInc: Get:", err)
	}

	if err := o.db.Put(key, o.value(weight), nil); err != nil {
		fmt.Println("LevelDB: WeightInc: Put:", err)
	}
}

func (o *LevelDBTable) key(wordIndex, wordIndex2 int) []byte {
	return (*[8]byte)(unsafe.Pointer(&[]int{wordIndex, wordIndex2}))[:]
}

func (o *LevelDBTable) value(weight float64) []byte {
	return (*[8]byte)(unsafe.Pointer(&weight))[:]
}

func (o *LevelDBTable) toValue(data []byte) float64 {
	return *(*float64)(unsafe.Pointer(&data))
}

func (o *LevelDBTable) Flush() {
	if o.db != nil {
		o.db.Close()
	}
}

package storecache

import (
	"github.com/boltdb/bolt"
	"fmt"
	"errors"
	"encoding/binary"
)

var db *bolt.DB = nil

func Open(path string) {
	var err error
	db, err = bolt.Open(path, 0600, nil)

	if err != nil {
		panic(fmt.Errorf("Error opening bolt db: %s", path))
	}
}

func Close() {
	if db != nil {
		db.Close()
	}
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}


func RetrieveCache(bzID int, currentDateTime string) (xmlContent *[]byte, errRes error) {

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(itob(bzID))
		if b == nil {
			errRes = errors.New("Not found")
			return errRes;
		}
		v := b.Get([]byte("lastDateTime"))
		if currentDateTime != "" && string(v) != currentDateTime {
			errRes = errors.New("The xml is outdated")
		}
		v = b.Get([]byte("xmlContent"))
		xmlContent = &v
		return nil
	})
	return xmlContent, errRes
}

func StoreCache(bzID int, lastDateTime string, xmlContent *[]byte) {
	var err error

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(itob(bzID))
		if b == nil {
			b, err = tx.CreateBucket(itob(bzID))
		}
		err := b.Put([]byte("lastDateTime"), []byte(lastDateTime))
		if err != nil {
			return err
		}
		err = b.Put([]byte("xmlContent"), *xmlContent)
		return err
	})
}


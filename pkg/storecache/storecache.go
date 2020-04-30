package storecache

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

var db *bolt.DB = nil

const credentialsBucket = "credentials"

func Open(path string) {
	var err error
	db, err = bolt.Open(path, 0600, &bolt.Options{Timeout: 5 * time.Second})

	if err != nil {
		fmt.Printf("Error opening bolt db: %s\n", path)
		fmt.Printf("Reason: %s", err)
		if err.Error() == "timeout" {
			fmt.Println("\n")
			fmt.Println("Possibly another instance of track has the database locked")
			fmt.Println("please close any other instance, or delete the database")
			fmt.Println("\n")
		}
		os.Exit(2)
	}
}

func Close() {
	if db != nil {
		db.Close()
	}
}

// itob returns an 8-byte big endian representation of v.
func itob(v int, isXML bool) []byte {
	if !isXML {
		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, uint64(v))
		return b
	} else {
		b := make([]byte, 9)
		binary.BigEndian.PutUint64(b, uint64(v))
		b[8] = 1
		return b
	}
}

func RetrieveCache(bzID int, currentDateTime string, isXml bool) (xmlContent *[]byte, errRes error) {

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(itob(bzID, isXml))
		if b == nil {
			errRes = errors.New("Not found")
			return errRes
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

func StoreCache(bzID int, lastDateTime string, xmlContent *[]byte, isXml bool) {
	var err error
	if xmlContent == nil {
		return
	}

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(itob(bzID, isXml))
		if b == nil {
			b, err = tx.CreateBucket(itob(bzID, isXml))
		}

		err := b.Put([]byte("lastDateTime"), []byte(lastDateTime))
		if err != nil {
			return err
		}
		err = b.Put([]byte("xmlContent"), *xmlContent)
		return err
	})
}

func StoreBzAuth(cookies []*http.Cookie, authToken string) {
	bucketName := []byte("credentials")

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			var err error
			b, err = tx.CreateBucket(bucketName)
			if err != nil {
				return err
			}
		}
		b.Put([]byte("token"), []byte(authToken))

		for _, cookie := range cookies {
			b.Put([]byte(cookie.Name), []byte(cookie.Value))
		}
		return nil
	})
}

func GetBzAuth() (cookies []*http.Cookie, authToken *string) {
	bucketName := []byte("credentials")
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			errRes := errors.New("Not found")
			return errRes
		}

		for _, n := range []string{"Bugzilla_login", "Bugzilla_logincookie"} {
			v := b.Get([]byte(n))
			if v != nil {
				cookies = append(cookies, &http.Cookie{
					Name:  n,
					Value: string(v),
				})

			}
		}
		v := b.Get([]byte("token"))

		if v != nil {
			tokenStr := string(v)
			authToken = &tokenStr
		}

		return nil
	})

	return cookies, authToken
}

func storeString(bucket string, key string, value string) {
	bucketName := []byte(bucket)

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			var err error
			b, err = tx.CreateBucket(bucketName)
			if err != nil {
				return err
			}
		}
		b.Put([]byte(key), []byte(value))

		return nil
	})
}

func getString(bucket string, key string) (value *string) {
	bucketName := []byte(bucket)
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			errRes := errors.New("Not found")
			return errRes
		}

		v := b.Get([]byte(key))
		if v != nil {
			tokenStr := string(v)
			value = &tokenStr
		}
		return nil
	})

	return value
}

func StoreTrelloToken(authToken string) {
	storeString(credentialsBucket, "trelloToken", authToken)
}

func GetTrelloToken() (authToken *string) {
	return getString(credentialsBucket, "trelloToken")
}

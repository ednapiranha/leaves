package db

import (
	"encoding/json"
	"log"

	"github.com/boltdb/bolt"
)

type Profile struct {
	Uid		string
	Name	string
	Phone	string
}

func NewDB(dbPath string) *bolt.DB {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Profile"))
		if (err != nil) {
			log.Fatal(err)
		}
		return nil
	})

	if (err != nil) {
		log.Fatal(err)
	}

	return db
}

func GetProfile(phone string, db *bolt.DB) (*Profile, error) {
	var profile *Profile

	err := db.View(func(tx *bolt.Tx) error {
		p := tx.Bucket([]byte("Profile"))
		acct := p.Get([]byte(phone))

		err := json.Unmarshal(acct, &profile)
		if (err != nil) {
			return err
		}
		return nil
	})

	if (err != nil) {
		return nil, err
	}

	return profile, nil
}

func UpdateProfile(uid string, name string, phone string, db *bolt.DB) (*Profile, error) {
	profile := &Profile{Uid: uid, Name: name, Phone: phone}

	encoded, err := json.Marshal(profile)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		p := tx.Bucket([]byte("Profile"))

		return p.Put([]byte(profile.Phone), encoded)
	})

	if (err != nil) {
		return nil, err
	}
	println("returning profile ", profile)
	return profile, err
}

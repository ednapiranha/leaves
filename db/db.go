package db

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

type Profile struct {
	Uid		string
	Name	string
	Phone	string
}

type Strains struct {
	Data []*Strain
	Meta struct {
		Pagination struct {
			Total			int `json:"total"`
			Count			int `json:"count"`
			Per_Page		int `json:"per_page"`
			Current_Page	int `json:"current_page"`
			Total_Pages 	int `json:"total_pages"`
			Links struct {
				Next		string `json:"next"`
			}
		}
	}
}

type StrainData struct {
	Data []Strain
}

type Strain struct {
	Name	string `json:"name"`
	Ucpc	string `json:"ucpc"`
	Link	string `json:"link"`
	Qr		string `json:"-"`
	Url		string `json:"url"`
	Image	string `json:"image"`
	/*
	SeedCompany struct {
		Name	string `json:"name"`
		Ucpc	string `json:"ucpc"`
		Link	string `json:"link"`
	}
	Genetics struct {
		Names	string `json:"names"`
		Ucpc	string `json:"ucpc"`
		Link	string `json:"link"`
	}
	*/
	SeedCompany	map[string]interface{} `json:"seedCompany"`
	Genetics	map[string]interface{} `json:"genetics"`
	Lineage		map[string]interface{} `json:"lineage"`
	Children	map[string]interface{} `json:"-"`
	Reviews		map[string]interface{} `json:"-"`
	CreatedAt	map[string]interface{} `json:"createdAt"`
	UpdatedAt	map[string]interface{} `json:"updatedAt"`
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

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Strain"))
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

func UpdateStrain(strain *Strain, db *bolt.DB) error {
	encoded, err := json.Marshal(strain)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("* added: ", strain.Name, strain.Ucpc)

	err = db.Update(func(tx *bolt.Tx) error {
		s := tx.Bucket([]byte("Strain"))

		return s.Put([]byte(strain.Ucpc), encoded)
	})

	if (err != nil) {
		return err
	}

	return nil
}

func (s *StrainData) AppendStrain(strain Strain) []Strain {
	s.Data = append(s.Data, strain)
	return s.Data
}

func GetAllStrains(db *bolt.DB) (StrainData, error) {
	var strain Strain
	strains := []Strain{}
	strn := StrainData{strains}

	err := db.View(func(tx *bolt.Tx) error {
		s := tx.Bucket([]byte("Strain"))
		c := s.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Println("+++++++++++++++++++++++ ", k)
			err := json.Unmarshal(v, &strain)
			if (err != nil) {
				log.Fatal(err)
			}

			strn.AppendStrain(strain)

			fmt.Println(strn.Data)
		}
		return nil
	})

	if (err != nil) {
		return strn, err
	}

	return strn, nil
}

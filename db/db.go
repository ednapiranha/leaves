package db

import (
	"fmt"
	"log"

	"github.com/asdine/storm"
)

type Profile struct {
	Uid		string `storm:"id"`
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

type Strain struct {
	Name	string `storm:"index"`
	Ucpc	string `storm:"id"`
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
	SeedCompany	map[string]interface{} `storm:"index"`
	Genetics	map[string]interface{} `json:"genetics"`
	Lineage		map[string]interface{} `json:"lineage"`
	Children	map[string]interface{} `json:"-"`
	Reviews		map[string]interface{} `json:"-"`
	CreatedAt	map[string]interface{} `storm:"index"`
	UpdatedAt	map[string]interface{} `storm:"index"`
}

func NewDB(dbPath string) *storm.DB {
	db, err := storm.Open(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func Close(db *storm.DB) {
	db.Close()
}

func GetProfile(phone string, db *storm.DB) (Profile, error) {
	var profile Profile

	err := db.One("Phone", phone, &profile)
	if (err != nil) {
		return profile, err
	}
	return profile, nil
}

func UpdateProfile(profile Profile, db *storm.DB) (Profile, error) {
	err := db.Save(&profile)
	if (err != nil) {
		return profile, err
	}
	return profile, err
}

func UpdateStrain(strain Strain, db *storm.DB) error {
	err := db.Save(&strain)
	if (err != nil) {
		fmt.Println(err)
		return err
	}
	return nil
}

func GetAllStrains(db *storm.DB) ([]Strain, error) {
	var strains []Strain

	err := db.AllByIndex("Name", &strains, storm.Limit(10))
	if (err != nil) {
		return strains, err
	}
	return strains, nil
}

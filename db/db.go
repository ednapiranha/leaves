package db

import (
	"fmt"
	"log"
	"time"

	"github.com/asdine/storm"
)

const limit = 20

type Profile struct {
	Uid		string `storm:"id"`
	Name	string `json:"name"`
	Phone	string `json:"phone"`
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
	CreatedAt	map[string]interface{} `json:"createdAt"`
	UpdatedAt	map[string]interface{} `json:"updatedAt"`
}

type Review struct {
	Rid			string `storm:"id"`
	Uid			string `storm:"index"`
	Group		string `storm:"index"`
	FiveMin		string `json:"fiveMin"`
	TenMin		string `json:"tenMin"`
	FifteenMin	string `json:"fifteenMin"`
	TwentyMin	string `json:"twentyMin"`
	Comments	string `json:"comments"`
	CreatedAt	time.Time `json:"createdAt"`
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

func GetLimit() int {
	return limit
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

func GetAllStrains(page int, db *storm.DB) ([]Strain, error) {
	var strains []Strain

	page = page - 1

	if (page < 0) {
		page = 0
	}

	err := db.AllByIndex("Name", &strains, storm.Limit(limit), storm.Skip(page * limit))
	if (err != nil) {
		return strains, err
	}
	return strains, nil
}

func UpdateReview(reviewStrain Review, reviewFeed Review, db *storm.DB) error {
	tx, err := db.Begin(true)

	err = db.Save(&reviewStrain)
	if (err != nil) {
		tx.Rollback()
		return err
	}

	err = db.Save(&reviewFeed)
	if (err != nil) {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func GetReviewsByStrain(id string, db *storm.DB) ([]Review, error) {
	var reviews []Review

	err := db.Find("Group", "strain", &reviews, storm.Limit(limit))
	if (err != nil) {
		return reviews, err
	}
	return reviews, nil
}

func GetReviewsByUser(uid string, db *storm.DB) ([]Review, error) {
	var reviews []Review

	err := db.Find("Uid", uid, &reviews, storm.Limit(limit))
	if (err != nil) {
		return reviews, err
	}
	return reviews, nil
}

func GetReviewsByFeed(db *storm.DB) ([]Review, error) {
	var reviews []Review

	err := db.Find("Group", "feed", &reviews, storm.Limit(limit))
	if (err != nil) {
		return reviews, err
	}
	return reviews, nil
}

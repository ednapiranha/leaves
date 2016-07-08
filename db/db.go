package db

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/asdine/storm"
	"github.com/nu7hatch/gouuid"
)

const limit = 24

type Profile struct {
	Uid   string `storm:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

/*
type Strains struct {
	Data []*StrainClean
	Meta struct {
		Pagination struct {
			Total        int `json:"total"`
			Count        int `json:"count"`
			Per_Page     int `json:"per_page"`
			Current_Page int `json:"current_page"`
			Total_Pages  int `json:"total_pages"`
			Links        struct {
				Next string `json:"next"`
			}
		}
	}
}
*/
type Strain struct {
	Name            string                 `storm:"index"`
	Ucpc            string                 `storm:"id"`
	Link            string                 `json:"link"`
	Qr              string                 `json:"-"`
	Url             string                 `json:"url"`
	Image           string                 `json:"image"`
	SeedCompany     map[string]interface{} `json:"genetics"`
	SeedCompanyName string                 `storm:"index"`
	SeedCompanyUcpc string                 `storm:"index"`
	Genetics        map[string]interface{} `json:"genetics"`
	GeneticsUcpc    string                 `json:"geneticsUcpc"`
	Lineage         map[string]interface{} `json:"lineage"`
	Children        map[string]interface{} `json:"-"`
	Reviews         map[string]interface{} `json:"-"`
	CreatedAt       map[string]interface{} `json:"createdAt"`
	UpdatedAt       map[string]interface{} `json:"updatedAt"`
}

type Fave struct {
	Rid       string `storm:"index"`
	Uid       string `storm:"index"`
	CreatedAt int32  `storm:"index"`
}

/*
type StrainClean struct {
	Name        string `storm:"index"`
	Ucpc        string `storm:"id"`
	Link        string `json:"link"`
	Qr          string `json:"-"`
	Url         string `json:"url"`
	Image       string `json:"image"`
	SeedCompany struct {
		Name string `storm:"index"`
		Ucpc string `storm:"index"`
		Link string `json:"-"`
	}
	Genetics struct {
		Names string `json:"-"`
		Ucpc  string `json:"ucpc"`
		Link  string `json:"-"`
	}
	Lineage   map[string]interface{} `json:"lineage"`
	Children  map[string]interface{} `json:"-"`
	Reviews   map[string]interface{} `json:"-"`
	CreatedAt map[string]interface{} `json:"createdAt"`
	UpdatedAt map[string]interface{} `json:"updatedAt"`
}
*/

type Review struct {
	Rid        string `storm:"id"`
	Ucpc       string `storm:"index"`
	Uid        string `storm:"index"`
	Strain     string `json:"strain"`
	Username   string `json:"username"`
	IsOwner    bool   `json:"isOwner"`
	Grower     string `storm:"index"`
	FiveMin    string `json:"fiveMin"`
	TenMin     string `json:"tenMin"`
	FifteenMin string `json:"fifteenMin"`
	TwentyMin  string `json:"twentyMin"`
	Comments   string `json:"comments"`
	CreatedAt  int32  `storm:"index"`
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
	if err != nil {
		return profile, err
	}
	return profile, nil
}

func UpdateProfile(profile Profile, db *storm.DB) (Profile, error) {
	err := db.Save(&profile)
	if err != nil {
		return profile, err
	}
	return profile, err
}

func UpdateStrain(strain Strain, db *storm.DB) error {
	err := db.Save(&strain)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func GetAllStrains(page int, db *storm.DB) ([]Strain, error) {
	var strains []Strain

	page = page - 1

	if page < 0 {
		page = 0
	}

	err := db.AllByIndex("Name", &strains, storm.Limit(limit), storm.Skip(page*limit))
	if err != nil {
		return strains, err
	}
	return strains, nil
}

func GetStrain(ucpc string, db *storm.DB) (Strain, error) {
	var strain Strain

	err := db.One("Ucpc", ucpc, &strain)
	if err != nil {
		return strain, err
	}
	return strain, nil
}

func AddReview(review Review, db *storm.DB) error {
	var rev Review

	u, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}

	id := hex.EncodeToString(u[:])

	rev = review
	rev.Rid = id
	rev.CreatedAt = int32(time.Now().Unix())

	err = db.Save(&rev)
	if err != nil {
		return err
	}

	return nil
}

func GetReview(id string, uid string, db *storm.DB) (Review, error) {
	var review Review
	var strain Strain
	var profile Profile

	err := db.One("Rid", id, &review)
	if err != nil {
		return review, err
	}

	err = db.One("Ucpc", review.Ucpc, &strain)
	if err != nil {
		return review, err
	}

	review.Strain = strain.Name

	if review.Uid == uid {
		review.IsOwner = true
	} else {
		review.IsOwner = false
	}

	err = db.One("Uid", review.Uid, &profile)
	if err != nil {
		return review, err
	}

	review.Username = profile.Name

	return review, nil
}

func RemoveReview(id string, uid string, db *storm.DB) (string, error) {
	var review Review

	err := db.One("Rid", id, &review)
	if err != nil {
		return "", err
	}

	if review.Uid != uid {
		return "", errors.New("Not the owner of this review. Cannot delete")
	}

	err = db.Remove(&review)
	if err != nil {
		return "", err
	}

	return review.Ucpc, nil
}

func SearchStrains(name string, page int, db *storm.DB) ([]Strain, error) {
	var strains []Strain

	page = page - 1

	if page < 0 {
		page = 0
	}

	err := db.Range("Name", name, name+"*", &strains, storm.Limit(limit), storm.Skip(page*limit))
	if err != nil {
		return strains, err
	}
	return strains, nil
}

func GetReviewsByStrain(id string, page int, uid string, db *storm.DB) ([]Review, error) {
	var reviews []Review
	var strain Strain

	page = page - 1

	if page < 0 {
		page = 0
	}

	err := db.Find("Ucpc", id, &reviews, storm.Limit(limit), storm.Skip(page*limit))
	if err != nil {
		return reviews, err
	}
	for i := range reviews {
		err = db.One("Ucpc", reviews[i].Ucpc, &strain)
		if err == nil {
			reviews[i].Strain = strain.Name
		}
		if reviews[i].Uid == uid {
			reviews[i].IsOwner = true
		} else {
			reviews[i].IsOwner = false
		}
	}
	return reviews, nil
}

func GetReviewsByGrower(grower string, page int, uid string, db *storm.DB) ([]Review, error) {
	var reviews []Review
	var strain Strain

	page = page - 1

	if page < 0 {
		page = 0
	}

	err := db.Find("Grower", grower, &reviews, storm.Limit(limit), storm.Skip(page*limit))
	if err != nil {
		return reviews, err
	}
	for i := range reviews {
		err = db.One("Ucpc", reviews[i].Ucpc, &strain)
		if err == nil {
			reviews[i].Strain = strain.Name
		}
		if reviews[i].Uid == uid {
			reviews[i].IsOwner = true
		} else {
			reviews[i].IsOwner = false
		}
	}
	return reviews, nil
}

func GetFeedByUser(uid string, page int, db *storm.DB) ([]Review, error) {
	var reviews []Review
	var strain Strain

	page = page - 1

	if page < 0 {
		page = 0
	}

	err := db.Find("Uid", uid, &reviews, storm.Limit(limit), storm.Skip(page*limit))
	if err != nil {
		return reviews, err
	}
	for i := range reviews {
		err = db.One("Ucpc", reviews[i].Ucpc, &strain)
		if err == nil {
			reviews[i].Strain = strain.Name
		}
		if reviews[i].Uid == uid {
			reviews[i].IsOwner = true
		} else {
			reviews[i].IsOwner = false
		}
	}
	return reviews, nil
}

func GetFeed(uid string, page int, db *storm.DB) ([]Review, error) {
	var reviews []Review
	var strain Strain

	page = page - 1

	if page < 0 {
		page = 0
	}

	err := db.AllByIndex("CreatedAt", &reviews, storm.Limit(limit), storm.Skip(page*limit))
	if err != nil {
		return reviews, err
	}
	for i := range reviews {
		fmt.Println(reviews[i].CreatedAt, reviews[i].Grower)
		err = db.One("Ucpc", reviews[i].Ucpc, &strain)
		if err == nil {
			reviews[i].Strain = strain.Name
		}
		if reviews[i].Uid == uid {
			reviews[i].IsOwner = true
		} else {
			reviews[i].IsOwner = false
		}
	}
	return reviews, nil
}

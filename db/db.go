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
const maxOrder = 90000000000

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
	SearchTerm      string                 `storm:"index"`
	ReviewsCount    int                    `json:"reviewsCount"`
	HasReviews      bool                   `json:"hasReviews"`
}

type Fave struct {
	Rid       string `storm:"index"`
	Uid       string `storm:"index"`
	CreatedAt int64  `storm:"index"`
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
	CreatedAt  int64  `storm:"index"`
	OrderId    int64  `storm:"index"`
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

func GetMaxOrder() int64 {
	return maxOrder
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

func GetAllStrainsNoPagination(db *storm.DB) ([]Strain, error) {
	var strains []Strain

	err := db.AllByIndex("Name", &strains)
	if err != nil {
		return strains, err
	}

	return strains, nil
}

func GetAllStrains(page int, db *storm.DB) ([]Strain, int, error) {
	var strains []Strain
	var reviews []Review
	strainsCount := 0

	page = page - 1

	if page < 0 {
		page = 0
	}

	err := db.AllByIndex("Name", &strains, storm.Limit(limit), storm.Skip((page+1)*limit))
	if err != nil {
		return strains, strainsCount, err
	}
	strainsCount = len(strains)

	err = db.AllByIndex("Name", &strains, storm.Limit(limit), storm.Skip(page*limit))
	if err != nil {
		return strains, strainsCount, err
	}

	for i := range strains {
		err = db.Find("Ucpc", strains[i].Ucpc, &reviews)
		if err == nil {
			strains[i].ReviewsCount = len(reviews)
			if strains[i].ReviewsCount > 99 {
				strains[i].ReviewsCount = 99
			}
			strains[i].HasReviews = true
		} else {
			strains[i].ReviewsCount = 0
			strains[i].HasReviews = false
		}
	}

	return strains, strainsCount, nil
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
	rev.CreatedAt = int64(time.Now().Unix())
	rev.OrderId = maxOrder - rev.CreatedAt

	err = db.Save(&rev)
	if err != nil {
		return err
	}

	return nil
}

func UpdateReview(review Review, db *storm.DB) error {
	err := db.Save(&review)
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

func SearchStrains(name string, page int, db *storm.DB) ([]Strain, int, error) {
	var strains []Strain
	var reviews []Review
	strainsCount := 0

	page = page - 1

	if page < 0 {
		page = 0
	}

	err := db.Range("SearchTerm", name, name+"*", &strains, storm.Limit(limit), storm.Skip((page+1)*limit))
	if err != nil {
		return strains, strainsCount, err
	}
	strainsCount = len(strains)

	err = db.Range("SearchTerm", name, name+"*", &strains, storm.Limit(limit), storm.Skip(page*limit))
	if err != nil {
		return strains, strainsCount, err
	}

	for i := range strains {
		err = db.Find("Ucpc", strains[i].Ucpc, &reviews)
		if err == nil {
			strains[i].ReviewsCount = len(reviews)
			if strains[i].ReviewsCount > 99 {
				strains[i].ReviewsCount = 99
			}
			strains[i].HasReviews = true
		} else {
			strains[i].ReviewsCount = 0
			strains[i].HasReviews = false
		}
	}

	return strains, strainsCount, nil
}

func GetReviewsByStrain(id string, page int, uid string, db *storm.DB) ([]Review, int, error) {
	var reviews []Review
	var strain Strain
	reviewsCount := 0

	page = page - 1

	if page < 0 {
		page = 0
	}

	err := db.Find("Ucpc", id, &reviews, storm.Limit(limit), storm.Skip((page+1)*limit))
	if err != nil {
		return reviews, reviewsCount, err
	}
	reviewsCount = len(reviews)

	err = db.Find("Ucpc", id, &reviews, storm.Limit(limit), storm.Skip(page*limit))
	if err != nil {
		return reviews, reviewsCount, err
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

	return reviews, reviewsCount, nil
}

func GetReviewsByGrower(grower string, page int, uid string, db *storm.DB) ([]Review, int, error) {
	var reviews []Review
	var strain Strain
	reviewsCount := 0

	page = page - 1

	if page < 0 {
		page = 0
	}

	err := db.Find("Grower", grower, &reviews, storm.Limit(limit), storm.Skip((page+1)*limit))
	if err != nil {
		return reviews, reviewsCount, err
	}
	reviewsCount = len(reviews)

	err = db.Find("Grower", grower, &reviews, storm.Limit(limit), storm.Skip(page*limit))
	if err != nil {
		return reviews, reviewsCount, err
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
	return reviews, reviewsCount, nil
}

func GetFeedByUser(uid string, page int, db *storm.DB) ([]Review, int, error) {
	var reviews []Review
	var strain Strain
	reviewsCount := 0

	page = page - 1

	if page < 0 {
		page = 0
	}

	err := db.Find("Uid", uid, &reviews, storm.Limit(limit), storm.Skip((page+1)*limit))
	if err != nil {
		return reviews, reviewsCount, err
	}
	reviewsCount = len(reviews)

	err = db.Find("Uid", uid, &reviews, storm.Limit(limit), storm.Skip(page*limit))
	if err != nil {
		return reviews, reviewsCount, err
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
	return reviews, reviewsCount, nil
}

func GetFeed(uid string, page int, db *storm.DB) ([]Review, int, error) {
	var reviews []Review
	var strain Strain
	reviewsCount := 0

	page = page - 1

	if page < 0 {
		page = 0
	}

	err := db.AllByIndex("OrderId", &reviews, storm.Limit(limit), storm.Skip((page+1)*limit))
	if err != nil {
		return reviews, reviewsCount, err
	}
	reviewsCount = len(reviews)

	err = db.AllByIndex("OrderId", &reviews, storm.Limit(limit), storm.Skip(page*limit))
	if err != nil {
		return reviews, reviewsCount, err
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
	return reviews, reviewsCount, nil
}

func GetFeedNoPagination(db *storm.DB) ([]Review, error) {
	var reviews []Review

	err := db.AllByIndex("CreatedAt", &reviews)
	if err != nil {
		return reviews, err
	}

	return reviews, nil
}

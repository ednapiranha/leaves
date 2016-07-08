package main

import (
	//"fmt"
	//"log"
	//"regexp"
	//"strings"

	"github.com/revolting/leaves/db"
)

var d = db.NewDB("./boltdb/leaves.db")

func main() {
	// cleans up strain names for searching
	/*
		strains, _ := db.GetAllStrainsNoPagination(d)

		reg, err := regexp.Compile("[()-/]")
		if err != nil {
			log.Fatal(err)
		}

		for i := range strains {
			strains[i].SearchTerm = reg.ReplaceAllString(strings.ToLower(strains[i].Name), "")
			db.UpdateStrain(strains[i], d)
		}
	*/
	// fix order id for older posts
	reviews, _ := db.GetFeedNoPagination(d)
	for i := range reviews {
		reviews[i].OrderId = db.GetMaxOrder() - reviews[i].CreatedAt
		db.UpdateReview(reviews[i], d)
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/revolting/leaves/db"
)

var strains *db.Strains

var d = db.NewDB("./boltdb/leaves.db")
var currUrl = "https://www.cannabisreports.com/api/v1.0/strains?page=1"

func main() {
	GetData(currUrl)
	defer d.Close()

	rate := time.Millisecond * 1500
	throttle := time.Tick(rate)

	for strains.Meta.Pagination.Total_Pages > strains.Meta.Pagination.Current_Page {
		<-throttle
		fmt.Println("adding next ", currUrl)
		go GetData(currUrl)
	}
}

func GetData(url string) {
	r, _ := http.Get(url)
	defer r.Body.Close()

	json.NewDecoder(r.Body).Decode(&strains)
	currUrl = strains.Meta.Pagination.Links.Next
	Update()
}

func Update() {
	for _, v := range strains.Data{
		err := db.UpdateStrain(*v, d)
		if (err != nil) {
			fmt.Println("could not update strain ", v.Name)
		}
	}
}

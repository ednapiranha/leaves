package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	//"reflect"

	"github.com/revolting/leaves/db"
)

var strains *db.Strains

var d = db.NewDB("./boltdb/leaves.db")

func main() {
	err := GetData("https://www.cannabisreports.com/api/v1.0/strains")
	if (err != nil) {
		log.Fatal(err)
	}
	//fmt.Println(strains.Data)

	for _, v := range strains.Data{
		err := db.UpdateStrain(&v, d)
		if (err != nil) {
			fmt.Println("could not update strain ", v.Name)
		}
	}
}

func GetData(url string) error {
	r, err := http.Get(url)
	if (err != nil) {
		log.Fatal(err)
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(&strains)
	//db.UpdateStrain(strain *Strain, d)
}

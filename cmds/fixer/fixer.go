package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/revolting/leaves/db"
)

var strains *db.Strains
var strClean db.StrainClean

var d = db.NewDB("./boltdb/leaves.db")

func main() {
	s, err := db.GetAllStrains(0, d)
	if err != nil {
		log.Fatal(err)
	}

	for i := range s {
		str, _ := json.Marshal(s[i].Genetics)
		_ = json.Unmarshal([]byte(str), &strClean)

		if reflect.ValueOf(strClean.Genetics.Names).Kind() != reflect.String {
			strClean.Genetics.Ucpc = ""
			fmt.Println("converting names to string")
		}

		if reflect.ValueOf(strClean.Genetics.Ucpc).Kind() != reflect.String {
			strClean.Genetics.Names = ""
			fmt.Println("converting ucpc to string")
		}

		s[i].GeneticsNames = strClean.Genetics.Names
		s[i].GeneticsUcpc = strClean.Genetics.Ucpc

		db.UpdateStrain(s[i], d)
	}

}

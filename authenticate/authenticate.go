package authenticate

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"github.com/asdine/storm"
	"github.com/nu7hatch/gouuid"
	"github.com/patrickmn/go-cache"
	"github.com/sfreiberg/gotwilio"

	"github.com/revolting/leaves/db"
)

// 5 minute cache for pin
var c = cache.New(5*time.Minute, 30*time.Second)

func init() {
	rand.Seed(time.Now().Unix())
}

func generatePin() string {
	pin := ""
	i := 0

	for ; i < 4; i++ {
		pin = pin + strconv.Itoa(rand.Intn(10))
	}
	println("pin:", pin)
	return pin
}

func fixPhone(phone string) string {
	reg, err := regexp.Compile("[^0-9+]")
	if err != nil {
		log.Fatal(err)
	}

	number := reg.ReplaceAllString(phone, "")
	regNA, err := regexp.MatchString("^[0-9]{10}$", number)
	if err != nil {
		log.Fatal(err)
	}

	regPl, err := regexp.MatchString("^+", number)
	if err != nil {
		log.Fatal(err)
	}

	if regNA {
		number = "+1" + number
	} else if regPl {
		number = "+" + number
	}
	println(number)
	return number
}

func SendPin(twilioSid string, twilioToken string, twilioPhone string, phone string) string {
	twilio := gotwilio.NewTwilioClient(twilioSid, twilioToken)
	from := twilioPhone
	to := fixPhone(phone)
	pin := generatePin()
	message := "Your PIN: " + pin
	twilio.SendSMS(from, to, message, "", "")

	c.Set(to, pin, cache.DefaultExpiration)
	return to
}

func ValidatePin(pin string, phone string) bool {
	phonePin, found := c.Get(phone)
	if found {
		p := phonePin.(string)
		if p == pin {
			return true
		}
	}
	return false
}

func CreateProfile(phone string, d *storm.DB) (db.Profile, error) {
	u, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}

	id := hex.EncodeToString(u[:])
	hash := md5.Sum([]byte(phone))
	phoneHash := hex.EncodeToString(hash[:])

	profile, err := db.GetProfile(phoneHash, d)

	if err != nil {
		fmt.Println("NEW USER")
		p := &db.Profile{Uid: id, Name: "???", Phone: phoneHash}
		profile, err = db.UpdateProfile(*p, d)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("FOUND USER ", profile, "\n", profile)
	}

	return profile, nil
}

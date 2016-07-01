package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/csrf"

	"github.com/revolting/leaves/authenticate"
	"github.com/revolting/leaves/db"
)

func Index(w http.ResponseWriter, req *http.Request) {
	session, err := s.Get(req, "uid")
	if (err != nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s := false

	if (session.Values["uid"] != nil) {
		fmt.Println(*session)
		s = true
	}

	r.HTML(w, http.StatusOK, "index", map[string]interface{}{
		"session": s,
	})
}

func Profile(w http.ResponseWriter, req *http.Request) {
	session, err := s.Get(req, "uid")
	if (err != nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s := false

	if (session.Values["uid"] != nil) {
		fmt.Println(*session)
		s = true
	}

	if (session.Values["uid"] == nil) {
		http.Redirect(w, req, "/", 301)
	}

	if (req.Method == "POST") {
		name := req.FormValue("name")
		uid := session.Values["uid"].(string)
		phone := session.Values["phone"].(string)

		p := &db.Profile{Uid: uid, Name: name, Phone: phone}
		profile, err := db.UpdateProfile(*p, d)
		if (err != nil) {
			log.Fatal(err)
		}

		if (err != nil) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session.Values["name"] = profile.Name;
		session.Save(req, w)
	}

	r.HTML(w, http.StatusOK, "profile", map[string]interface{}{
		"session": s,
		"uid": session.Values["uid"],
		"name": session.Values["name"],
		csrf.TemplateTag: csrf.TemplateField(req),
	})
}

func Directory(w http.ResponseWriter, req *http.Request) {
	page := 1
	prev := "1"

	p := req.URL.Query().Get("page")
	if (p != "") {
		pg, err := strconv.Atoi(p)
		if (err != nil) {
			page = 1
		} else {
			page = pg
		}
	}

	next := strconv.Itoa(page + 1)
	prevInt := page - 1

	if (prevInt >= 1) {
		prev = strconv.Itoa(prevInt)
	}

	strains, err := db.GetAllStrains(page, d)
	if (err != nil) {
		log.Fatal(err)
	}

	if (len(strains) < db.GetLimit()) {
		next = strconv.Itoa(page)
	}

	fmt.Println(len(strains))

	r.HTML(w, http.StatusOK, "directory", map[string]interface{}{
		"session": s,
		"strains": strains,
		"prev": prev,
		"next": next,
	})
}

func Authenticate(w http.ResponseWriter, req *http.Request) {
	if (req.Method == "POST") {
		session, err := s.Get(req, "uid")
		if (err != nil) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		decoder := req.FormValue("phone")
		phone := authenticate.SendPin(*twilioSid, *twilioToken, *twilioPhone, decoder)
		session.Values["phone"] = phone;
		session.Save(req, w)

		http.Redirect(w, req, "/validate", 301)
	} else {
		r.HTML(w, http.StatusOK, "authenticate", map[string]interface{}{
			csrf.TemplateTag: csrf.TemplateField(req),
		})
	}
}

func Validate(w http.ResponseWriter, req *http.Request) {
	if (req.Method == "POST") {
		session, err := s.Get(req, "uid")
		if (err != nil) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pin := req.FormValue("pin")
		phone := session.Values["phone"].(string)
		pinVerify := authenticate.ValidatePin(pin, phone)

		if (pinVerify) {
			profile, err := authenticate.CreateProfile(phone, d)
			if (err != nil) {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			session.Values["phone"] = profile.Phone
			session.Values["uid"] = profile.Uid
			session.Values["name"] = profile.Name
			session.Save(req, w)

			http.Redirect(w, req, "/", 301)
		} else {
			r.HTML(w, http.StatusOK, "validate", nil)
		}
	} else {
		r.HTML(w, http.StatusOK, "validate", map[string]interface{}{
			csrf.TemplateTag: csrf.TemplateField(req),
		})
	}
}

func Logout(w http.ResponseWriter, req *http.Request) {
	session, err := s.Get(req, "uid")
	if (err != nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["phone"] = nil
	session.Values["uid"] = nil
	session.Values["name"] = nil
	session.Save(req, w)
	http.Redirect(w, req, "/", 301)
}

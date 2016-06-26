package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/csrf"

	"app/authenticate"
	"app/db"
	"app/utils"
)

var r = utils.GetRender()
var s = utils.GetSession()

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

		profile, err := db.UpdateProfile(uid, name, phone)
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
	r.HTML(w, http.StatusOK, "directory", nil)
}

func Authenticate(w http.ResponseWriter, req *http.Request) {
	if (req.Method == "POST") {
		session, err := s.Get(req, "uid")
		if (err != nil) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		decoder := req.FormValue("phone")
		phone := authenticate.SendPin(decoder)
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
			profile, err := authenticate.CreateProfile(phone)
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

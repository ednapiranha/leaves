package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/csrf"

	"app/authenticate"
	"app/utils"
)

var r = utils.GetRender()
var s = utils.GetSession()

func Index(w http.ResponseWriter, req *http.Request) {
	session, err := s.Get(req, "phone")
	if (err != nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s := false

	if (session.Values["phone"] != nil) {
		fmt.Println(*session)
		s = true
	}

	r.HTML(w, http.StatusOK, "index", map[string]interface{}{
		"session": s,
	})
}

func Directory(w http.ResponseWriter, req *http.Request) {
	r.HTML(w, http.StatusOK, "directory", nil)
}

func Authenticate(w http.ResponseWriter, req *http.Request) {
	if (req.Method == "POST") {
		session, err := s.Get(req, "phone")
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
		session, err := s.Get(req, "phone")
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
	session, err := s.Get(req, "phone")
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
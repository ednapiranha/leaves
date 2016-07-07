package main

import (
	"net/http"
	"strconv"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"

	"github.com/revolting/leaves/authenticate"
	"github.com/revolting/leaves/db"
)

func Index(w http.ResponseWriter, req *http.Request) {
	uid := ""

	session, err := s.Get(req, "uid")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s := false

	if session.Values["uid"] != nil {
		uid = session.Values["uid"].(string)
		s = true
	}

	reviews, _ := db.GetFeed(uid, d)

	r.HTML(w, http.StatusOK, "index", map[string]interface{}{
		"session": s,
		"uid":     uid,
		"reviews": reviews,
	})
}

func Profile(w http.ResponseWriter, req *http.Request) {
	uid := ""
	session, err := s.Get(req, "uid")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s := false

	if session.Values["uid"] != nil {
		s = true
		uid = session.Values["uid"].(string)
	} else {
		http.Redirect(w, req, "/", 301)
	}

	if req.Method == "POST" {
		name := req.FormValue("name")
		phone := session.Values["phone"].(string)

		p := &db.Profile{Uid: uid, Name: name, Phone: phone}
		profile, err := db.UpdateProfile(*p, d)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session.Values["name"] = profile.Name
		session.Save(req, w)
	}

	reviews, _ := db.GetFeedByUser(uid, d)

	r.HTML(w, http.StatusOK, "profile", map[string]interface{}{
		"session":        s,
		"uid":            uid,
		"name":           session.Values["name"],
		"reviews":        reviews,
		csrf.TemplateTag: csrf.TemplateField(req),
	})
}

func Directory(w http.ResponseWriter, req *http.Request) {
	var strains []db.Strain

	session, err := s.Get(req, "uid")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s := false

	if session.Values["uid"] != nil {
		s = true
	}

	page := 1
	prev := "1"

	p := req.URL.Query().Get("page")
	if p != "" {
		pg, err := strconv.Atoi(p)
		if err != nil {
			page = 1
		} else {
			page = pg
		}
	}

	next := strconv.Itoa(page + 1)
	prevInt := page - 1
	name := ""

	if prevInt >= 1 {
		prev = strconv.Itoa(prevInt)
	}

	if req.Method == "POST" || len(req.URL.Query().Get("keyword")) > 1 {
		if req.Method == "POST" {
			name = req.FormValue("name")
		} else {
			name = req.URL.Query().Get("keyword")
		}
		strains, _ = db.SearchStrains(name, page, d)
	} else {
		strains, _ = db.GetAllStrains(page, d)
	}

	if len(strains) < db.GetLimit() {
		next = strconv.Itoa(page)
	}

	r.HTML(w, http.StatusOK, "directory", map[string]interface{}{
		"search":         name,
		"session":        s,
		"strains":        strains,
		"prev":           prev,
		"next":           next,
		csrf.TemplateTag: csrf.TemplateField(req),
	})
}

func StrainDetail(w http.ResponseWriter, req *http.Request) {
	uid := ""

	session, err := s.Get(req, "uid")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s := false

	if session.Values["uid"] != nil {
		uid = session.Values["uid"].(string)
		s = true
	}

	vars := mux.Vars(req)
	st, err := db.GetStrain(vars["ucpc"], d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	reviews, _ := db.GetReviewsByStrain(vars["ucpc"], uid, d)

	r.HTML(w, http.StatusOK, "strain", map[string]interface{}{
		"session":        s,
		"strain":         st,
		"reviews":        reviews,
		csrf.TemplateTag: csrf.TemplateField(req),
	})
}

func UpdateReview(w http.ResponseWriter, req *http.Request) {
	session, err := s.Get(req, "uid")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if session.Values["uid"] == nil {
		http.Redirect(w, req, "/", 301)
	}

	vars := mux.Vars(req)

	review := &db.Review{
		Uid:        session.Values["uid"].(string),
		Ucpc:       vars["ucpc"],
		Grower:     req.FormValue("grower"),
		FiveMin:    req.FormValue("fiveMin"),
		TenMin:     req.FormValue("tenMin"),
		FifteenMin: req.FormValue("fifteenMin"),
		TwentyMin:  req.FormValue("twentyMin"),
		Comments:   req.FormValue("comments"),
	}

	err = db.AddReview(*review, d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, "/strain/"+vars["ucpc"], 301)
}

func GetReview(w http.ResponseWriter, req *http.Request) {
	uid := ""

	session, err := s.Get(req, "uid")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s := false

	if session.Values["uid"] != nil {
		uid = session.Values["uid"].(string)
		s = true
	}

	vars := mux.Vars(req)

	review, err := db.GetReview(vars["rid"], uid, d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	r.HTML(w, http.StatusOK, "review", map[string]interface{}{
		"session":        s,
		"review":         review,
		csrf.TemplateTag: csrf.TemplateField(req),
	})
}

func DeleteReview(w http.ResponseWriter, req *http.Request) {
	session, err := s.Get(req, "uid")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if session.Values["uid"] == nil {
		http.Redirect(w, req, "/", 301)
		return
	}

	uid := session.Values["uid"].(string)
	vars := mux.Vars(req)

	ucpc, err := db.RemoveReview(vars["rid"], uid, d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	http.Redirect(w, req, "/strain/"+ucpc, 301)
}

func Authenticate(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		session, err := s.Get(req, "uid")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		decoder := req.FormValue("phone")
		phone := authenticate.SendPin(*twilioSid, *twilioToken, *twilioPhone, decoder)
		session.Values["phone"] = phone
		session.Save(req, w)

		http.Redirect(w, req, "/validate", 301)
	} else {
		r.HTML(w, http.StatusOK, "authenticate", map[string]interface{}{
			csrf.TemplateTag: csrf.TemplateField(req),
		})
	}
}

func Validate(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		session, err := s.Get(req, "uid")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pin := req.FormValue("pin")
		phone := session.Values["phone"].(string)
		pinVerify := authenticate.ValidatePin(pin, phone)

		if pinVerify {
			profile, err := authenticate.CreateProfile(phone, d)
			if err != nil {
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
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["phone"] = nil
	session.Values["uid"] = nil
	session.Values["name"] = nil
	session.Save(req, w)
	http.Redirect(w, req, "/", 301)
}

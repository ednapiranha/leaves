package main

import (
	"net/http"
	"strconv"
	"unicode/utf8"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"

	"github.com/revolting/leaves/authenticate"
	"github.com/revolting/leaves/db"
)

func setPage(p string) (string, string, int) {
	page := 1
	prev := "1"

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

	if prevInt >= 1 {
		prev = strconv.Itoa(prevInt)
	}

	return prev, next, page
}

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

	prev, next, page := setPage(req.URL.Query().Get("page"))
	reviews, nextCount, _ := db.GetFeed(uid, page, d)

	if len(reviews) < db.GetLimit() {
		next = strconv.Itoa(page)
	}

	showNext := true
	showPrev := true

	if nextCount == 0 {
		showNext = false
	}

	if page == 1 {
		showPrev = false
	}

	r.HTML(w, http.StatusOK, "index", map[string]interface{}{
		"session":  s,
		"uid":      uid,
		"reviews":  reviews,
		"prev":     prev,
		"next":     next,
		"showNext": showNext,
		"showPrev": showPrev,
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

	prev, next, page := setPage(req.URL.Query().Get("page"))

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

	reviews, nextCount, _ := db.GetFeedByUser(uid, page, d)

	if len(reviews) < db.GetLimit() {
		next = strconv.Itoa(page)
	}

	showNext := true
	showPrev := true

	if nextCount == 0 {
		showNext = false
	}

	if page == 1 {
		showPrev = false
	}

	r.HTML(w, http.StatusOK, "profile", map[string]interface{}{
		"session":        s,
		"uid":            uid,
		"name":           session.Values["name"],
		"reviews":        reviews,
		"prev":           prev,
		"next":           next,
		"showNext":       showNext,
		"showPrev":       showPrev,
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

	prev, next, page := setPage(req.URL.Query().Get("page"))
	name := ""
	nextCount := 0

	if req.Method == "POST" || utf8.RuneCountInString(req.URL.Query().Get("keyword")) > 0 {
		if req.Method == "POST" {
			name = req.FormValue("name")
		} else {
			name = req.URL.Query().Get("keyword")
		}
		strains, nextCount, _ = db.SearchStrains(name, page, d)
	} else {
		strains, nextCount, _ = db.GetAllStrains(page, d)
	}

	if len(strains) < db.GetLimit() {
		next = strconv.Itoa(page)
	}

	showNext := true
	showPrev := true

	if nextCount == 0 {
		showNext = false
	}

	if page == 1 {
		showPrev = false
	}

	r.HTML(w, http.StatusOK, "directory", map[string]interface{}{
		"search":         name,
		"session":        s,
		"strains":        strains,
		"prev":           prev,
		"next":           next,
		"showNext":       showNext,
		"showPrev":       showPrev,
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

	prev, next, page := setPage(req.URL.Query().Get("page"))
	reviews, nextCount, _ := db.GetReviewsByStrain(vars["ucpc"], page, uid, d)

	showNext := true
	showPrev := true

	if nextCount == 0 {
		showNext = false
	}

	if page == 1 {
		showPrev = false
	}

	r.HTML(w, http.StatusOK, "strain", map[string]interface{}{
		"session":        s,
		"strain":         st,
		"reviews":        reviews,
		"prev":           prev,
		"next":           next,
		"showNext":       showNext,
		"showPrev":       showPrev,
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

func UpdateLike(w http.ResponseWriter, req *http.Request) {
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

	err = db.UpdateLike(vars["rid"], uid, d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
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

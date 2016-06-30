package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var m = mux.NewRouter().StrictSlash(true)
var w = httptest.NewRecorder()

func TestGetIndex(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	m.HandleFunc("/", Index).Methods("GET")
	m.ServeHTTP(w, r)

	fmt.Println(w)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetDirectory(t *testing.T) {
	r, err := http.NewRequest("GET", "/directory", nil)
	if err != nil {
		t.Fatal(err)
	}

	m.HandleFunc("/directory", Directory).Methods("GET")
	m.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetAuthenticate(t *testing.T) {
	r, err := http.NewRequest("GET", "/authenticate", nil)
	if err != nil {
		t.Fatal(err)
	}

	m.HandleFunc("/authenticate", Authenticate).Methods("GET")
	m.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetProfile(t *testing.T) {
	r, err := http.NewRequest("GET", "/profile", nil)
	if err != nil {
		t.Fatal(err)
	}

	m.HandleFunc("/profile", Profile).Methods("GET")
	m.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

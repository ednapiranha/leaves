package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var (
	httpPort 	= flag.String("port", ":8080", "Listen address")
	serverEnv 	= flag.Bool("isDev", true, "Server environment mode")
	twilioSid 	= flag.String("twilioSid", "111", "Twilio SID")
	twilioToken = flag.String("twilioToken", "111", "Twilio Token")
	twilioPhone	= flag.String("twilioPhone", "+15555555", "Twilio Phone Number")

	r			= render.New(render.Options{
					Directory: "templates",
					Extensions: []string{".html"},
					IsDevelopment: *serverEnv,
				})
)

func main() {
	flag.Parse()

	router := NewRouter()
	router.PathPrefix("/media/").Handler(http.StripPrefix("/media/", http.FileServer(http.Dir("./media/"))))

	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.UseHandler(router)

	log.Fatal(http.ListenAndServe(*httpPort, n))
}

func Index(w http.ResponseWriter, req *http.Request) {
	r.HTML(w, http.StatusOK, "index", nil)
}

func Directory(w http.ResponseWriter, req *http.Request) {
	r.HTML(w, http.StatusOK, "directory", nil)
}

func Authenticate(w http.ResponseWriter, req *http.Request) {
	if (req.Method == "POST") {
		decoder := req.FormValue("phone")
		pin, err := sendPin(decoder)
		if (err != nil) {
			log.Fatal(err)
		}
	}
	r.HTML(w, http.StatusOK, "authenticate", nil)
}

func Validate(w http.ResponseWriter, req *http.Request) {
	r.HTML(w, http.StatusOK, "validate", nil)
}

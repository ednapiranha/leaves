package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
	"github.com/unrolled/secure"
	"github.com/urfave/negroni"

	"github.com/revolting/leaves/db"
)

var (
	httpPort		= flag.String("port", ":8080", "Listen address")
	isDev			= flag.Bool("isDev", true, "Server environment mode")
	twilioSid		= flag.String("twilioSid", "111", "Twilio SID")
	twilioToken		= flag.String("twilioToken", "111", "Twilio token")
	twilioPhone		= flag.String("twilioPhone", "+15555555", "Twilio phone number")
	cookieSecret	= flag.String("cookie", "secret", "Session cookie secret")
	csrfSecret		= flag.String("csrfSecret", "something-that-is-32-bytes------", "CSRF secret")
	dbPath			= flag.String("db", "./boltdb/leaves.db", "Database path")
	s				= sessions.NewCookieStore([]byte(*cookieSecret))

	r				= render.New(render.Options{
						Directory: "templates",
						Extensions: []string{".tmpl"},
						Layout: "layout",
						IsDevelopment: *isDev,
					})
	d				= db.NewDB(*dbPath)
)

func main() {
	flag.Parse()

	router := NewRouter()
	router.PathPrefix("/media/").Handler(http.StripPrefix("/media/",
		http.FileServer(http.Dir("./media/"))))

	csrf := csrf.Protect(
		[]byte(*csrfSecret),
		csrf.Secure(!*isDev),
	)

	csp := secure.New(secure.Options{
		AllowedHosts: []string{"localhost" + *httpPort, "leaves.revolting.me", "fonts.googleapis.com"},
		FrameDeny: true,
		IsDevelopment: *isDev,
	})

	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.Use(negroni.HandlerFunc(csp.HandlerFuncWithNext))
	n.UseHandler(router)

	log.Fatal(http.ListenAndServe(*httpPort, csrf(n)))
}

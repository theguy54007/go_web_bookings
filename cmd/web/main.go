package main

import (
	"encoding/gob"
	"fmt"
	"go_udemy/bookings/internal/config"
	"go_udemy/bookings/internal/handlers"
	"go_udemy/bookings/internal/helpers"
	"go_udemy/bookings/internal/models"
	renders "go_udemy/bookings/internal/render"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
	// gob.Register(models.Reservation{})

	// app.InProduction = false

	// session = scs.New()
	// session.Lifetime = 24 * time.Hour
	// session.Cookie.Persist = true
	// session.Cookie.SameSite = http.SameSiteLaxMode
	// session.Cookie.Secure = app.InProduction

	// app.Session = session

	// tc, err := render.CreateTemplateCache()
	// if err != nil {
	// 	log.Fatal("cannot create template cache")
	// }

	// app.TemplateCache = tc
	// app.UseCache = false

	// repo := handlers.NewRepo(&app)
	// handlers.NewHandlers(repo)

	// render.NewTemplates(&app)

	fmt.Printf("Starting application on port %s\n", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() error {
	gob.Register(models.Reservation{})

	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := renders.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	renders.NewTemplates(&app)
	helpers.NewHelpers(&app)

	return nil
}

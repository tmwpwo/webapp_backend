package main

import (
	"fmt"
	"go_server/pkg/config"
	"go_server/pkg/driver"
	"go_server/pkg/handlers"
	"go_server/pkg/renders"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

var port = ":8080"
var app config.AppConfig
var session *scs.SessionManager

func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	fmt.Printf("strting app on port %s", port)

	srv := &http.Server{
		Addr:    port,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database
	log.Println("connecting to database..........")

	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=test user=postgres password=")
	if err != nil {
		log.Fatal("cannot connect to database! bruh...")
		return db, err
	}
	log.Println("connectd to dataabse")

	templateCache, err := renders.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return db, err
	}

	app.TemplateCache = templateCache

	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	renders.NewTemplates(&app)

	return db, nil
}

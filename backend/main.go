package main

import (
	"flag"
	"log"

	"github.com/inchingforward/logbook/backend/handlers"
	"github.com/inchingforward/logbook/backend/models"
	"github.com/inchingforward/logbook/backend/view"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	_ "github.com/lib/pq"
)

func init() {
	var err error

	models.DB, err = sqlx.Connect("postgres", "user=logbook dbname=logbook sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	debug := false

	flag.BoolVar(&debug, "debug", false, "true to enable debug")
	flag.Parse()

	e := echo.New()

	view.SetRenderer(e, debug)

	e.Static("/static", "static")

	handlers.AddHandlers(e)

	log.Println("Listening on 4003...")
	e.Run(standard.New(":4003"))
}

package main

import (
	"flag"
	"log"

	"github.com/inchingforward/logbook/server/handlers"
	"github.com/inchingforward/logbook/server/view"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
)

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

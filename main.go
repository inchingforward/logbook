package main

import (
	"net/http"

	"github.com/labstack/echo"
)

func notYetImplemented(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented)
}

func main() {
	e := echo.New()

	e.GET("/", notYetImplemented)
	e.GET("/about", notYetImplemented)
	e.GET("/contact", notYetImplemented)
	e.GET("/login", notYetImplemented)
	e.POST("/login", notYetImplemented)
	e.POST("/logout", notYetImplemented)
	e.GET("/logbook", notYetImplemented)
	e.GET("/logbook/add", notYetImplemented)
	e.POST("/logbook/add", notYetImplemented)
	e.GET("/logbook/:uuid", notYetImplemented)
	e.POST("/logbook/:uuid", notYetImplemented)
	e.POST("/fetchtitle", notYetImplemented)
	e.GET("/users/:username", notYetImplemented)
	e.GET("/users/:username/logbook", notYetImplemented)
	e.GET("/users/:username/logbook/:uuid", notYetImplemented)

	e.Logger.Fatal(e.Start(":8000"))
}

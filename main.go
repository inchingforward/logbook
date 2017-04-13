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
	e.GET("/bookmarks", notYetImplemented)
	e.GET("/bookmarks/add", notYetImplemented)
	e.POST("/bookmarks/add", notYetImplemented)
	e.GET("/bookmarks/:uuid", notYetImplemented)
	e.POST("/bookmarks/:uuid", notYetImplemented)
	e.POST("/fetchtitle", notYetImplemented)
	e.GET("/users/:username", notYetImplemented)
	e.GET("/users/:username/bookmarks", notYetImplemented)
	e.GET("/users/:username/bookmarks/:uuid", notYetImplemented)

	e.Logger.Fatal(e.Start(":8000"))
}

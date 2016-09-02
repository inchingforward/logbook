package handlers

import (
	"fmt"
	"net/http"

	"github.com/inchingforward/logbook/backend/view"
	"github.com/labstack/echo"
)

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AddHandlers add the Logbook handler functions to the Echo engine.
func AddHandlers(e *echo.Echo) {
	e.GET("/", index)
	e.GET("/about", about)
	e.POST("/login", login)
}

func index(c echo.Context) error {
	return view.RenderTemplate(c, "index.html")
}

func about(c echo.Context) error {
	return view.RenderTemplate(c, "about.html")
}

func login(c echo.Context) error {
	login := new(Login)
	if err := c.Bind(login); err != nil {
		return err
	}

	if login.Username == "" || login.Password == "" {
		return c.String(http.StatusOK, "ERROR: Invalid login")
	}

	return c.String(http.StatusOK, fmt.Sprintf("You sent '%v'", login))
}

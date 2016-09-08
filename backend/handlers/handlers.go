package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/inchingforward/logbook/backend/models"
	"github.com/inchingforward/logbook/backend/view"
	"github.com/labstack/echo"
)

var secret string

// A Login holds a username and password a user attempted
// to login with.
type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// A Result holds the result of making a Logbook REST call.  If
// the call succeeded, Success will be true and Message is optional. If
// the call failed, Success will be false and Message will contain
// an error message.
type Result struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func init() {
	secret := os.Getenv("LOGBOOK_SECRET")
	if secret == "" {
		log.Fatal("Required environment variable LOGBOOK_SECRET missing")
	}
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
		return c.JSON(http.StatusOK, &Result{false, "Invalid login"})
	}

	log.Printf("Login for user %v", login.Username)
	user, err := models.Authenticate(login.Username, login.Password)
	if err != nil {
		return err
	}

	log.Printf("user %v successfully logged in!\n", user)

	return c.JSON(http.StatusOK, &Result{true, "Logged in"})
}

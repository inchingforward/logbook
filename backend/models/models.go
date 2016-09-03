package models

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// DB represents a connection to the database.
var DB *sqlx.DB

// A LogbookUser is an account that has been created by a user.
type LogbookUser struct {
	ID          uint64      `db:"id"`
	Username    string      `db:"username"`
	Password    string      `db:"password"`
	DisplayName string      `db:"display_name"`
	Active      bool        `db:"active"`
	CreatedAt   time.Time   `db:"created_at"`
	LastLoginAt pq.NullTime `db:"last_login_at"`
}

// A LogbookEntry is a captured note or URL made by a LogbookUser.
type LogbookEntry struct {
	ID        uint64    `db:"id"`
	Title     string    `db:"title"`
	URL       string    `db:"url"`
	Notes     string    `db:"notes"`
	Private   bool      `db:"private"`
	UserID    uint64    `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	EntryAt   time.Time `db:"entry_at"`
	tags      string    `db:"tags"`
}

// Authenticate returns the LogbookUser corresponding to username and password.
func Authenticate(username, password string) (*LogbookUser, error) {
	user := &LogbookUser{}

	err := DB.Get(user, "select * from logbook_user where username = $1", username)

	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, errors.New("Invalid user or password")
	}

	return user, err
}

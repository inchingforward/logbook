package models

import (
	"time"
)

type LogbookUser struct {
	ID          uint64    `db:"id"`
	Username    string    `db:"username"`
	Password    string    `db:"password"`
	DisplayName string    `db:"display_name"`
	Active      bool      `db:"active"`
	CreatedAt   time.Time `db:"created_at"`
	LastLoginAt time.Time `db:"last_login_at"`
}

type LogbookEntry struct {
	ID        uint64    `db:"id"`
	Title     string    `db:"title"`
	Url       string    `db:"url"`
	Notes     string    `db:"notes"`
	Private   bool      `db:"private"`
	UserID    uint64    `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	EntryAt   time.Time `db:"entry_at"`
	tags      string    `db:"tags"`
}

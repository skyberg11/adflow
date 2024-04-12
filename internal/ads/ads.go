package ads

import (
	"errors"
	"time"
)

type Ad struct {
	ID           int64
	Title        string `validate:"min:1;max:100"`
	Text         string `validate:"min:1;max:500"`
	AuthorID     int64
	Published    bool
	CreationTime time.Time
	UpdateTime   time.Time
}

type User struct {
	FirstName  string `validate:"min:1;max:100"`
	SecondName string `validate:"min:1;max:100"`
	Nickname   string `validate:"min:1;max:100"`
	Password   string `validate:"min:1;max:100"`
	Email      string `validate:"min:1;max:100"`
	Phone      string `validate:"min:1;max:100"`
	ID         int64
}

type Filter struct {
	Published    any
	AuthorID     any
	TitlePrefix  any
	CreationTime any
}

var ErrBadRequest = errors.New("BadRequest")
var ErrAccessDenied = errors.New("Forbidden")

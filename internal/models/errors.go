package models

import (
	"errors"
)

var ErrNoRecord = errors.New("Models: No records were found")

var ErrInvalidCredentials = errors.New("Models: Invalid credentials")

var ErrDuplicateEmail = errors.New("Models: Duplicate email")

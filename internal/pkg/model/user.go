package model

import (
	"encoding/json"

	"github.com/mono83/maybe"
	maybejson "github.com/mono83/maybe/json"
)

const (
	sexMale   = "М"
	sexFemale = "Ж"
)

// Sex ...
type Sex int32

// SexFromString ...
func SexFromString(s string) Sex {
	switch s {
	case sexMale:
		return 1
	case sexFemale:
		return 2
	}

	return 0
}

// String ...
func (s Sex) String() string {
	switch s {
	case 1:
		return "М"
	case 2:
		return "Ж"
	}

	return ""
}

// User ...
type User struct {
	ID         int32
	Name       string
	TelegramID int64
	Email      maybe.Maybe[string]
	Phone      maybe.Maybe[string]
	State      int32
	Birthdate  maybe.Maybe[string]
	Sex        maybe.Maybe[Sex]
}

// MarshalJSON ...
func (u User) MarshalJSON() ([]byte, error) {
	type wrapperUser struct {
		ID         int32
		Name       string
		TelegramID int64
		Email      maybejson.Maybe[string]
		Phone      maybejson.Maybe[string]
		State      int32
		Birthdate  maybejson.Maybe[string]
		Sex        maybejson.Maybe[Sex]
	}

	wu := wrapperUser{
		ID:         u.ID,
		Name:       u.Name,
		TelegramID: u.TelegramID,
		Email:      maybejson.Wrap(u.Email),
		Phone:      maybejson.Wrap(u.Phone),
		State:      u.State,
		Birthdate:  maybejson.Wrap(u.Birthdate),
		Sex:        maybejson.Wrap(u.Sex),
	}
	return json.Marshal(wu)
}

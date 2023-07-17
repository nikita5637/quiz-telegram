package model

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
	ID        int32
	Email     string
	Name      string
	Phone     string
	State     int32
	Birthdate string
	Sex       Sex
}

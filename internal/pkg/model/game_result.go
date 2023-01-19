package model

// ResultPlace ...
type ResultPlace uint32

const (
	// InvalidPlace ...
	InvalidPlace ResultPlace = iota
	// FirstPlace ...
	FirstPlace
	// SecondPlace ...
	SecondPlace
	// ThrirdPlace ...
	ThrirdPlace
)

// String ...
func (p ResultPlace) String() string {
	switch p {
	case FirstPlace:
		return "🥇"
	case SecondPlace:
		return "🥈"
	case ThrirdPlace:
		return "🥉"
	}

	return ""
}

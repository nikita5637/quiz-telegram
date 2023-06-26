package model

// TODO use from proto
const (
	// LeagueQuizPlease ...
	LeagueQuizPlease = 1
	// LeagueSquiz ...
	LeagueSquiz = 2
)

// League ...
type League struct {
	ID        int32
	Name      string
	ShortName string
	LogoLink  string
	WebSite   string
}

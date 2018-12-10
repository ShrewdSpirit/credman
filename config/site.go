package config

type Site struct {
	Name              string
	Password          string
	Email             string
	Username          string
	Notes             string
	SecurityQuestions [5]string
}

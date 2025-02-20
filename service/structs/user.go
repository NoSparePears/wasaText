package structs

import "regexp"

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	UserPropic64 string `json:"photo"`
}

// validates username according to set rules
func (u *User) IsValid() bool {
	username := u.Username
	validUser := regexp.MustCompile(`[a-z][a-z0-9]{2,13}$`)
	return validUser.MatchString(username)
}

package serializers

import (
	"fmt"
)

type TwitterUser struct {
	Name         string
	Username     string
	AccessToken  string
	AccessSecret string
}

func (u *TwitterUser) GetDisplayName() string {
	return fmt.Sprintf("%s (@%s)", u.Name, u.Username)
}

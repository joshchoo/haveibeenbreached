package haveibeenbreached

import (
	"fmt"
	"regexp"
	"strings"
)

var accountEntityType = "Account"

type AccountItem struct {
	PK       string
	SK       string
	Type     string
	Username string
	Breaches []string
}

func (a AccountItem) isDBItem() bool {
	return true
}

type Account struct {
	Username Username
	Breaches []string
}

type Username interface {
	String() string
	PartitionKey() string
	SortKey() string
}

func (a Account) Item() AccountItem {
	return AccountItem{
		PK:       a.Username.PartitionKey(),
		SK:       a.Username.SortKey(),
		Type:     accountEntityType,
		Username: a.Username.String(),
		Breaches: a.Breaches,
	}
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Email struct {
	Domain string
	Alias  string
}

func NewEmail(emailStr string) (Email, error) {
	if !emailRegex.MatchString(emailStr) {
		return Email{}, fmt.Errorf("not a valid email address: %s", emailStr)
	}
	email := strings.Split(emailStr, "@")
	return Email{
		Alias:  email[0],
		Domain: email[1],
	}, nil
}

func (e Email) String() string {
	return fmt.Sprintf("%s@%s", e.Alias, e.Domain)
}

func (e Email) PartitionKey() string {
	return fmt.Sprintf("EMAIL#%s", e.Domain)
}

func (e Email) SortKey() string {
	return fmt.Sprintf("EMAIL#%s", e.Alias)
}

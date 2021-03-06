package haveibeenbreached

import (
	"fmt"
	"strings"
)

var accountEntityType = "Account"

type Username interface {
	String() string
	PartitionKey() string
	SortKey() string
}

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

func (a AccountItem) GetUsername() (Username, error) {
	username, err := NewEmailAccount(a.Username)
	if err != nil {
		return nil, err
	}
	return username, nil
}

func (a AccountItem) ToAccount() (Account, error) {
	username, err := a.GetUsername()
	if err != nil {
		return Account{}, err
	}
	return Account{
		Username: username,
		Breaches: a.Breaches,
	}, nil
}

type Account struct {
	Username Username
	Breaches []string
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

func (a Account) ToItem() DBItem {
	return a.Item()
}

type EmailAccount struct {
	Domain string
	Alias  string
}

func NewEmailAccount(emailStr string) (EmailAccount, error) {
	if !IsValidEmail(emailStr) {
		return EmailAccount{}, fmt.Errorf("not a valid email address: %s", emailStr)
	}
	email := strings.Split(emailStr, "@")
	return EmailAccount{
		Alias:  email[0],
		Domain: email[1],
	}, nil
}

func (e EmailAccount) String() string {
	return fmt.Sprintf("%s@%s", e.Alias, e.Domain)
}

func (e EmailAccount) PartitionKey() string {
	return fmt.Sprintf("EMAIL#%s", e.Domain)
}

func (e EmailAccount) SortKey() string {
	return fmt.Sprintf("EMAIL#%s", e.Alias)
}

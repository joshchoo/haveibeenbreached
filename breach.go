package haveibeenbreached

import (
	"fmt"
	"time"
)

var breachEntityType = "Breach"

/// BreachItem represents the schema stored in the database.
type BreachItem struct {
	PK               string
	SK               string
	Type             string
	BreachName       string
	Title            string
	Domain           string
	Description      string
	BreachDate       time.Time
	BreachedAccounts []string
}

func (b BreachItem) isDBItem() bool {
	return true
}

func (b BreachItem) ToBreach() Breach {
	return Breach{
		BreachName:       b.BreachName,
		Title:            b.Title,
		Domain:           b.Domain,
		Description:      b.Description,
		BreachDate:       b.BreachDate,
		BreachedAccounts: b.BreachedAccounts,
	}
}

/// Breach represents the breach domain model.
type Breach struct {
	BreachName       string
	Title            string
	Domain           string
	Description      string
	BreachDate       time.Time
	BreachedAccounts []string
}

func (b Breach) PartitionKey() string {
	return BreachPartitionKey(b.BreachName)
}

func (b Breach) SortKey() string {
	return BreachSortKey(b.BreachName)
}

func (b Breach) Item() BreachItem {
	return BreachItem{
		PK:               b.PartitionKey(),
		SK:               b.SortKey(),
		Type:             breachEntityType,
		BreachName:       b.BreachName,
		Title:            b.Title,
		Domain:           b.Domain,
		Description:      b.Description,
		BreachDate:       b.BreachDate,
		BreachedAccounts: b.BreachedAccounts,
	}
}

func (b Breach) ToItem() DBItem {
	return b.Item()
}

func (b Breach) AddAccounts(accounts []string) Breach {
	b.BreachedAccounts = accounts
	return b
}

func BreachPartitionKey(breachName string) string {
	return fmt.Sprintf("BREACH#%s", breachName)
}

func BreachSortKey(breachName string) string {
	return fmt.Sprintf("BREACH#%s", breachName)
}

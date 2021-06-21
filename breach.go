package haveibeenbreached

import (
	"fmt"
	"time"
)

var breachEntityType = "Breach"

/// BreachItem represents the schema stored in the database.
type BreachItem struct {
	PK          string
	SK          string
	Type        string
	BreachName  string
	Title       string
	Domain      string
	Description string
	BreachDate  time.Time
}

func (b BreachItem) isDBItem() bool {
	return true
}

/// Breach represents the breach domain model.
type Breach struct {
	BreachName  string
	Title       string
	Domain      string
	Description string
	BreachDate  time.Time
}

func (b Breach) PartitionKey() string {
	return fmt.Sprintf("BREACH#%s", b.BreachName)
}

func (b Breach) SortKey() string {
	return fmt.Sprintf("BREACH#%s", b.BreachName)
}

func (b Breach) Item() BreachItem {
	return BreachItem{
		PK:          b.PartitionKey(),
		SK:          b.SortKey(),
		Type:        breachEntityType,
		BreachName:  b.BreachName,
		Title:       b.Title,
		Domain:      b.Domain,
		Description: b.Description,
		BreachDate:  b.BreachDate,
	}
}

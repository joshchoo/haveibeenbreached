package haveibeenbreached

import (
	"fmt"
)

var subscriberEntityType = "Subscriber"

type SubscriberItem struct {
	PK    string
	SK    string
	Type  string
	Email string
}

func (s SubscriberItem) isDBItem() bool {
	return true
}

func (s SubscriberItem) ToSubscriber() Subscriber {
	return Subscriber{
		Email: s.Email,
	}
}

type Subscriber struct {
	Email string
}

func NewSubscriber(email string) (Subscriber, error) {
	if !IsValidEmail(email) {
		return Subscriber{}, fmt.Errorf("%s is not a valid email", email)
	}
	return Subscriber{Email: email}, nil
}

func (s Subscriber) PartitionKey() string {
	return fmt.Sprintf("SUB#%s", s.Email)
}

func (s Subscriber) SortKey() string {
	return fmt.Sprintf("SUB#%s", s.Email)
}

func (s Subscriber) Item() SubscriberItem {
	return SubscriberItem{
		PK:    s.PartitionKey(),
		SK:    s.SortKey(),
		Type:  subscriberEntityType,
		Email: s.Email,
	}
}

func (s Subscriber) ToItem() DBItem {
	return s.Item()
}

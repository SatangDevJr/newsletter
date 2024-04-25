package entity

import "time"

type Subscribers struct {
	ID               int64      `json:"id" sql:"id"`
	Email            string     `json:"email" sql:"email"`
	Name             string     `json:"name" sql:"name"`
	IsSubscribed     bool       `json:"isSubscribed" sql:"isSubscribed"`
	SubscribedDate   *time.Time `json:"subscribedDate" sql:"subscribedDate"`
	UnsubscribedDate *time.Time `json:"unsubscribedDate" sql:"unsubscribedDate"`
	DelFlag          *bool      `json:"delFlag" sql:"delFlag"`
}

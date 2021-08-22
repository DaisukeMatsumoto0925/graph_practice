package resolver

import (
	"github.com/DaisukeMatsumoto0925/backend/src/graphql/subscriber"
	"github.com/jinzhu/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Subscribers struct {
	Message *subscriber.MessageSubscriber
}
type Resolver struct {
	db          *gorm.DB
	subscribers Subscribers
}

func New(
	db *gorm.DB,
	subscribers Subscribers,
) *Resolver {
	return &Resolver{
		db:          db,
		subscribers: subscribers,
	}
}

package resolver

import "github.com/jinzhu/gorm"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	db *gorm.DB
}

func New(
	db *gorm.DB,
) *Resolver {
	return &Resolver{
		db: db,
	}
}

package controller

import (
	"github.com/DaisukeMatsumoto0925/backend/src/infra/redis"
	"github.com/jinzhu/gorm"
)

type Controller struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func NewController(db *gorm.DB, redisClient *redis.Client) *Controller {
	return &Controller{
		db:          db,
		redisClient: redisClient,
	}
}

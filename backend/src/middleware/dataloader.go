package middleware

import (
	"github.com/DaisukeMatsumoto0925/backend/src/dataloader"
	"github.com/labstack/echo"

	"github.com/jinzhu/gorm"
)

type Dataloader struct {
	DB *gorm.DB
}

func NewDataloader(db *gorm.DB) *Dataloader {
	return &Dataloader{DB: db}
}

func (loader Dataloader) InjectStoreStatusLoader() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			loader := dataloader.CreateUserLoader(loader.DB)
			newCtx := dataloader.SetUserLoader(ctx.Request().Context(), loader)

			ctx.SetRequest(ctx.Request().WithContext(newCtx))
			return h(ctx)
		}
	}
}

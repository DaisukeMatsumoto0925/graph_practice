package dataloader

import (
	"app/graph/generated"
	"app/graph/model"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
)

const userLoadersKey = "userLoader"

func UserLoaderMiddleware(db *gorm.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		loader := generated.NewUserLoader(generated.UserLoaderConfig{
			MaxBatch: 100,
			Wait:     1 * time.Millisecond,
			Fetch: func(ids []int) ([]*model.User, []error) {
				var users []*model.User
				errors := make([]error, len(ids))
				err := db.Debug().Where("id IN (?)", ids).Find(&users).Error
				if err != nil {
					for i := 0; i < len(ids); i++ {
						errors[i] = err
					}
				}
				return users, nil
			},
		},
		)
		ctx := context.WithValue(r.Context(), userLoadersKey, loader)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func User(ctx context.Context, id int) (*model.User, error) {
	v := ctx.Value(userLoadersKey)
	loader, ok := v.(*generated.UserLoader)

	if !ok {
		return nil, errors.New("failed to get loader from current context")
	}
	return loader.Load(id)
}

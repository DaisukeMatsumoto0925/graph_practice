package dataloader

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/DaisukeMatsumoto0925/backend2/graph/generated"
	gmodel "github.com/DaisukeMatsumoto0925/backend2/graph/model"
)

const userLoadersKey = "userLoader"

func CreateUserLoader(db *gorm.DB) *generated.UserLoader {
	return generated.NewUserLoader(generated.UserLoaderConfig{
		MaxBatch: 100,
		Wait:     1 * time.Millisecond,
		Fetch: func(ids []int) ([]*gmodel.User, []error) {
			var users []*gmodel.User
			errors := make([]error, len(ids))
			err := db.Where("id IN (?)", ids).Find(&users).Error
			if err != nil {
				for i := 0; i < len(ids); i++ {
					errors[i] = err
				}
			}

			userIDs := map[int]*gmodel.User{}
			for _, user := range users {
				idInt, _ := strconv.Atoi(user.ID)
				userIDs[idInt] = user
			}

			results := make([]*gmodel.User, len(ids))
			for i, id := range ids {
				results[i] = userIDs[id]
			}
			return results, nil
		},
	})
}

func User(ctx context.Context, id int) (*gmodel.User, error) {
	v := ctx.Value(userLoadersKey)
	if v == nil {
		panic("not found operator loader, must inject")
	}
	loader, ok := v.(*generated.UserLoader)
	if !ok {
		return nil, errors.New("failed to get loader from current context")
	}
	return loader.Load(id)
}

func SetUserLoader(ctx context.Context, userLoader *generated.UserLoader) context.Context {
	return context.WithValue(ctx, userLoadersKey, userLoader)
}

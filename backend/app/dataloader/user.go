package dataloader

import (
	"app/graph/generated"
	"app/graph/model"
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

const userLoadersKey = "userLoader"

func CreateUserLoader(db *gorm.DB) *generated.UserLoader {
	return generated.NewUserLoader(generated.UserLoaderConfig{
		MaxBatch: 100,
		Wait:     1 * time.Millisecond,
		Fetch: func(ids []int) ([]*model.User, []error) {
			var users []*model.User
			errors := make([]error, len(ids))
			err := db.Where("id IN (?)", ids).Find(&users).Error
			if err != nil {
				for i := 0; i < len(ids); i++ {
					errors[i] = err
				}
			}

			userID := map[int]*model.User{}
			for _, user := range users {
				idInt, _ := strconv.Atoi(user.ID)
				userID[idInt] = user
			}

			results := make([]*model.User, len(ids))
			for i, id := range ids {
				results[i] = userID[id]
			}
			return users, nil
		},
	})
}

func User(ctx context.Context, id int) (*model.User, error) {
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

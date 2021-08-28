package resolver

import (
	"context"
	"fmt"

	"github.com/DaisukeMatsumoto0925/backend/graph/generated"
	gmodel "github.com/DaisukeMatsumoto0925/backend/graph/model"
)

func (r *userResolver) ID(ctx context.Context, obj *gmodel.User) (string, error) {
	return fmt.Sprintf("%s:%s", "USER", obj.ID), nil
}

func (r *mutationResolver) JoinUser(ctx context.Context, input gmodel.JoinUserInput) (*gmodel.User, error) {
	var user gmodel.User
	if err := r.db.First(&user, input.UserID).Error; err != nil {
		return nil, err
	}

	if err := r.subscribers.Message.SetNx(ctx, user.ID); err != nil {
		return nil, err
	}

	return &user, nil
}

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }

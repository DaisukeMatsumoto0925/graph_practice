package resolver

import (
	"app/graph/generated"
	"app/graph/model"
	"context"
	"fmt"
)

func (r *userResolver) ID(ctx context.Context, obj *model.User) (string, error) {
	return fmt.Sprintf("%s:%s", "USER", obj.ID), nil
}

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }

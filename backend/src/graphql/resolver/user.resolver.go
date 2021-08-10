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

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }

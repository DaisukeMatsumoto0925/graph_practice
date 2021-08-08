package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/DaisukeMatsumoto0925/backend2/graph/generated"
	gmodel "github.com/DaisukeMatsumoto0925/backend2/graph/model"
)

func (r *taskResolver) ID(ctx context.Context, obj *gmodel.Task) (string, error) {
	return "id", nil
}

func (r *taskResolver) User(ctx context.Context, obj *gmodel.Task) (*gmodel.User, error) {
	panic(fmt.Errorf("not implemented"))
}

// Task returns generated.TaskResolver implementation.
func (r *Resolver) Task() generated.TaskResolver { return &taskResolver{r} }

type taskResolver struct{ *Resolver }

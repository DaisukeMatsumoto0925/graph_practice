package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/DaisukeMatsumoto0925/backend2/graph/generated"
	gmodel "github.com/DaisukeMatsumoto0925/backend2/graph/model"
)

func (r *mutationResolver) CreateTask(ctx context.Context, input gmodel.NewTask) (*gmodel.Task, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateTask(ctx context.Context, input gmodel.UpdateTask) (*gmodel.Task, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Tasks(ctx context.Context, input gmodel.PaginationInput) (*gmodel.TaskConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Task(ctx context.Context, id string) (*gmodel.Task, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

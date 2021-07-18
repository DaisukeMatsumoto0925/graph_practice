package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"app/graph/generated"
	"app/graph/model"
	"context"
	"errors"
	"fmt"
	"time"
)

func (r *mutationResolver) CreateTask(ctx context.Context, input model.NewTask) (*model.Task, error) {
	task := model.Task{
		Title:     input.Title,
		Note:      input.Note,
		Completed: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	r.DB.Create(&task)

	return &task, nil
}

func (r *mutationResolver) UpdateTask(ctx context.Context, input model.UpdateTask) (*model.Task, error) {
	var task model.Task
	r.DB.First(&task, input.ID)

	if input.Title == nil && input.Note == nil && input.Completed == nil {
		return nil, errors.New("could not update a task: params error")
	}

	if input.Title != nil {
		task.Title = *input.Title
	}
	if input.Note != nil {
		task.Note = *input.Note
	}
	if input.Completed != nil {
		task.Completed = *input.Completed
	}

	if err := r.DB.Save(task).Error; err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *queryResolver) Tasks(ctx context.Context, input model.PaginationInput) (*model.TaskConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Task(ctx context.Context, id int) (*model.Task, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

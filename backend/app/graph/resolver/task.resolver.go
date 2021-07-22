package resolver

import (
	"app/graph/model"
	"context"
	"errors"
	"strconv"
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
	// validation
	if input.First == nil && input.Last == nil {
		return nil, errors.New("input.First or input.Last is required: input error")
	}
	if input.First != nil && input.Last != nil {
		return nil, errors.New("input.First or input.Last is required: input error")
	}
	var limit int
	if input.First != nil {
		limit = *input.First
	} else {
		limit = *input.Last
	}

	var tasksSizeLimit = 30
	if input.First != nil && *input.First > tasksSizeLimit {
		return nil, errors.New("input.First exceeds tasksSizeLimit: input error ")
	}
	if input.Last != nil && *input.Last > tasksSizeLimit {
		return nil, errors.New("input.Last exceeds tasksSizeLimit: input error ")
	}

	// DB処理
	var tasks []*model.Task
	// var task *model.Task

	var cursorID int
	db := r.DB
	if input.After != nil {
		cursorID = *input.First
		db.Where("id >= ?", cursorID)
	} else {
		cursorID = *input.Last
		db.Where("id <= ?", cursorID).Order("id desc")
	}

	if err := db.Limit(limit + 1).Find(&tasks).Error; err != nil {
		return nil, errors.New("error")
	}

	startCursor := strconv.Itoa(tasks[0].ID)
	endCursor := strconv.Itoa(tasks[len(tasks)-1].ID)

	if len(tasks) > limit {
		endCursor := strconv.Itoa(tasks[limit-1].ID)
		return &model.TaskConnection{
			PageInfo: &model.PageInfo{StartCursor: &startCursor, EndCursor: &endCursor, HasNextPage: true, HasPreviousPage: true},
			Edges:    []*model.TaskEdge{},
			Nodes:    tasks,
		}, nil
	}

	return &model.TaskConnection{
		PageInfo: &model.PageInfo{
			StartCursor:     &startCursor,
			EndCursor:       &endCursor,
			HasNextPage:     len(tasks) <= limit,
			HasPreviousPage: false,
		},
		Edges: []*model.TaskEdge{},
		Nodes: tasks,
	}, nil
}

func (r *queryResolver) Task(ctx context.Context, id int) (*model.Task, error) {
	var task model.Task
	if err := r.DB.First(&task, id).Error; err != nil {
		return nil, err
	}

	return &task, nil
}

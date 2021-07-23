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
		return nil, errors.New("passing input.First and input.Last is not supported: input error")
	}
	if input.Before != nil && input.After != nil {
		return nil, errors.New("passing input.Before and input.After is not supported: input error")
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

	// connection情報作成
	var cursorID int

	db := r.DB
	if input.After != nil {
		cursorID, _ = strconv.Atoi(*input.After)
		db = db.Where("id > ?", cursorID)
	}

	if input.Before != nil {
		cursorID, _ = strconv.Atoi(*input.Before)
		db = db.Where("id < ?", cursorID).Order("id desc")
	}

	// SELECT
	var tasks []*model.Task
	if err := db.Limit(limit + 1).Find(&tasks).Error; err != nil {
		return nil, errors.New("could not find tasks: data base error ")
	}

	if len(tasks) == 0 {
		return &model.TaskConnection{
			PageInfo: &model.PageInfo{
				StartCursor:     nil,
				EndCursor:       nil,
				HasNextPage:     false,
				HasPreviousPage: false,
			},
			Edges: []*model.TaskEdge{},
			Nodes: []*model.Task{},
		}, nil
	}

	// TODO: hasPreviousPageの判定
	edges := make([]*model.TaskEdge, len(tasks))

	for i, task := range tasks {
		newEdge := &model.TaskEdge{
			Cursor: strconv.Itoa(task.ID),
			Node:   task,
		}
		edges[i] = newEdge
	}

	startCursor := edges[0].Cursor
	endCursor := edges[len(edges)-2].Cursor

	// limitより件数が多いのでHasNextPageがtrueになる
	if len(tasks) > limit {
		return &model.TaskConnection{
			PageInfo: &model.PageInfo{
				StartCursor:     &startCursor,
				EndCursor:       &endCursor,
				HasNextPage:     true,
				HasPreviousPage: false,
			},
			Edges: edges[:len(edges)-1],
			Nodes: tasks[:len(tasks)-1],
		}, nil
	}

	return &model.TaskConnection{
		PageInfo: &model.PageInfo{
			StartCursor:     &startCursor,
			EndCursor:       &endCursor,
			HasNextPage:     false,
			HasPreviousPage: false,
		},
		Edges: edges,
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

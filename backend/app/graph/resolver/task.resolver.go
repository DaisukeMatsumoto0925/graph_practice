package resolver

import (
	"app/dataloader"
	"app/domain"
	"app/graph/generated"
	"app/graph/model"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"
)

func (r *mutationResolver) CreateTask(ctx context.Context, input model.NewTask) (*model.Task, error) {
	var user model.User
	if err := r.DB.Where("name = ?", "ADMIN").First(&user).Error; err != nil {
		return nil, err
	}

	userID, err := strconv.Atoi(user.ID)
	if err != nil {
		return nil, err
	}

	task := domain.Task{
		ID:        0,
		UserID:    userID,
		Title:     input.Title,
		Note:      input.Note,
		Completed: 0,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	if err := r.DB.Create(&task).Error; err != nil {
		return nil, err
	}

	graphTask := model.Task{
		ID:        strconv.Itoa(task.ID),
		UserID:    "USER:" + strconv.Itoa(userID),
		Title:     task.Title,
		Note:      task.Note,
		Completed: task.Completed,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}

	return &graphTask, nil
}

func (r *mutationResolver) UpdateTask(ctx context.Context, input model.UpdateTask) (*model.Task, error) {
	var task model.Task
	if err := r.DB.First(&task, input.ID).Error; err != nil {
		return nil, err
	}

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

	if err := r.DB.Save(&task).Error; err != nil {
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
	var tasksSizeLimit = 100
	if input.First != nil && *input.First > tasksSizeLimit {
		return nil, errors.New("input.First exceeds tasksSizeLimit: input error ")
	}
	if input.Last != nil && *input.Last > tasksSizeLimit {
		return nil, errors.New("input.Last exceeds tasksSizeLimit: input error ")
	}

	db := r.DB
	var tasks []*model.Task

	if input.After != nil {
		db = db.Where("id > ?", *input.After)
	}

	if input.Before != nil {
		db = db.Where("id < ?", *input.Before).Order("id desc")
	}

	if input.Last != nil {
		db = db.Order("id desc")
	}

	if err := db.Limit(limit + 1).Find(&tasks).Error; err != nil {
		return nil, errors.New("could not find tasks: data base error ")
	}

	//検索結果0の場合
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

	edges := make([]*model.TaskEdge, len(tasks))
	nodes := make([]*model.Task, len(tasks))

	// last, before 指定の時はスライスの後ろから入れていく
	if input.First != nil || input.After != nil {
		for i, task := range tasks {
			newEdge := &model.TaskEdge{
				Cursor: task.ID,
				Node:   task,
			}
			nodes = tasks
			edges[i] = newEdge
		}
	} else {
		for i, task := range tasks {
			newEdge := &model.TaskEdge{
				Cursor: task.ID,
				Node:   task,
			}
			nodes[len(tasks)-1-i] = tasks[i]
			edges[len(edges)-1-i] = newEdge
		}
	}

	startCursor := edges[0].Cursor
	endCursor := edges[len(edges)-1].Cursor

	var hasPreviousPage bool
	var hasNextPage bool

	// startCursorのIDより前に1件でもデータがある場合はpreviousPageはtrue
	startCursorInt, _ := strconv.Atoi(startCursor)
	endCursorInt, _ := strconv.Atoi(endCursor)
	var task model.Task
	if input.First != nil {
		if err := r.DB.Where("id <= ?", startCursorInt-1).First(&task).Error; err == nil {
			hasPreviousPage = true
		}
	} else {
		if err := r.DB.Where("id >= ?", endCursorInt+1).First(&task).Error; err == nil {
			hasNextPage = true
		}
	}

	// Firstが渡された場合 if limit以上 else limit以下
	if input.First != nil && limit < len(nodes) {
		endCursor = edges[len(edges)-2].Cursor
		hasNextPage = true
		return &model.TaskConnection{
			PageInfo: &model.PageInfo{
				StartCursor:     &startCursor,
				EndCursor:       &endCursor,
				HasNextPage:     hasNextPage,
				HasPreviousPage: hasPreviousPage,
			},
			Edges: edges[:len(edges)-1],
			Nodes: nodes[:len(nodes)-1],
		}, nil
	} else if input.First != nil && limit >= len(nodes) {
		return &model.TaskConnection{
			PageInfo: &model.PageInfo{
				StartCursor:     &startCursor,
				EndCursor:       &endCursor,
				HasNextPage:     hasNextPage,
				HasPreviousPage: hasPreviousPage,
			},
			Edges: edges,
			Nodes: nodes,
		}, nil
	}

	// Lastが渡された場合 if limit以上 else limit以下
	if input.Last != nil && limit < len(nodes) {
		startCursor = edges[len(edges)-limit].Cursor
		hasPreviousPage = true
		return &model.TaskConnection{
			PageInfo: &model.PageInfo{
				StartCursor:     &startCursor,
				EndCursor:       &endCursor,
				HasNextPage:     hasNextPage,
				HasPreviousPage: hasPreviousPage,
			},
			Edges: edges[len(edges)-limit:],
			Nodes: nodes[len(nodes)-limit:],
		}, nil
	} else if input.Last != nil && limit >= len(nodes) {
		return &model.TaskConnection{
			PageInfo: &model.PageInfo{
				StartCursor:     &startCursor,
				EndCursor:       &endCursor,
				HasNextPage:     hasNextPage,
				HasPreviousPage: hasPreviousPage,
			},
			Edges: edges,
			Nodes: nodes,
		}, nil
	}

	return nil, nil
}

func (r *queryResolver) Task(ctx context.Context, id string) (*model.Task, error) {
	var task model.Task
	if err := r.DB.First(&task, id).Error; err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *taskResolver) ID(ctx context.Context, obj *model.Task) (string, error) {
	return fmt.Sprintf("%s:%s", "TASK", obj.ID), nil
}

func (r *taskResolver) User(ctx context.Context, obj *model.Task) (*model.User, error) {
	// var user model.User
	// if err := r.DB.Debug().First(&user, obj.UserID).Error; err != nil {
	// 	return nil, err
	// }
	// return &user, nil
	idInt, _ := strconv.Atoi(obj.UserID)
	user, err := dataloader.User(ctx, idInt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Task returns generated.TaskResolver implementation.
func (r *Resolver) Task() generated.TaskResolver { return &taskResolver{r} }

type taskResolver struct{ *Resolver }

package resolver

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/DaisukeMatsumoto0925/backend/graph/generated"
	gmodel "github.com/DaisukeMatsumoto0925/backend/graph/model"
	"github.com/DaisukeMatsumoto0925/backend/src/dataloader"
	"github.com/DaisukeMatsumoto0925/backend/src/domain"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func (r *mutationResolver) CreateTask(ctx context.Context, input gmodel.NewTask) (*gmodel.Task, error) {
	var user gmodel.User
	if err := r.db.Where("name = ?", "ADMIN").First(&user).Error; err != nil {
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

	if err := r.db.Create(&task).Error; err != nil {
		return nil, err
	}

	graphTask := gmodel.Task{
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

func (r *mutationResolver) UpdateTask(ctx context.Context, input gmodel.UpdateTask) (*gmodel.Task, error) {
	var task gmodel.Task
	if err := r.db.First(&task, input.ID).Error; err != nil {
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

	if err := r.db.Save(&task).Error; err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *queryResolver) Tasks(ctx context.Context, input gmodel.PaginationInput) (*gmodel.TaskConnection, error) {
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

	db := r.db
	var tasks []*gmodel.Task

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
		return &gmodel.TaskConnection{
			PageInfo: &gmodel.PageInfo{
				StartCursor:     nil,
				EndCursor:       nil,
				HasNextPage:     false,
				HasPreviousPage: false,
			},
			Edges: []*gmodel.TaskEdge{},
			Nodes: []*gmodel.Task{},
		}, nil
	}

	edges := make([]*gmodel.TaskEdge, len(tasks))
	nodes := make([]*gmodel.Task, len(tasks))

	// last, before 指定の時はスライスの後ろから入れていく
	if input.First != nil || input.After != nil {
		for i, task := range tasks {
			newEdge := &gmodel.TaskEdge{
				Cursor: task.ID,
				Node:   task,
			}
			nodes = tasks
			edges[i] = newEdge
		}
	} else {
		for i, task := range tasks {
			newEdge := &gmodel.TaskEdge{
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
	var task gmodel.Task
	if input.First != nil {
		if err := r.db.Where("id <= ?", startCursorInt-1).First(&task).Error; err == nil {
			hasPreviousPage = true
		}
	} else {
		if err := r.db.Where("id >= ?", endCursorInt+1).First(&task).Error; err == nil {
			hasNextPage = true
		}
	}

	// Firstが渡された場合 if limit以上 else limit以下
	if input.First != nil && limit < len(nodes) {
		endCursor = edges[len(edges)-2].Cursor
		hasNextPage = true
		return &gmodel.TaskConnection{
			PageInfo: &gmodel.PageInfo{
				StartCursor:     &startCursor,
				EndCursor:       &endCursor,
				HasNextPage:     hasNextPage,
				HasPreviousPage: hasPreviousPage,
			},
			Edges: edges[:len(edges)-1],
			Nodes: nodes[:len(nodes)-1],
		}, nil
	} else if input.First != nil && limit >= len(nodes) {
		return &gmodel.TaskConnection{
			PageInfo: &gmodel.PageInfo{
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
		return &gmodel.TaskConnection{
			PageInfo: &gmodel.PageInfo{
				StartCursor:     &startCursor,
				EndCursor:       &endCursor,
				HasNextPage:     hasNextPage,
				HasPreviousPage: hasPreviousPage,
			},
			Edges: edges[len(edges)-limit:],
			Nodes: nodes[len(nodes)-limit:],
		}, nil
	} else if input.Last != nil && limit >= len(nodes) {
		return &gmodel.TaskConnection{
			PageInfo: &gmodel.PageInfo{
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

func (r *queryResolver) Task(ctx context.Context, id string) (*gmodel.Task, error) {
	var task gmodel.Task

	if err := r.db.First(&task, id).Error; err != nil {
		graphql.AddError(ctx, &gqlerror.Error{
			Path:    graphql.GetPath(ctx),
			Message: fmt.Sprintf("Error %s", err),
			Extensions: map[string]interface{}{
				"code": "code1",
			},
		})
		return nil, nil
	}

	return &task, nil
}

func (r *taskResolver) ID(ctx context.Context, obj *gmodel.Task) (string, error) {
	return fmt.Sprintf("%s:%s", "TASK", obj.ID), errors.New("err")
}

func (r *taskResolver) User(ctx context.Context, obj *gmodel.Task) (*gmodel.User, error) {
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

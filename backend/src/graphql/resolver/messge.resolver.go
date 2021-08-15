package resolver

import (
	"context"

	gmodel "github.com/DaisukeMatsumoto0925/backend/graph/model"
)

func (r *mutationResolver) PostMessage(ctx context.Context, input gmodel.PostMessageInput) (*gmodel.Message, error) {
	return nil, nil
}

package resolver

import (
	"context"

	gmodel "github.com/DaisukeMatsumoto0925/backend/graph/model"
)

func (r *subscriptionResolver) MessagePosted(ctx context.Context, userID *string) (<-chan *gmodel.Message, error) {
	return nil, nil
}
func (r *subscriptionResolver) UserJoined(ctx context.Context, userID *string) (<-chan *gmodel.User, error) {
	return nil, nil
}

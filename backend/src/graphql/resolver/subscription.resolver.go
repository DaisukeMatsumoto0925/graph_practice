package resolver

import (
	"context"

	gmodel "github.com/DaisukeMatsumoto0925/backend/graph/model"
)

func (r *subscriptionResolver) MessagePosted(ctx context.Context, userID *string) (<-chan *gmodel.Message, error) {
	msgSubscriber := r.subscribers.Message
	isJoined, err := msgSubscriber.CheckJoined(ctx, *userID)
	if err != nil || !isJoined {
		return nil, err
	}

	msgChan := msgSubscriber.MakeChan(ctx, userID)
	return msgChan, nil
}
func (r *subscriptionResolver) UserJoined(ctx context.Context, userID *string) (<-chan *gmodel.User, error) {
	return nil, nil
}

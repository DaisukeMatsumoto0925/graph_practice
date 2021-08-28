package resolver

import (
	"context"

	gmodel "github.com/DaisukeMatsumoto0925/backend/graph/model"
)

func (r *mutationResolver) PostMessage(ctx context.Context, input gmodel.PostMessageInput) (*gmodel.Message, error) {
	msgSubscriber := r.subscribers.Message
	isJoined, err := msgSubscriber.CheckJoined(ctx, input.UserID)
	if err != nil || !isJoined {
		return nil, err
	}

	isSet, err := msgSubscriber.SetExpire(ctx, input.UserID)
	if err != nil || !isSet {
		return nil, err
	}

	var user gmodel.User
	if err := r.db.First(&user, input.UserID).Error; err != nil {
		return nil, err
	}

	m := &gmodel.Message{
		ID:      user.ID,
		User:    &user,
		Message: input.Message,
	}

	if err := msgSubscriber.PublishMsg(ctx, m); err != nil {
		return nil, err
	}

	return m, nil
}

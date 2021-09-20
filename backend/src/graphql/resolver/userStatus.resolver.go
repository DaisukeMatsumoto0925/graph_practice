package resolver

import (
	"context"
	"time"

	gmodel "github.com/DaisukeMatsumoto0925/backend/graph/model"
)

func (r *mutationResolver) UpdateUserStatus(ctx context.Context, input gmodel.UpdateUserStatusInput) (*gmodel.UserStatus, error) {
	userStatusSubs := r.subscribers.UserStatus

	if err := userStatusSubs.Client.Set(
		ctx, input.UserID,
		string(gmodel.StatusOnline),
		time.Millisecond*time.Duration(6000),
	).Err(); err != nil {
		return nil, err
	}

	return &gmodel.UserStatus{
		UserID: input.UserID,
		Status: gmodel.StatusOnline,
	}, nil
}

func (r *subscriptionResolver) UserStatusChanged(ctx context.Context, userID string) (<-chan *gmodel.UserStatus, error) {
	userStatusSubs := r.subscribers.UserStatus

	userStatusSubs.Mutex.Lock()
	channels, ok := userStatusSubs.UserStatusChannels[userID]
	if !ok {
		channels = make(chan *gmodel.UserStatus)
		userStatusSubs.UserStatusChannels[userID] = channels
	}
	userStatusSubs.Mutex.Unlock()

	go func() {
		<-ctx.Done()
		userStatusSubs.Mutex.Lock()
		delete(userStatusSubs.UserStatusChannels, userID)
		userStatusSubs.Mutex.Unlock()
	}()

	return channels, nil
}

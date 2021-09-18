package resolver

import (
	"context"

	gmodel "github.com/DaisukeMatsumoto0925/backend/graph/model"
)

func (r *mutationResolver) UpdateUserStatus(ctx context.Context, input gmodel.UpdateUserStatusInput) (*gmodel.UserStatus, error) {
	panic("")
}

func (r *subscriptionResolver) UserStatusChanged(ctx context.Context, userID string) (<-chan *gmodel.UserStatus, error) {
	panic("")
}

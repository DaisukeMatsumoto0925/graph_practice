package subscriber

import (
	"context"
	"fmt"
	"strings"
	"sync"

	gmodel "github.com/DaisukeMatsumoto0925/backend/graph/model"
	"github.com/go-redis/redis/v8"
)

type UserStatusSubscriber struct {
	Client             *redis.Client
	UserStatusChannels map[string]chan *gmodel.UserStatus
	Mutex              sync.Mutex
}

func NewUserStatusSubscriber(ctx context.Context, client *redis.Client) *UserStatusSubscriber {
	subscriber := &UserStatusSubscriber{
		Client:             client,
		UserStatusChannels: map[string]chan *gmodel.UserStatus{},
		Mutex:              sync.Mutex{},
	}
	subscriber.startSubscribingRedis(ctx)
	return subscriber
}

func (m *UserStatusSubscriber) startSubscribingRedis(ctx context.Context) error {
	go func() {
		pubsub := m.Client.PSubscribe(ctx, "__keyspace@0__:*")
		defer pubsub.Close()
		ch := pubsub.Channel()

		for {
			select {
			case <-ctx.Done():
			case msg := <-ch:
				switch msg.Payload {
				case "set":
					prefix := "__keyspace@0__:"
					userID := strings.TrimPrefix(msg.Channel, prefix)
					status, err := m.Client.Get(ctx, userID).Result()
					if err != nil {
						fmt.Println("Redis Error GET:", err)
						continue
					}

					userStatus := &gmodel.UserStatus{
						UserID: userID,
					}
					switch status {
					case "ONLINE":
						userStatus.Status = gmodel.StatusOnline
					}

					m.Mutex.Lock()
					for _, ch := range m.UserStatusChannels {
						ch <- userStatus
					}
					m.Mutex.Unlock()
				case "expired":
					prefix := "__keyspace@0__:"
					userID := strings.TrimPrefix(msg.Channel, prefix)
					userStatus := &gmodel.UserStatus{
						UserID: userID,
						Status: gmodel.StatusOffline,
					}

					m.Mutex.Lock()
					for _, ch := range m.UserStatusChannels {
						ch <- userStatus
					}
					m.Mutex.Unlock()
				// case "expire":
				// case "del":
				default:
				}
			}
		}
	}()

	return nil
}

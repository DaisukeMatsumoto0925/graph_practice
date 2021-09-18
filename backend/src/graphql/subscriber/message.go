package subscriber

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	gmodel "github.com/DaisukeMatsumoto0925/backend/graph/model"
	"github.com/go-redis/redis/v8"
)

type MessageSubscriber struct {
	client       *redis.Client
	msgChannels  map[string]chan *gmodel.Message
	userChannels map[string]chan *gmodel.User
	mutex        sync.Mutex
}

func NewMessageSubscriber(ctx context.Context, client *redis.Client) *MessageSubscriber {
	subscriber := &MessageSubscriber{
		client:       client,
		msgChannels:  map[string]chan *gmodel.Message{},
		userChannels: map[string]chan *gmodel.User{},
		mutex:        sync.Mutex{},
	}
	subscriber.startSubscribingRedis(ctx)
	return subscriber
}

func (m *MessageSubscriber) CheckJoined(ctx context.Context, userID string) (bool, error) {
	val, err := m.client.Exists(ctx, userID).Result()
	if err != nil {
		log.Println(err)
		return false, err
	}

	if val == 1 {
		return true, nil
	}
	return false, nil
}

func (m *MessageSubscriber) SetExpire(ctx context.Context, userID string) (bool, error) {
	val, err := m.client.SetXX(ctx, userID, userID, 60*time.Minute).Result()
	if !val {
		return val, err
	}
	return true, nil

}

func (m *MessageSubscriber) PublishMsg(ctx context.Context, msg *gmodel.Message) error {
	mb, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	m.client.Publish(ctx, "room", mb)
	return nil
}

func (m *MessageSubscriber) MakeChan(ctx context.Context, userID *string) <-chan *gmodel.Message {
	messageChan := make(chan *gmodel.Message, 1)
	m.mutex.Lock()
	m.msgChannels[*userID] = messageChan
	m.mutex.Unlock()

	go func() {
		<-ctx.Done()
		m.mutex.Lock()
		delete(m.msgChannels, *userID)
		m.mutex.Unlock()
		m.client.Del(ctx, *userID)
	}()

	return messageChan
}

func (m *MessageSubscriber) SetNx(ctx context.Context, userID string) error {
	val, err := m.client.SetNX(ctx, userID, userID, 60*time.Minute).Result()
	if err != nil {
		log.Println(err)
		return err
	}
	if !val {
		return errors.New("this user name has already used")
	}

	m.mutex.Lock()
	for _, ch := range m.userChannels {
		ch <- &gmodel.User{ID: userID}
	}
	m.mutex.Unlock()
	return nil
}

func (m *MessageSubscriber) startSubscribingRedis(ctx context.Context) error {
	var err error
	go func() {
		pubsub := m.client.Subscribe(ctx, "room")
		defer pubsub.Close()

		for {
			msgi, err := pubsub.Receive(ctx)
			if err != nil {
				continue
			}
			switch msg := msgi.(type) {
			case *redis.Message:
				var ms gmodel.Message
				if err := json.Unmarshal([]byte(msg.Payload), &ms); err != nil {
					continue
				}

				m.mutex.Lock()
				for _, ch := range m.msgChannels {
					ch <- &ms
				}
				m.mutex.Unlock()
			default:
			}
		}
	}()

	if err != nil {
		return err
	}

	return nil
}

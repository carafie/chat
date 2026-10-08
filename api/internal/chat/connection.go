package chat

import (
	"slices"
	"sync"
	"time"
	"uuid"

	"github.com/carafie/chat/api/internal/event"
	"github.com/carafie/chat/api/pkg/errcode"
	"github.com/carafie/chat/api/pkg/pubsub"
)

var (
	ErrConnectionGuestUnset  = errcode.New("GUEST_UNSET", "connection guest is unset")
	ErrConnectionRoomUnknown = errcode.New("ROOM_UNKNOWN", "connection room is unknown")
)

const rawEventsBuffer = 16

type Connection struct {
	guest       Guest
	domain      Domain
	joinedRooms []joinedRoom
	pubsub      *pubsub.Client
	rawEvents   chan []byte
	close       sync.Once
}

type joinedRoom struct {
	room         Room
	subscription pubsub.Subscription
}

func NewConnection(domain Domain, pubsub *pubsub.Client) *Connection {
	return &Connection{
		domain:      domain,
		joinedRooms: make([]joinedRoom, 0),
		pubsub:      pubsub,
		rawEvents:   make(chan []byte, rawEventsBuffer),
	}
}

func (c *Connection) Inbound(raw []byte) {
	in, err := event.Deserialize[event.Inbound](raw)
	if err != nil {
		c.error(in.ID, err)
		return
	}

	switch in.Kind {
	case event.KindJoinRoom:
		err = c.joinRoom(in)
	case event.KindLeaveRoom:
		err = c.leaveRoom(in)
	case event.KindSetGuest:
		err = c.setGuest(in)
	case event.KindSendMessage:
		err = c.sendMessage(in)
	default:
		err = event.ErrInvalid
	}
	if err != nil {
		c.error(in.ID, err)
	}
}

func (c *Connection) joinRoom(in event.Inbound) error {
	data, err := event.Deserialize[event.InboundJoinRoom](in.Data)
	if err != nil {
		return err
	}
	room, err := NewRoom(data.RoomName)
	if err != nil {
		return err
	}

	for _, joinedRoom := range c.joinedRooms {
		if joinedRoom.room.Name == room.Name {
			return nil
		}
	}
	topic := c.domain.Name + "/" + room.Name
	subscription := c.pubsub.Subscribe(topic, c.rawEvents)
	c.joinedRooms = append(c.joinedRooms, joinedRoom{
		room:         room,
		subscription: subscription,
	})
	return nil
}

func (c *Connection) leaveRoom(in event.Inbound) error {
	data, err := event.Deserialize[event.InboundLeaveRoom](in.Data)
	if err != nil {
		return err
	}
	room, err := NewRoom(data.RoomName)
	if err != nil {
		return err
	}

	for i, joinedRoom := range c.joinedRooms {
		if joinedRoom.room.Name == room.Name {
			c.pubsub.Unsubscribe(joinedRoom.subscription)
			c.joinedRooms = slices.Delete(c.joinedRooms, i, i+1)
			return nil
		}
	}
	return nil
}

func (c *Connection) setGuest(in event.Inbound) error {
	data, err := event.Deserialize[event.InboundSetGuest](in.Data)
	if err != nil {
		return err
	}
	guest, err := NewGuest(uuid.NewV7(), data.GuestName, data.GuestColor)
	if err != nil {
		return err
	}

	c.guest = guest
	return nil
}

func (c *Connection) sendMessage(in event.Inbound) error {
	if c.guest == (Guest{}) {
		return ErrConnectionGuestUnset
	}

	data, err := event.Deserialize[event.InboundSendMessage](in.Data)
	if err != nil {
		return err
	}
	room, err := NewRoom(data.RoomName)
	if err != nil {
		return err
	}
	message, err := NewMessage(data.Message)
	if err != nil {
		return err
	}

	for _, joinedRoom := range c.joinedRooms {
		if joinedRoom.room.Name == room.Name {
			return c.message(joinedRoom, message)
		}
	}
	return ErrConnectionRoomUnknown
}

func (c *Connection) message(joinedRoom joinedRoom, message Message) error {
	out := event.Outbound{
		Kind: event.KindMessage,
		Data: event.OutboundMessage{
			ID:         uuid.NewV7(),
			RoomName:   joinedRoom.room.Name,
			GuestID:    c.guest.ID,
			GuestName:  c.guest.Name,
			GuestColor: c.guest.Color,
			Message:    message,
			SentAt:     time.Now().UTC(),
		},
	}
	raw, err := event.Serialize(out)
	if err != nil {
		return err
	}
	c.pubsub.Publish(joinedRoom.subscription.Topic, raw)
	return nil
}

func (c *Connection) error(id string, err error) {
	out := event.Outbound{
		Kind: event.KindError,
		Data: event.OutboundError{ID: id, Code: errcode.Code(err)},
	}
	if raw, err := event.Serialize(out); err == nil {
		c.rawEvents <- raw
	}
}

func (c *Connection) Outbound() (raw <-chan []byte) {
	return c.rawEvents
}

func (c *Connection) Close() {
	c.close.Do(func() {
		for _, joinedRoom := range c.joinedRooms {
			c.pubsub.Unsubscribe(joinedRoom.subscription)
		}
		close(c.rawEvents)
	})
}

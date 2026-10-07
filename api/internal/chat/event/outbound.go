package event

import (
	"time"
	"uuid"
)

// Outbound event kind.
const (
	KindMessage Kind = "message"
	KindError   Kind = "error"
)

// Outbound is an envelope for outgoing events.
type Outbound struct {
	Kind Kind `json:"kind"`
	Data any  `json:"data"`
}

type OutboundMessage struct {
	ID         uuid.UUID `json:"id"`
	RoomName   string    `json:"room_name"`
	GuestID    uuid.UUID `json:"guest_id"`
	GuestName  string    `json:"guest_name"`
	GuestColor string    `json:"guest_color"`
	Message    string    `json:"message"`
	SentAt     time.Time `json:"sent_at"`
}

type OutboundError struct {
	ID   string `json:"id"`
	Code string `json:"code"`
}

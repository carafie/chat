package event

import "encoding/json/jsontext"

// Inbound event kind.
const (
	KindJoinRoom    Kind = "join_room"
	KindLeaveRoom   Kind = "leave_room"
	KindSetGuest    Kind = "set_guest"
	KindSendMessage Kind = "send_message"
)

// Inbound is an envelope for incoming events.
type Inbound struct {
	ID   string         `json:"id"`
	Kind Kind           `json:"kind"`
	Data jsontext.Value `json:"data"`
}

type InboundJoinRoom struct {
	RoomName string `json:"room_name"`
}

type InboundLeaveRoom struct {
	RoomName string `json:"room_name"`
}

type InboundSetGuest struct {
	GuestName  string `json:"guest_name"`
	GuestColor string `json:"guest_color"`
}

type InboundSendMessage struct {
	RoomName string `json:"room_name"`
	Message  string `json:"message"`
}

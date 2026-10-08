package web

import (
	"context"
	"net/http"
	"time"

	"github.com/carafie/chat/api/internal/chat"
	"github.com/carafie/chat/api/pkg/pubsub"
	"github.com/coder/websocket"
)

const (
	connectPattern = "GET /connect"
	wsWriteTimeout = 10 * time.Second
)

type Handler struct {
	pubsub *pubsub.Client
}

func NewHandler(pubsub *pubsub.Client) *Handler {
	return &Handler{pubsub: pubsub}
}

func (h *Handler) Connect(w http.ResponseWriter, r *http.Request) {
	domain, err := chat.NewDomain(r.Header.Get("Origin"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	wsConn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		return
	}
	defer wsConn.CloseNow()

	chatConn := chat.NewConnection(domain, h.pubsub)
	defer chatConn.Close()

	ctx, cancel := context.WithCancel(context.Background())

	outboundDone := make(chan struct{})
	go func() {
		defer close(outboundDone)
		defer cancel()

		for {
			select {
			case <-ctx.Done():
				return
			case raw, ok := <-chatConn.Outbound():
				if !ok {
					return
				}
				writeCtx, writeCancel := context.WithTimeout(ctx, wsWriteTimeout)
				err := wsConn.Write(writeCtx, websocket.MessageText, raw)
				writeCancel()
				if err != nil {
					return
				}
			}
		}
	}()

	for {
		_, raw, err := wsConn.Read(ctx)
		if err != nil {
			break
		}
		chatConn.Inbound(raw)
	}
	cancel()
	<-outboundDone
}

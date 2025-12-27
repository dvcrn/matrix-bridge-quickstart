package connector

import (
	"context"
	"fmt"
	"time"

	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/networkid"
	"maunium.net/go/mautrix/bridgev2/simplevent"
	"maunium.net/go/mautrix/event"
)

// File where you can put all the events from the upstream network
// For example when you receive a new message
// This file is responsible for bridging those upstream things to matrix
//
// Some things you can do here:
// - Connect to an upstream websocket
// - Poll for messages

// QueueRemoteMessage shows the preferred Remote -> Matrix flow using the bridge event queue.
func (nc *MyNetworkClient) QueueRemoteMessage(ctx context.Context, portalID networkid.PortalID, body string) {
	evt := &simplevent.Message[string]{
		EventMeta: simplevent.EventMeta{
			Type:      bridgev2.RemoteEventMessage,
			PortalKey: networkid.PortalKey{ID: portalID},
			Sender: bridgev2.EventSender{
				Sender:   networkid.UserID("example-ghost"),
				IsFromMe: false,
			},
			CreatePortal: true,
			Timestamp:    time.Now(),
		},
		Data: body,
		ID:   networkid.MessageID(fmt.Sprintf("remote-%s-%d", portalID, time.Now().UnixNano())),
		ConvertMessageFunc: func(ctx context.Context, portal *bridgev2.Portal, intent bridgev2.MatrixAPI, data string) (*bridgev2.ConvertedMessage, error) {
			return &bridgev2.ConvertedMessage{
				Parts: []*bridgev2.ConvertedMessagePart{{
					Type: event.EventMessage,
					Content: &event.MessageEventContent{
						MsgType: event.MsgText,
						Body:    data,
					},
				}},
			}, nil
		},
	}

	nc.bridge.QueueRemoteEvent(nc.login, evt)
}

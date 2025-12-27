package connector

import (
	"context"
	"time"

	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/networkid"
	"maunium.net/go/mautrix/event"
)

// BackfillingNetworkAPI is responsible for loading historic messages
var _ bridgev2.BackfillingNetworkAPI = (*MyNetworkClient)(nil)

// FetchMessages implements [bridgev2.BackfillingNetworkAPI].
// This wil get called when the user opens a room and wants to load historical messages
func (nc *MyNetworkClient) FetchMessages(ctx context.Context, fetchParams bridgev2.FetchMessagesParams) (*bridgev2.FetchMessagesResponse, error) {
	portal := fetchParams.Portal
	now := time.Now().UTC()

	backfillMsg := &bridgev2.BackfillMessage{
		ID:        networkid.MessageID("static-backfill-" + string(portal.ID)),
		Timestamp: now.Add(-1 * time.Minute),
		Sender: bridgev2.EventSender{
			Sender:   networkid.UserID("example-ghost"),
			IsFromMe: false,
		},
		ConvertedMessage: &bridgev2.ConvertedMessage{
			Parts: []*bridgev2.ConvertedMessagePart{{
				Type: event.EventMessage,
				Content: &event.MessageEventContent{
					MsgType: event.MsgText,
					Body:    "Backfilled history (placeholder).",
				},
			}},
		},
	}

	return &bridgev2.FetchMessagesResponse{
		Messages:                []*bridgev2.BackfillMessage{backfillMsg},
		HasMore:                 false,
		Forward:                 fetchParams.Forward,
		MarkRead:                !fetchParams.Forward,
		AggressiveDeduplication: true,
	}, nil
}

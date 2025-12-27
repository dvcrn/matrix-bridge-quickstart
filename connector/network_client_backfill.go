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

	log := nc.log.With().
		Str("portal_id", string(portal.ID)).
		Str("portal_mxid", string(portal.MXID)).
		Bool("forward", fetchParams.Forward).
		Logger()
	ctx = log.WithContext(ctx)
	log.Info().Msg("FetchMessages called")

	backfillMsg1 := &bridgev2.BackfillMessage{
		ID:        networkid.MessageID("static-backfill-1-" + string(portal.ID)),
		Timestamp: now.Add(-3 * time.Minute),
		Sender: bridgev2.EventSender{
			Sender:   networkid.UserID("example-ghost"),
			IsFromMe: false,
		},
		ConvertedMessage: &bridgev2.ConvertedMessage{
			Parts: []*bridgev2.ConvertedMessagePart{{
				Type: event.EventMessage,
				Content: &event.MessageEventContent{
					MsgType: event.MsgText,
					Body:    "Backfilled history (placeholder) #1.",
				},
			}},
		},
	}

	backfillMsg2 := &bridgev2.BackfillMessage{
		ID:        networkid.MessageID("static-backfill-2-" + string(portal.ID)),
		Timestamp: now.Add(-2 * time.Minute),
		Sender: bridgev2.EventSender{
			Sender:   networkid.UserID("example-ghost"),
			IsFromMe: false,
		},
		ConvertedMessage: &bridgev2.ConvertedMessage{
			Parts: []*bridgev2.ConvertedMessagePart{{
				Type: event.EventMessage,
				Content: &event.MessageEventContent{
					MsgType: event.MsgText,
					Body:    "Backfilled history (placeholder) #2.",
				},
			}},
		},
	}

	backfillMsg3 := &bridgev2.BackfillMessage{
		ID:        networkid.MessageID("static-backfill-3-" + string(portal.ID)),
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
					Body:    "Backfilled history (placeholder) #3.",
				},
			}},
		},
	}

	return &bridgev2.FetchMessagesResponse{
		Messages:                []*bridgev2.BackfillMessage{backfillMsg1, backfillMsg2, backfillMsg3},
		HasMore:                 false,
		Forward:                 fetchParams.Forward,
		MarkRead:                !fetchParams.Forward,
		AggressiveDeduplication: true,
	}, nil
}

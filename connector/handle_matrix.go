package connector

import (
	"context"
	"fmt"

	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/event"
)

// HandleMatrixMessage handles incoming messages from Matrix for this user.
func (nc *MyNetworkClient) HandleMatrixMessage(ctx context.Context, msg *bridgev2.MatrixMessage) (*bridgev2.MatrixMessageResponse, error) {
	log := nc.log.With().
		Str("portal_id", string(msg.Portal.ID)).
		Str("sender_mxid", string(msg.Event.Sender)).
		Str("event_id", string(msg.Event.ID)).
		Logger()
	ctx = log.WithContext(ctx)

	log.Info().Msg("HandleMatrixMessage called")

	_, err := nc.bridge.GetExistingUserByMXID(ctx, msg.Event.Sender)
	if err != nil {
		log.Err(err).Str("user_mxid", string(msg.Event.Sender)).Msg("Failed to get user object, ignoring message")
		return nil, nil
	}

	nc.QueueRemoteMessage(ctx, msg.Portal.ID, "Hi there too")

	return &bridgev2.MatrixMessageResponse{}, nil
}

// GetUserInfo is not implemented for this simple connector.
func (nc *MyNetworkClient) GetUserInfo(ctx context.Context, ghost *bridgev2.Ghost) (*bridgev2.UserInfo, error) {
	return nil, fmt.Errorf("user info not available")
}

// GetChatInfo is not implemented for this simple connector.
func (nc *MyNetworkClient) GetChatInfo(ctx context.Context, portal *bridgev2.Portal) (*bridgev2.ChatInfo, error) {
	return nil, fmt.Errorf("chat info not available")
}

// GetCapabilities returns the supported features for chats handled by this client.
func (nc *MyNetworkClient) GetCapabilities(ctx context.Context, portal *bridgev2.Portal) *event.RoomFeatures {
	return &event.RoomFeatures{
		MaxTextLength: 65536,
		Formatting: event.FormattingFeatureMap{
			event.FmtBold:          event.CapLevelFullySupported,
			event.FmtItalic:        event.CapLevelFullySupported,
			event.FmtUnderline:     event.CapLevelFullySupported,
			event.FmtStrikethrough: event.CapLevelFullySupported,
			event.FmtInlineCode:    event.CapLevelFullySupported,
			event.FmtCodeBlock:     event.CapLevelFullySupported,
		},
		Edit:   event.CapLevelFullySupported,
		Reply:  event.CapLevelFullySupported,
		Thread: event.CapLevelFullySupported,
	}
}

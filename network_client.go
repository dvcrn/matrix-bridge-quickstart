package main

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/networkid"
	"maunium.net/go/mautrix/event"
)

// Ensure SimpleNetworkClient implements NetworkAPI
var _ bridgev2.NetworkAPI = (*SimpleNetworkClient)(nil)

// SimpleNetworkClient implements the bridgev2.NetworkAPI for interacting
// with the simple network on behalf of a specific user login.
type SimpleNetworkClient struct {
	log       zerolog.Logger
	bridge    *bridgev2.Bridge
	login     *bridgev2.UserLogin
	connector *SimpleNetworkConnector // Reference back to the connector if needed
}

// Connect is a no-op for this simple connector.
func (nc *SimpleNetworkClient) Connect(ctx context.Context) {
	nc.log.Info().Msg("SimpleNetworkClient Connect called (no-op)")
}

// Disconnect is a no-op for this simple connector.
func (nc *SimpleNetworkClient) Disconnect() {
	nc.log.Info().Msg("SimpleNetworkClient Disconnect called (no-op)")
}

// LogoutRemote is a no-op for this simple connector.
func (nc *SimpleNetworkClient) LogoutRemote(ctx context.Context) {
	nc.log.Info().Msg("SimpleNetworkClient LogoutRemote called (no-op)")
}

// IsThisUser checks if the given remote network user ID belongs to this client instance.
func (nc *SimpleNetworkClient) IsThisUser(ctx context.Context, userID networkid.UserID) bool {
	// Compare with the username stored during login
	return string(userID) == nc.login.RemoteName
}

// IsLoggedIn always returns true for this simple connector.
func (nc *SimpleNetworkClient) IsLoggedIn() bool {
	// In a real bridge, this would check the connection status to the remote network.
	return true
}

// GetUserInfo is not implemented for this simple connector.
func (c *SimpleNetworkClient) GetUserInfo(ctx context.Context, ghost *bridgev2.Ghost) (*bridgev2.UserInfo, error) {
	return nil, fmt.Errorf("user info not available")
}

// GetChatInfo is not implemented for this simple connector.
func (c *SimpleNetworkClient) GetChatInfo(ctx context.Context, portal *bridgev2.Portal) (*bridgev2.ChatInfo, error) {
	return nil, fmt.Errorf("chat info not available")
}

// GetCapabilities returns the supported features for chats handled by this client.
func (c *SimpleNetworkClient) GetCapabilities(ctx context.Context, portal *bridgev2.Portal) *event.RoomFeatures {
	return &event.RoomFeatures{
		MaxTextLength: 65536,
		// Explicitly declare support for text messages and basic formatting/features
		Formatting: event.FormattingFeatureMap{
			event.FmtBold:          event.CapLevelFullySupported,
			event.FmtItalic:        event.CapLevelFullySupported,
			event.FmtUnderline:     event.CapLevelFullySupported,
			event.FmtStrikethrough: event.CapLevelFullySupported,
			event.FmtInlineCode:    event.CapLevelFullySupported,
			event.FmtCodeBlock:     event.CapLevelFullySupported,
		},
		// Add support for basic message features
		Edit:   event.CapLevelFullySupported,
		Reply:  event.CapLevelFullySupported,
		Thread: event.CapLevelFullySupported,
	}
}

// HandleMatrixMessage handles incoming messages from Matrix for this user.
func (nc *SimpleNetworkClient) HandleMatrixMessage(ctx context.Context, msg *bridgev2.MatrixMessage) (*bridgev2.MatrixMessageResponse, error) {
	log := nc.log.With().
		Str("portal_id", string(msg.Portal.ID)).
		Str("sender_mxid", string(msg.Event.Sender)).
		Str("event_id", string(msg.Event.ID)).
		Logger()
	ctx = log.WithContext(ctx)

	log.Info().Msg("HandleMatrixMessage called")

	// Try to get the user associated with the sender MXID
	_, err := nc.bridge.GetExistingUserByMXID(ctx, msg.Event.Sender)
	if err != nil {
		log.Err(err).Str("user_mxid", string(msg.Event.Sender)).Msg("Failed to get user object, ignoring message")
		// ignoring this because we only reply to user messages
		return nil, nil
	}

	// Get the Network User ID for the ghost
	ghostNetworkUserID := networkid.UserID(fmt.Sprintf("%s_ghosty_ghost", nc.connector.GetNetworkID()))

	// Get the ghost object
	ghost, err := nc.bridge.GetExistingGhostByID(ctx, ghostNetworkUserID)
	if err != nil {
		log.Err(err).Str("ghost_network_user_id", string(ghostNetworkUserID)).Msg("Failed to get ghost object to reply, ignoring message")
		return nil, err
	}

	// Prepare the reply message
	replyContent := &event.MessageEventContent{
		MsgType: event.MsgText,
		Body:    "OK",
	}
	content := &event.Content{Parsed: replyContent}

	// Use Portal.MXID which should be the id.RoomID for the portal
	roomID := msg.Portal.MXID
	err = ghost.Intent.EnsureJoined(ctx, roomID)
	if err != nil {
		log.Err(err).Msg("Failed to ensure ghost was joined before replying")
		// For this simple case, we'll try sending anyway.
	}

	// Send the reply using the ghost's intent
	respSendEvent, err := ghost.Intent.SendMessage(ctx, roomID, event.EventMessage, content, nil)
	if err != nil {
		log.Err(err).Msg("Failed to send reply message from ghost")
		// Return the error so the bridge core knows the handler failed
		return nil, fmt.Errorf("failed to send ghost reply: %w", err)
	}

	log.Info().Str("reply_event_id", string(respSendEvent.EventID)).Msg("Successfully sent reply message from ghost")

	// Return an empty response as we don't need to store a corresponding remote message
	return &bridgev2.MatrixMessageResponse{}, nil
}

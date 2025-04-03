package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/rs/zerolog"
	"go.mau.fi/util/configupgrade"
	"go.mau.fi/util/ptr"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/networkid"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

// Ensure SimpleNetworkConnector implements NetworkConnector
var _ bridgev2.NetworkConnector = (*SimpleNetworkConnector)(nil)

// SimpleNetworkConnector implements the NetworkConnector interface
type SimpleNetworkConnector struct {
	log    zerolog.Logger
	bridge *bridgev2.Bridge
}

// NewSimpleNetworkConnector creates a new instance of SimpleNetworkConnector
func NewSimpleNetworkConnector(log zerolog.Logger) *SimpleNetworkConnector {
	return &SimpleNetworkConnector{
		log: log.With().Str("component", "network-connector").Logger(),
	}
}

// Init initializes the connector with the bridge instance
// This method might be called by the bridge core even if not strictly in the interface.
func (c *SimpleNetworkConnector) Init(br *bridgev2.Bridge) {
	c.bridge = br
	c.log = c.bridge.Log
	c.log.Info().Msg("SimpleNetworkConnector Init called")
}

// GetName implements bridgev2.NetworkConnector
func (c *SimpleNetworkConnector) GetName() bridgev2.BridgeName {
	// These values should probably come from config eventually
	return bridgev2.BridgeName{
		DisplayName:          "Simple Bridge",
		NetworkURL:           "https://example.org", // Placeholder
		NetworkIcon:          "",                    // Placeholder
		NetworkID:            "simplenetwork",       // Placeholder
		BeeperBridgeType:     "simple",              // Placeholder
		DefaultPort:          29319,                 // Placeholder
		DefaultCommandPrefix: "!simple",             // Placeholder
	}
}

// GetNetworkID implements bridgev2.NetworkConnector
func (c *SimpleNetworkConnector) GetNetworkID() string {
	return c.GetName().NetworkID
}

// GetCapabilities implements bridgev2.NetworkConnector
// This returns NetworkGeneralCapabilities as required by the interface.
func (c *SimpleNetworkConnector) GetCapabilities() *bridgev2.NetworkGeneralCapabilities {
	// Return general capabilities for the connector itself (not room-specific)
	return &bridgev2.NetworkGeneralCapabilities{
		// DisappearingMessages: false, // Set if supported
		// AggressiveUpdateInfo: false, // Set if needed
	}
}

// GetDBMetaTypes implements bridgev2.NetworkConnector
// Added based on linter error in main.go
func (c *SimpleNetworkConnector) GetDBMetaTypes() database.MetaTypes {
	// Include database import if needed for MetaTypes
	return database.MetaTypes{} // No custom meta types for this simple network connector
}

// --- Login Flow Implementation ---

// GetLoginFlows implements bridgev2.NetworkConnector
func (c *SimpleNetworkConnector) GetLoginFlows() []bridgev2.LoginFlow {
	return []bridgev2.LoginFlow{{
		ID:          LoginFlowIDUsernamePassword,
		Name:        "Username & Password",
		Description: "Log in using a username and password (no actual validation).",
	}}
}

// CreateLogin implements bridgev2.NetworkConnector
func (c *SimpleNetworkConnector) CreateLogin(ctx context.Context, user *bridgev2.User, flowID string) (bridgev2.LoginProcess, error) {
	if flowID != LoginFlowIDUsernamePassword {
		return nil, fmt.Errorf("unsupported login flow ID: %s", flowID)
	}
	// Now returns SimpleLogin defined in login.go
	return &SimpleLogin{
		User: user,
		Main: c, // Pass the connector instance
		Log:  user.Log.With().Str("action", "login").Str("flow", flowID).Logger(),
	}, nil
}

// SimpleLogin struct and methods removed - now in login.go

// --- End Login Flow Implementation ---

// GetConfig implements bridgev2.NetworkConnector
func (c *SimpleNetworkConnector) GetConfig() (string, any, configupgrade.Upgrader) {
	// TODO: Implement proper config handling if network-specific config is needed
	return "simple-config.yaml", nil, nil // Placeholder config file name
}

// GetBridgeInfoVersion implements bridgev2.NetworkConnector
func (c *SimpleNetworkConnector) GetBridgeInfoVersion() (int, int) {
	return 1, 0 // Example major and minor version
}

// Start implements bridgev2.NetworkConnector
// Keeping context here as it was originally present.
func (c *SimpleNetworkConnector) Start(ctx context.Context) error {
	c.log.Info().Msg("SimpleNetworkConnector Start called")
	// TODO: Implement actual startup logic if needed (e.g., connect to the network)
	return nil
}

// Stop implements bridgev2.NetworkConnector
// Keeping context here as it was originally present.
func (c *SimpleNetworkConnector) Stop(ctx context.Context) error {
	c.log.Info().Msg("SimpleNetworkConnector Stop called")
	// TODO: Implement actual shutdown logic if needed
	return nil
}

// LoadUserLogin implements bridgev2.NetworkConnector
// Note: Receives *bridgev2.UserLogin which wraps *database.UserLogin
func (c *SimpleNetworkConnector) LoadUserLogin(ctx context.Context, login *bridgev2.UserLogin) error {
	c.log.Info().
		Str("user_id", string(login.ID)).
		Str("remote_name", login.RemoteName).
		Str("mxid", string(login.User.MXID)).
		Msg("LoadUserLogin called")

	// Create a new SimpleNetworkClient for this login
	client := &SimpleNetworkClient{
		log:       c.log.With().Str("user_id", string(login.ID)).Logger(),
		bridge:    c.bridge,
		login:     login,
		connector: c,
	}

	// Store the client in the login object
	login.Client = client

	c.log.Info().
		Str("user_id", string(login.ID)).
		Str("remote_name", login.RemoteName).
		Interface("client_type", client).
		Msg("Created and stored SimpleNetworkClient")

	// Run welcome logic *after* the login is fully established and loaded
	go c.createWelcomeRoomAndSendIntro(login)

	return nil
}

// createWelcomeRoomAndSendIntro performs the room creation and ghost interaction logic.
func (c *SimpleNetworkConnector) createWelcomeRoomAndSendIntro(login *bridgev2.UserLogin) {
	ctx := context.Background() // Use a background context for this async task
	user := login.User          // Get the User object
	log := c.log.With().Str("user_mxid", string(user.MXID)).Str("login_id", string(login.ID)).Logger()
	ctx = log.WithContext(ctx) // Add logger to context for subsequent calls

	// Seed random number generator for greetings
	// Note: Better to seed once at startup, but fine for this example
	rand.Seed(time.Now().UnixNano())

	log.Info().Msg("Starting welcome room creation process")

	portalId := networkid.PortalID("some-remote-id")
	portalKey := networkid.PortalKey{
		ID: portalId,
	}

	portal, err := c.bridge.GetPortalByKey(ctx, portalKey)
	if err != nil {
		log.Err(err).Str("portal_key", string(portalKey.ID)).Msg("Failed to get portal")
		return
	}

	log.Info().Str("portal_id", string(portal.ID)).Msg("Successfully retrieved portal")

	// 1. Define Ghost details
	ghostNetworkUserID := networkid.UserID(fmt.Sprintf("%s_ghosty_ghost", c.GetNetworkID()))
	ghostDisplayName := "Ghosty Ghost"

	ghost, err := c.bridge.GetGhostByID(ctx, ghostNetworkUserID)
	if err != nil {
		log.Err(err).Str("ghost_network_user_id", string(ghostNetworkUserID)).Msg("Failed to get ghost")
		return
	}

	// Update ghost info if needed
	ghostInfo := &bridgev2.UserInfo{
		Name: &ghostDisplayName,
		// Add other fields as needed - e.g., avatar
	}
	// Use UpdateInfo to update the profile
	ghost.UpdateInfo(ctx, ghostInfo)

	log = log.With().Str("ghost_mxid", string(ghost.ID)).Logger() // MXID field should exist
	log.Info().Msg("Successfully retrieved/provisioned ghost")

	// 2. Create the room
	roomAliasName := fmt.Sprintf("welcome-%s-%d", login.RemoteName, time.Now().UnixNano()) // Unique alias attempt

	// Determine if auto-join via initial_members is likely available
	// You might want to make this check more robust based on MatrixConnector capabilities
	autoJoin := c.bridge.Matrix.GetCapabilities().AutoJoinInvites
	inviteList := []id.UserID{user.MXID} // Always invite the user

	createReq := &mautrix.ReqCreateRoom{
		Visibility:      "private",
		Preset:          "private_chat",
		Name:            fmt.Sprintf("Welcome %s!", login.RemoteName),
		Topic:           "Your special welcome room.",
		Invite:          inviteList,
		IsDirect:        false, // Not a DM (usually management rooms are DM, welcome might not be)
		CreationContent: map[string]interface{}{},
		InitialState:    []*event.Event{},
		RoomAliasName:   roomAliasName, // Local part of the alias
	}

	// If auto-join seems available, use initial_members
	if autoJoin {
		// TODO: Use BeeperInitialMembers when available/needed
		// createReq.BeeperInitialMembers = []id.UserID{user.MXID}
		// createReq.BeeperAutoJoinInvites = true // May still be needed depending on homeserver
		log.Info().Msg("Attempting room creation with auto-join flags")
	} else {
		log.Info().Msg("Auto-join flags not set, will attempt join via double puppet if available")
	}

	resp, err := c.bridge.Bot.CreateRoom(ctx, createReq)
	if err != nil {
		log.Err(err).Msg("Failed to create welcome room")
		return
	}
	roomID := resp
	log = log.With().Str("room_id", string(roomID)).Logger()
	log.Info().Msg("Successfully created welcome room")

	doublePuppet := user.DoublePuppet(ctx) // Get the double puppet intent
	// 3. Attempt to join the room using Double Puppet if auto-join didn't handle it
	if autoJoin {
		if doublePuppet != nil {
			log.Info().Msg("Attempting to join welcome room using double puppet intent")
			err = doublePuppet.EnsureJoined(ctx, roomID)
			if err != nil {
				// Log the error, but continue - the user might join manually
				log.Err(err).Msg("Failed to join welcome room using double puppet intent (user might need to accept invite)")
			} else {
				log.Info().Msg("Successfully joined welcome room using double puppet intent")
			}
		} else {
			log.Warn().Msg("Double puppet intent not available, user will need to manually accept invite")
		}
	} else {
		log.Info().Msg("Skipping explicit join, assuming auto-join worked")
	}

	err = portal.CreateMatrixRoom(ctx, user.GetDefaultLogin(), &bridgev2.ChatInfo{
		Name:  ptr.Ptr(fmt.Sprintf("Welcome %s again!", login.RemoteName)),
		Topic: ptr.Ptr("Your special welcome room."),
		Members: &bridgev2.ChatMemberList{Members: []bridgev2.ChatMember{
			{
				EventSender: bridgev2.EventSender{
					Sender:      networkid.UserID(user.MXID),
					SenderLogin: networkid.UserLoginID(user.GetDefaultLogin().ID),
				},
				Membership: event.MembershipJoin,
				Nickname:   ptr.Ptr(login.RemoteName),
				PowerLevel: ptr.Ptr(100),
				UserInfo: &bridgev2.UserInfo{
					Name: ptr.Ptr(login.RemoteName),
				},
			},
		}},
	})
	if err != nil {
		log.Err(err).Msg("Failed to create matrix room")
		return
	}
	log.Info().Msg("Successfully created matrix room")

	// 4. Invite the Ghost (Ghost's MXID) using the Bot's Intent API
	// Ensure the ghost isn't already invited/joined if createReq didn't handle it
	// (Not strictly necessary here as ghost isn't in initial invite list)
	// _, err = c.bridge.Bot.InviteUser(ctx, roomID, &mautrix.ReqInviteUser{UserID: ghost.MXID})
	// if err != nil {
	// 	// Don't fail the whole process, maybe ghost was already there? Log and continue.
	// 	log.Err(err).Msg("Failed to invite ghost to welcome room (continuing)")
	// } else {
	// 	log.Info().Msg("Invited ghost to welcome room")
	// }
	// Optional: Make the ghost auto-accept if needed via ghost.EnsureJoined(ctx, roomID)

	// 5. Send Message from Ghost using the Ghost's Intent API
	greetings := []string{"Hello there!", "Welcome!", "Greetings!", "Hi!", "Hey!"}
	randomGreeting := greetings[rand.Intn(len(greetings))]
	messageContent := &event.MessageEventContent{
		MsgType: event.MsgText,
		Body:    fmt.Sprintf("%s I'm %s, your friendly welcome bot for the %s bridge.", randomGreeting, ghostDisplayName, c.GetName().DisplayName),
	}

	// Wrap the specific content type in the general event.Content struct
	content := &event.Content{Parsed: messageContent}

	// Ensure ghost is joined before sending message
	err = ghost.Intent.EnsureJoined(ctx, roomID)
	if err != nil {
		log.Err(err).Msg("Failed to ensure ghost was joined before sending message")
		// Decide if you should return or try sending anyway
		// return
	}

	_, err = ghost.Intent.SendMessage(ctx, roomID, event.EventMessage, content, nil) // Pass wrapped content
	if err != nil {
		log.Err(err).Msg("Failed to send welcome message from ghost")
		return
	}
	log.Info().Msg("Successfully sent welcome message from ghost")
}

// HandleMatrixMessage implements NetworkAPI
func (c *SimpleNetworkConnector) HandleMatrixMessage(ctx context.Context, msg *bridgev2.MatrixMessage) (*bridgev2.MatrixMessageResponse, error) {
	log := c.log.With().
		Str("portal_id", string(msg.Portal.ID)).
		Str("sender_mxid", string(msg.Event.Sender)).
		Str("event_id", string(msg.Event.ID)).
		Logger()
	ctx = log.WithContext(ctx)

	log.Info().Msg("HandleMatrixMessage called")

	// try to get the user
	_, err := c.bridge.GetExistingUserByMXID(ctx, msg.Event.Sender)
	if err != nil {
		log.Err(err).Str("user_mxid", string(msg.Event.Sender)).Msg("Failed to get user object, ignoring message")
		// ignoring this because we only reply to user messages
		return nil, nil
	}

	// Get the Network User ID for the ghost
	ghostNetworkUserID := networkid.UserID(fmt.Sprintf("%s_ghosty_ghost", c.GetNetworkID()))

	// Calculate the expected Matrix User ID for the ghost
	ghost, err := c.bridge.GetExistingGhostByID(ctx, ghostNetworkUserID)
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
		// Decide if you should return error or try sending anyway
		// For this simple case, we'll try sending anyway.
	}

	// Use roomID obtained from Portal.MXID
	respSendEvent, err := ghost.Intent.SendMessage(ctx, roomID, event.EventMessage, content, nil)
	if err != nil {
		log.Err(err).Msg("Failed to send reply message from ghost")
		// Return the error so the bridge core knows the handler failed
		return nil, fmt.Errorf("failed to send ghost reply: %w", err)
	}

	log.Info().Str("reply_event_id", string(respSendEvent.EventID)).Msg("Successfully sent reply message from ghost")
	log.Info().Str("reply_event_id", string(respSendEvent.EventID)).Msg("Successfully sent reply message from ghost") // Extract EventID

	// Return nil, nil as we don't need to store a corresponding remote message
	// A real bridge might return a database.Message representing the "OK" confirmation.
	return &bridgev2.MatrixMessageResponse{}, nil
}

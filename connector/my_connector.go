package connector

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/rs/zerolog"
	"go.mau.fi/util/configupgrade"
	"go.mau.fi/util/ptr"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/networkid"
	"maunium.net/go/mautrix/event"
)

// Ensure MyConnector implements NetworkConnector.
var _ bridgev2.NetworkConnector = (*MyConnector)(nil)

// MyConnector implements the NetworkConnector interface.
type MyConnector struct {
	log    zerolog.Logger
	bridge *bridgev2.Bridge
}

// NewMyConnector creates a new instance of MyConnector.
func NewMyConnector(log zerolog.Logger) *MyConnector {
	return &MyConnector{
		log: log.With().Str("component", "network-connector").Logger(),
	}
}

// Init initializes the connector with the bridge instance.
func (c *MyConnector) Init(br *bridgev2.Bridge) {
	c.bridge = br
	c.log = c.bridge.Log
	c.log.Info().Msg("MyConnector Init called")
}

// GetName implements bridgev2.NetworkConnector.
func (c *MyConnector) GetName() bridgev2.BridgeName {
	return bridgev2.BridgeName{
		DisplayName:          "Simple Bridge",
		NetworkURL:           "https://example.org",
		NetworkIcon:          "",
		NetworkID:            "simplenetwork",
		BeeperBridgeType:     "simple",
		DefaultPort:          29319,
		DefaultCommandPrefix: "!simple",
	}
}

// GetNetworkID implements bridgev2.NetworkConnector.
func (c *MyConnector) GetNetworkID() string {
	return c.GetName().NetworkID
}

// GetCapabilities implements bridgev2.NetworkConnector.
func (c *MyConnector) GetCapabilities() *bridgev2.NetworkGeneralCapabilities {
	return &bridgev2.NetworkGeneralCapabilities{}
}

// GetDBMetaTypes implements bridgev2.NetworkConnector.
func (c *MyConnector) GetDBMetaTypes() database.MetaTypes {
	return database.MetaTypes{}
}

// GetLoginFlows implements bridgev2.NetworkConnector.
func (c *MyConnector) GetLoginFlows() []bridgev2.LoginFlow {
	return []bridgev2.LoginFlow{{
		ID:          LoginFlowIDUsernamePassword,
		Name:        "Username & Password",
		Description: "Log in using a username and password (no actual validation).",
	}}
}

// CreateLogin implements bridgev2.NetworkConnector.
func (c *MyConnector) CreateLogin(ctx context.Context, user *bridgev2.User, flowID string) (bridgev2.LoginProcess, error) {
	if flowID != LoginFlowIDUsernamePassword {
		return nil, fmt.Errorf("unsupported login flow ID: %s", flowID)
	}
	return &SimpleLogin{
		User: user,
		Main: c,
		Log:  user.Log.With().Str("action", "login").Str("flow", flowID).Logger(),
	}, nil
}

// GetConfig implements bridgev2.NetworkConnector.
func (c *MyConnector) GetConfig() (string, any, configupgrade.Upgrader) {
	return "simple-config.yaml", nil, nil
}

// GetBridgeInfoVersion implements bridgev2.NetworkConnector.
func (c *MyConnector) GetBridgeInfoVersion() (int, int) {
	return 1, 0
}

// Start implements bridgev2.NetworkConnector.
func (c *MyConnector) Start(ctx context.Context) error {
	c.log.Info().Msg("MyConnector Start called")
	return nil
}

// Stop implements bridgev2.NetworkConnector.
func (c *MyConnector) Stop(ctx context.Context) error {
	c.log.Info().Msg("MyConnector Stop called")
	return nil
}

// LoadUserLogin implements bridgev2.NetworkConnector.
func (c *MyConnector) LoadUserLogin(ctx context.Context, login *bridgev2.UserLogin) error {
	c.log.Info().
		Str("user_id", string(login.ID)).
		Str("remote_name", login.RemoteName).
		Str("mxid", string(login.User.MXID)).
		Msg("LoadUserLogin called")

	client := &MyNetworkClient{
		log:       c.log.With().Str("user_id", string(login.ID)).Logger(),
		bridge:    c.bridge,
		login:     login,
		connector: c,
	}

	login.Client = client

	c.log.Info().
		Str("user_id", string(login.ID)).
		Str("remote_name", login.RemoteName).
		Interface("client_type", client).
		Msg("Created and stored MyNetworkClient")

	go c.createWelcomeRoomAndSendIntro(login)

	return nil
}

// createWelcomeRoomAndSendIntro performs the room creation and ghost interaction logic.
func (c *MyConnector) createWelcomeRoomAndSendIntro(login *bridgev2.UserLogin) {
	ctx := context.Background()
	user := login.User
	log := c.log.With().Str("user_mxid", string(user.MXID)).Str("login_id", string(login.ID)).Logger()
	ctx = log.WithContext(ctx)

	rand.Seed(time.Now().UnixNano())

	log.Info().Msg("Starting welcome room creation process")

	portalID := networkid.PortalID("welcome-room")
	portalKey := networkid.PortalKey{
		ID: portalID,
	}

	portal, err := c.bridge.GetPortalByKey(ctx, portalKey)
	if err != nil {
		log.Err(err).Str("portal_key", string(portalKey.ID)).Msg("Failed to get portal")
		return
	}

	log.Info().Str("portal_id", string(portal.ID)).Msg("Successfully retrieved portal")

	ghostNetworkUserID := networkid.UserID(fmt.Sprintf("%s_ghosty_ghost", c.GetNetworkID()))
	ghostDisplayName := "Ghosty Ghost"

	ghost, err := c.bridge.GetGhostByID(ctx, ghostNetworkUserID)
	if err != nil {
		log.Err(err).Str("ghost_network_user_id", string(ghostNetworkUserID)).Msg("Failed to get ghost")
		return
	}

	ghostInfo := &bridgev2.UserInfo{
		Name: &ghostDisplayName,
	}
	ghost.UpdateInfo(ctx, ghostInfo)

	log = log.With().Str("ghost_mxid", string(ghost.ID)).Logger()
	log.Info().Msg("Successfully retrieved/provisioned ghost")

	err = portal.CreateMatrixRoom(ctx, user.GetDefaultLogin(), &bridgev2.ChatInfo{
		Name:  ptr.Ptr(fmt.Sprintf("Welcome %s!", login.RemoteName)),
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

	err = ghost.Intent.EnsureJoined(ctx, portal.MXID)
	if err != nil {
		log.Err(err).Msg("Failed to ensure ghost was joined before sending message")
		return
	}

	greetings := []string{"Hello there!", "Welcome!", "Greetings!", "Hi!", "Hey!"}
	randomGreeting := greetings[rand.Intn(len(greetings))]
	messageContent := &event.MessageEventContent{
		MsgType: event.MsgText,
		Body:    fmt.Sprintf("%s I'm %s, your friendly welcome bot for the %s bridge.", randomGreeting, ghostDisplayName, c.GetName().DisplayName),
	}

	content := &event.Content{Parsed: messageContent}

	_, err = ghost.Intent.SendMessage(ctx, portal.MXID, event.EventMessage, content, nil)
	if err != nil {
		log.Err(err).Msg("Failed to send welcome message from ghost")
		return
	}
	log.Info().Msg("Successfully sent welcome message from ghost")
}

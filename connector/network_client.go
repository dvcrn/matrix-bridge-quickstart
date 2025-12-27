package connector

import (
	"context"

	"github.com/rs/zerolog"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/networkid"
)

// Ensure MyNetworkClient implements NetworkAPI.
var _ bridgev2.NetworkAPI = (*MyNetworkClient)(nil)

// MyNetworkClient implements the bridgev2.NetworkAPI for interacting
// with the simple network on behalf of a specific user login.
type MyNetworkClient struct {
	log       zerolog.Logger
	bridge    *bridgev2.Bridge
	login     *bridgev2.UserLogin
	connector *MyConnector
}

// Connect is a no-op for this simple connector.
func (nc *MyNetworkClient) Connect(ctx context.Context) {
	nc.log.Info().Msg("MyNetworkClient Connect called (no-op)")
}

// Disconnect is a no-op for this simple connector.
func (nc *MyNetworkClient) Disconnect() {
	nc.log.Info().Msg("MyNetworkClient Disconnect called (no-op)")
}

// LogoutRemote is a no-op for this simple connector.
func (nc *MyNetworkClient) LogoutRemote(ctx context.Context) {
	nc.log.Info().Msg("MyNetworkClient LogoutRemote called (no-op)")
}

// IsThisUser checks if the given remote network user ID belongs to this client instance.
func (nc *MyNetworkClient) IsThisUser(ctx context.Context, userID networkid.UserID) bool {
	return string(userID) == nc.login.RemoteName
}

// IsLoggedIn always returns true for this simple connector.
func (nc *MyNetworkClient) IsLoggedIn() bool {
	return true
}

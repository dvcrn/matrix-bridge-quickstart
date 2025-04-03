package main

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"go.mau.fi/util/configupgrade"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	// "maunium.net/go/mautrix/bridgev2/database" // Only needed if LoadUserLogin uses DB types
	// "maunium.net/go/mautrix/bridgev2/networkid" // Only needed if methods use networkid types directly
)

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
	return &bridgev2.NetworkGeneralCapabilities{} // Empty capabilities for now
}

// GetDBMetaTypes implements bridgev2.NetworkConnector
// Added based on linter error in main.go
func (c *SimpleNetworkConnector) GetDBMetaTypes() database.MetaTypes {
	// Include database import if needed for MetaTypes
	return database.MetaTypes{} // No custom meta types for this simple network connector
}

// GetLoginFlows implements bridgev2.NetworkConnector
func (c *SimpleNetworkConnector) GetLoginFlows() []bridgev2.LoginFlow {
	return []bridgev2.LoginFlow{} // No login flows supported
}

// CreateLogin implements bridgev2.NetworkConnector
func (c *SimpleNetworkConnector) CreateLogin(ctx context.Context, user *bridgev2.User, identifier string) (bridgev2.LoginProcess, error) {
	return nil, fmt.Errorf("login not supported in this simple implementation")
}

// GetConfig implements bridgev2.NetworkConnector
func (c *SimpleNetworkConnector) GetConfig() (string, any, configupgrade.Upgrader) {
	// TODO: Implement proper config handling
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
func (c *SimpleNetworkConnector) LoadUserLogin(ctx context.Context, login *bridgev2.UserLogin) error {
	// This simple bridge doesn't handle user logins persistently
	c.log.Debug().Str("user_id", string(login.ID)).Msg("LoadUserLogin called (no-op)")
	return nil
}

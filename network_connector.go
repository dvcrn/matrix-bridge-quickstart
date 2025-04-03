package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"go.mau.fi/util/configupgrade"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/networkid"
	"maunium.net/go/mautrix/bridgev2/status"
	// "maunium.net/go/mautrix/bridgev2/database" // Only needed if LoadUserLogin uses DB types
	// "maunium.net/go/mautrix/bridgev2/networkid" // Only needed if methods use networkid types directly
)

// Login Flow/Step IDs
const (
	LoginFlowIDUsernamePassword = "user-pass"
	LoginStepIDUsernamePassword = "user-pass-input"
	LoginStepIDComplete         = "complete"
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
	return &SimpleLogin{
		User: user,
		Main: c,
		Log:  user.Log.With().Str("action", "login").Str("flow", flowID).Logger(),
	}, nil
}

// SimpleLogin represents an ongoing username/password login attempt.
type SimpleLogin struct {
	User *bridgev2.User
	Main *SimpleNetworkConnector
	Log  zerolog.Logger
}

// Ensure SimpleLogin implements the required interface
var _ bridgev2.LoginProcessUserInput = (*SimpleLogin)(nil)

// Start implements bridgev2.LoginProcessUserInput
func (sl *SimpleLogin) Start(ctx context.Context) (*bridgev2.LoginStep, error) {
	sl.Log.Debug().Msg("Starting username/password login flow")
	return &bridgev2.LoginStep{
		Type:         bridgev2.LoginStepTypeUserInput,
		StepID:       LoginStepIDUsernamePassword,
		Instructions: "Enter your username and password for the 'Simple Network'.",
		UserInputParams: &bridgev2.LoginUserInputParams{
			Fields: []bridgev2.LoginInputDataField{
				{
					Type: bridgev2.LoginInputFieldTypeUsername, // Correct type based on login.go
					ID:   "username",
					Name: "Username",
				},
				{
					Type: bridgev2.LoginInputFieldTypePassword, // Correct type based on login.go
					ID:   "password",
					Name: "Password",
				},
			},
		},
	}, nil
}

// SubmitUserInput implements bridgev2.LoginProcessUserInput
func (sl *SimpleLogin) SubmitUserInput(ctx context.Context, input map[string]string) (*bridgev2.LoginStep, error) {
	username := input["username"]
	password := input["password"] // We don't actually use the password here
	_ = password                  // Explicitly ignore unused variable

	if username == "" {
		return nil, fmt.Errorf("username cannot be empty") // Basic validation example
	}

	sl.Log.Info().Str("username", username).Msg("Received login credentials (no actual validation performed)")

	// In a real bridge, you would authenticate with the remote network here.
	// Since this is simple, we just generate a unique ID based on the username.
	// We need a stable way to generate the LoginID for a given remote identifier (username).
	// Using a UUID based on a namespace and the username ensures this.
	// IMPORTANT: Do NOT just use the raw username, as it might contain invalid characters for a Matrix Localpart.
	// Also, avoid collisions if usernames are not unique across different contexts (not an issue here).

	// Use a fixed UUID namespace for this bridge type
	namespace := uuid.MustParse("f7a4f3e3-5d5a-4a9e-8d8a-3b0b9e8a1b2c") // Example namespace UUID
	loginIDStr := uuid.NewSHA1(namespace, []byte(strings.ToLower(username))).String()
	// Correct type is networkid.UserLoginID
	var loginID networkid.UserLoginID = networkid.UserLoginID(loginIDStr)

	// Create the UserLogin entry in the bridge database
	ul, err := sl.User.NewLogin(ctx, &database.UserLogin{
		ID:         loginID,
		RemoteName: username, // Use the provided username as the display name
		RemoteProfile: status.RemoteProfile{
			Name: username,
			// Add other profile fields if known (e.g., avatar URL)
		},
		// Metadata: &YourUserLoginMetadata{ ... }, // Add if you have custom metadata
	}, &bridgev2.NewLoginParams{
		DeleteOnConflict: false, // Or true if you want relogins to replace old ones
	})
	if err != nil {
		sl.Log.Err(err).Msg("Failed to create user login entry")
		return nil, fmt.Errorf("failed to create user login: %w", err)
	}

	sl.Log.Info().Str("login_id", string(ul.ID)).Msg("Successfully 'logged in' and created user login")

	// Load the user login into memory (important!)
	// In a real bridge, this would trigger connecting to the remote network for this user.
	err = sl.Main.LoadUserLogin(ctx, ul)
	if err != nil {
		// Log the error, but maybe still return success to the user? Depends on desired UX.
		sl.Log.Err(err).Msg("Failed to load user login after creation (this might indicate an issue)")
		// Optionally delete the login record if loading failed critically:
		// sl.User.DeleteLogin(ctx, ul.ID)
		// return nil, fmt.Errorf("failed to activate user login: %w", err)
	}

	return &bridgev2.LoginStep{
		Type:         bridgev2.LoginStepTypeComplete,
		StepID:       LoginStepIDComplete,
		Instructions: fmt.Sprintf("Successfully logged in as '%s'", username),
		CompleteParams: &bridgev2.LoginCompleteParams{
			UserLoginID: ul.ID,
			UserLogin:   ul, // Pass the loaded UserLogin back
		},
	}, nil
}

// Cancel implements bridgev2.LoginProcessUserInput
func (sl *SimpleLogin) Cancel() {
	sl.Log.Debug().Msg("Login process cancelled")
	// Add any cleanup logic here if needed (e.g., aborting network connections)
}

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
	// In a real bridge, this is where you would establish a connection to the
	// remote network for the specific user 'login'.
	// You'd store the connection object/client within the login.Client field
	// or manage it separately mapped by login.ID.
	c.log.Info().Str("user_id", string(login.ID)).Str("remote_name", login.RemoteName).Msg("LoadUserLogin called")

	// Example: Storing a dummy client object
	// login.Client = &YourNetworkClient{ ... connection details ... }

	// Note: SetRemoteStatus doesn't exist on User or UserLogin.
	// Status is typically managed via login.BridgeState.Set(...) in a real bridge.

	return nil
}

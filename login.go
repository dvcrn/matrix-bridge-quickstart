package main

import (
	"context"
	"fmt"
	"strings"

	// Added time for createWelcomeRoomAndSendIntro call
	"github.com/google/uuid"
	"github.com/rs/zerolog" // Added ptr for createWelcomeRoomAndSendIntro call
	// Added mautrix for createWelcomeRoomAndSendIntro call
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/networkid"
	"maunium.net/go/mautrix/bridgev2/status"
	// Added event for createWelcomeRoomAndSendIntro call
	// Added id for createWelcomeRoomAndSendIntro call
)

// Login Flow/Step IDs - Copied from network_connector.go as they are used here
const (
	LoginFlowIDUsernamePassword = "user-pass"
	LoginStepIDUsernamePassword = "user-pass-input"
	LoginStepIDComplete         = "complete"
)

// SimpleLogin represents an ongoing username/password login attempt.
type SimpleLogin struct {
	User *bridgev2.User
	Main *SimpleNetworkConnector // Needs access to the connector for LoadUserLogin and createWelcomeRoom...
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
	err = sl.Main.LoadUserLogin(ctx, ul) // Calls method on SimpleNetworkConnector
	if err != nil {
		// Log the error, but maybe still return success to the user? Depends on desired UX.
		sl.Log.Err(err).Msg("Failed to load user login after creation (this might indicate an issue)")
		// Optionally delete the login record if loading failed critically:
		// sl.User.DeleteLogin(ctx, ul.ID)
		// return nil, fmt.Errorf("failed to activate user login: %w", err)
	}

	// Run welcome logic *after* the login is fully established and loaded
	// This needs access to the connector instance (sl.Main)
	go sl.Main.createWelcomeRoomAndSendIntro(ul)

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

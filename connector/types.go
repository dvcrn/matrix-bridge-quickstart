package connector

import (
	"time"

	"maunium.net/go/mautrix/bridgev2/networkid"
	"maunium.net/go/mautrix/id"
)

// LoginMetadata stores additional remote-auth/session data in the bridge database.
// Use this for connector-specific fields that don't fit in core bridge tables.
type LoginMetadata struct {
	RemoteUserID string     `json:"remote_user_id,omitempty"`
	AccessToken  string     `json:"access_token,omitempty"`
	ExpiresAt    time.Time  `json:"expires_at,omitempty"`
	DeviceID     string     `json:"device_id,omitempty"`
	Scopes       []string   `json:"scopes,omitempty"`
	LastSyncAt   *time.Time `json:"last_sync_at,omitempty"`
}

// New creates a new instance for database registration.
func (m *LoginMetadata) New() any {
	return &LoginMetadata{}
}

// GhostMetadata stores additional remote metadata for a Matrix ghost.
type GhostMetadata struct {
	RemoteUserID string `json:"remote_user_id,omitempty"`
	RemoteName   string `json:"remote_name,omitempty"`
	AvatarURL    string `json:"avatar_url,omitempty"`
}

// New creates a new instance for database registration.
func (m *GhostMetadata) New() any {
	return &GhostMetadata{}
}

// PortalMetadata stores additional remote metadata for a Matrix portal (room).
type PortalMetadata struct {
	RemoteRoomID  string              `json:"remote_room_id,omitempty"`
	OtherUserID   networkid.UserID    `json:"other_user_id,omitempty"`
	InitialName   string              `json:"initial_name,omitempty"`
	InitialAvatar id.ContentURIString `json:"initial_avatar_mxc,omitempty"`
	CreatedAt     time.Time           `json:"created_at,omitempty"`
	LastSyncedAt  *time.Time          `json:"last_synced_at,omitempty"`
	RemoteTopic   string              `json:"remote_topic,omitempty"`
	LastMessageID string              `json:"last_message_id,omitempty"`
	Tags          []string            `json:"tags,omitempty"`
	Notes         map[string]string   `json:"notes,omitempty"`
}

// New creates a new instance for database registration.
func (m *PortalMetadata) New() any {
	return &PortalMetadata{}
}

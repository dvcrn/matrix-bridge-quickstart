# Matrix Bridge Development Notes

## Room Creation in mautrix Bridges

There are two distinct ways to create a room in the mautrix bridge architecture:

1. **Direct Matrix Room Creation** (via `c.bridge.Bot.CreateRoom`):
   - Creates a standard Matrix room using the bridge bot's credentials
   - The room exists only in the Matrix network
   - It has no connection to the remote network or the bridge's routing system
   - Suitable for management rooms, welcome rooms, or other bridge-utility rooms
   - **Messages in these rooms won't trigger `HandleMatrixMessage`**

2. **Portal Room Creation** (via `portal.CreateMatrixRoom`):
   - Creates a "portal" room that's registered with the bridge
   - Links a Matrix room to a specific remote chat/conversation
   - Stores the mapping in the bridge database
   - Sets up proper event routing for messages, reactions, etc.
   - **Only messages in these rooms will trigger `HandleMatrixMessage`**

## Why `HandleMatrixMessage` Isn't Being Called

If your `HandleMatrixMessage` function isn't being called when a new message arrives, it's likely because:

1. The messages are being sent in a room that wasn't created as a portal room
2. The bridge doesn't know to route events from that room to your handler

### Solution:

To properly create a portal room:

```go
// Create a portal with a remote room identifier
portalKey := &database.PortalKey{
    NetworkID: c.GetNetworkID(),
    RoomID:    networkid.RoomID("some-remote-id"), // Remote room identifier
}

portal, err := c.bridge.GetPortalByKey(ctx, portalKey, true)
if err != nil {
    log.Err(err).Msg("Failed to get/create portal")
    return
}

// Link the Matrix room to this portal
// You can either create a new room:
err = portal.CreateMatrixRoom(ctx, nil)
// Or link an existing room:
err = portal.CreateMatrixRoom(ctx, nil, &existingRoomID)
```

The bridge's architecture routes messages through portal objects. Without properly registering your room as a portal, the messages won't reach your `HandleMatrixMessage` function.

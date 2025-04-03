package main

import (
	// Keep only necessary imports for mxmain setup
	"maunium.net/go/mautrix/bridgev2/matrix/mxmain"
	// We might need zerolog if we create the connector instance here with a logger,
	// but mxmain likely handles logging setup internally.
)

// Build time variables (optional but good practice)
var (
	Tag       = "unknown"
	Commit    = "unknown"
	BuildTime = "unknown"
)

func main() {
	// Create the network connector instance
	// mxmain will handle logger injection if the connector needs it during its Init phase.
	connector := &SimpleNetworkConnector{}

	// Create and configure the BridgeMain helper
	m := mxmain.BridgeMain{
		Name:        "minibridge", // Choose a name
		Description: "A minimal mautrix-go bridge example.",
		Version:     "0.1.0",   // Choose a version
		URL:         "",        // Add your repo URL if you have one
		Connector:   connector, // Pass the network connector

		// Optional hooks (like in matrix-whatsapp example)
		// PostInit: func() { ... },
		// PostStart: func() { ... },
	}

	// Initialize version info and run the bridge
	m.InitVersion(Tag, Commit, BuildTime) // Pass build vars
	m.Run()
}

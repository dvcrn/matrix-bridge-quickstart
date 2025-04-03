package main

import (
	// Import the main bridge helper package
	"maunium.net/go/mautrix/bridgev2/matrix/mxmain"
	// mxmain handles logging setup internally.
)

// Build time variables (optional but good practice)
var (
	Tag       = "unknown"
	Commit    = "unknown"
	BuildTime = "unknown"
)

func main() {
	// Create the network connector instance.
	// mxmain will handle logger injection during Init.
	connector := &SimpleNetworkConnector{}

	// Create and configure the BridgeMain helper.
	m := mxmain.BridgeMain{
		Name:        "minibridge",
		Description: "A minimal mautrix-go bridge example.",
		Version:     "0.1.0",
		URL:         "", // Optional: Add your repo URL
		Connector:   connector,

		// Optional hooks can be added here if needed.
		// PostInit: func() { ... },
	}

	// Initialize version info and run the bridge.
	m.InitVersion(Tag, Commit, BuildTime)
	m.Run()
}

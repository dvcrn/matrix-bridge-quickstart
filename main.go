package main

import (
	"maunium.net/go/mautrix/bridgev2/matrix/mxmain"

	"github.com/dvcrn/matrix-bridge-quickstart/connector"
)

var (
	Tag       = "unknown"
	Commit    = "unknown"
	BuildTime = "unknown"
)

func main() {
	connector := &connector.MyConnector{}

	m := mxmain.BridgeMain{
		Name:        "minibridge",
		Description: "A minimal mautrix-go bridge example.",
		Version:     "0.1.0",
		URL:         "",
		Connector:   connector,
	}

	m.InitVersion(Tag, Commit, BuildTime)
	m.Run()
}

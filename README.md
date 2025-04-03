# Minimal Mautrix-Go Bridge: Quickstart Template 🚀

[![Go Reference](https://pkg.go.dev/badge/github.com/mautrix/go.svg)](https://pkg.go.dev/github.com/mautrix/go)

Welcome! This project provides a minimal, bare-bones template for creating a [Matrix](https://matrix.org/) bridge using the powerful [mautrix-go](https://github.com/mautrix/go) library, specifically leveraging its modern `bridgev2` framework.

**What is a Matrix Bridge?**

A Matrix bridge connects the decentralized, open Matrix communication network to other, often proprietary, chat networks (like WhatsApp, Telegram, Discord, etc.). It acts as a translator, allowing users on Matrix to communicate with users on the other network, and vice-versa.

Building a bridge involves a lot of standard setup: handling Matrix connections, managing user logins, storing data, processing configuration, etc. This template handles that common boilerplate for you, letting you jump straight into the interesting part: connecting to *your* specific target network.

## ⚙️ Project Structure

Here's a breakdown of the key files:

*   **`main.go`**:
    *   The main entry point. Handles command-line flags, configuration, logging, and the bridge's start/stop lifecycle using `mxmain`.
    *   You usually **won't need to modify this** much initially.

*   **`network_connector.go`**:
    *   **⭐ This is the main implementation of the network side of the bridge! ⭐**
    *   Contains the `SimpleNetworkConnector` struct, which implements the `bridgev2.NetworkConnector` interface.
    *   This file defines the bridge's core properties (`GetName`, `GetCapabilities`), handles loading user sessions (`LoadUserLogin`), and initiates the login process (`GetLoginFlows`, `CreateLogin`).
    *   It also contains the main logic for handling events *from* Matrix (`HandleMatrixMessage`, etc. - though these might be delegated).

*   **`login.go`**:
    *   Contains the logic for specific login flows (e.g., `SimpleLogin` for username/password).
    *   Implements `bridgev2.LoginProcess` interfaces to handle steps like asking for user input (`Start`, `SubmitUserInput`) and finalizing login.

*   **`network_client.go`** (Optional but Recommended):
    *   This file typically holds the client logic for interacting with the *remote network* for a *specific logged-in user*.
    *   You'd create a struct (e.g., `SimpleNetworkClient`) that implements `bridgev2.NetworkClient`.
    *   Methods here would handle sending messages *to* the remote network, fetching user/room info, handling typing notifications, etc., based on Matrix events forwarded from `network_connector.go`.
    *   The `LoadUserLogin` method in `network_connector.go` would instantiate this client.

---

## 🚀 Getting Started: Building Your Bridge

Follow these steps to get your basic bridge running:

1.  **Clone/Copy Template:**
    *   Get a local copy of this template directory (e.g., `git clone ...` or download ZIP).

2.  **Implement Your Connector (`network_connector.go`):**
    *   Open `network_connector.go`. This is where you'll spend most of your time.
    *   **Goal:** Replace the placeholder logic with real code to interact with your target network.
    *   Start by filling in:
        *   `GetName()`: Provide accurate details about your bridge and the network it connects to.
        *   `GetCapabilities()`: Define what features your bridge supports (e.g., message formatting, read receipts).
        *   `GetLoginFlows()` / `CreateLogin()`: Implement the actual login mechanism for your target network. The current example is just a placeholder!
        *   `LoadUserLogin()`: This is crucial. When a user logs in, this function should establish their *persistent* connection to the remote network.
        *   `Start()` / `Stop()`: Add any global setup/teardown logic for your network connection.
    *   **Configuration:** If your network needs API keys or other settings, implement `GetConfig()` to load them from a file (like `simple-config.yaml`) and create that YAML file.

3.  **Generate Registration File:**
    *   Open your terminal in the project directory.
    *   Run: `go run . -g -c config.yaml -r registration.yaml`
    *   This creates the initial `registration.yaml`. **Keep this file safe!**

4.  **Configure the Bridge (`config.yaml`):**
    *   Edit `config.yaml`.
    *   Set `homeserver.address` (e.g., `https://matrix.example.com`) and `homeserver.domain` (e.g., `matrix.example.com`).
    *   **Crucial:** Copy the `id`, `as_token`, `hs_token` from the *generated* `registration.yaml` into the `appservice` section of `config.yaml`. Also, copy `bot.username` and potentially adjust `username_template`.
    *   Review and adjust `database` (default is `./simple-bridge.db`), `logging`, and `permissions` as needed.
    *   If you created a network-specific config file (e.g., `simple-config.yaml`), configure its settings now.

5.  **Configure Your Homeserver:**
    *   Copy the generated `registration.yaml` file to your Matrix homeserver's configuration directory.
        *   For Synapse, this is often `/etc/synapse/conf.d/` or similar. Check your homeserver's documentation.
    *   **Restart your homeserver** software (e.g., `systemctl restart synapse`). This makes it load the registration file and know about your bridge.

6.  **Build the Bridge:**
    *   In the project directory, run: `go build`
    *   This creates an executable binary (e.g., `minibridge`).

7.  **Run the Bridge:**
    *   Execute the binary, pointing it to your config files:
        ```bash
        ./minibridge -c config.yaml -r registration.yaml
        ```
    *   Check the terminal output for logs and potential errors.

🎉 **Congratulations!** You have a running (though perhaps very basic) Matrix bridge.

---

## ⏭️ Next Steps

*   **Flesh out `network_connector.go`:** Implement message handling, user/room synchronization, presence, typing notifications, etc.
*   **Consult `mautrix-go` Docs:** Explore the `bridgev2` package documentation for detailed information on interfaces and helpers: [pkg.go.dev/maunium.net/go/mautrix/bridgev2](https://pkg.go.dev/maunium.net/go/mautrix/bridgev2)
*   **Study Other Bridges:** Look at the source code of other `mautrix-go` based bridges (like `mautrix-whatsapp`, `mautrix-telegram`) for inspiration and examples.
*   **Testing:** Implement unit and integration tests for your connector logic.
*   **Refine Configuration:** Make your bridge more robust by handling configuration validation and updates.

Good luck with your bridge development!
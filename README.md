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
    *   **⭐ This is the heart of your bridge logic! ⭐**
    *   Contains the `SimpleNetworkConnector` struct, which implements the `bridgev2.NetworkConnector` interface.
    *   You'll fill in the methods here (`Start`, `Stop`, `GetName`, `GetCapabilities`, `CreateLogin`, `LoadUserLogin`, etc.) to:
        *   Connect to your target network's API.
        *   Handle authentication (logging users in).
        *   Manage user identities and profiles.
        *   Translate messages and events between Matrix and the remote network.
    *   Includes a *basic, non-functional placeholder* for username/password login to demonstrate the flow.

    *   **Placeholder Login Flow:** The `SimpleLogin` struct implements a basic username/password flow. It generates a unique ID from the username and saves the login details (linking Matrix user to remote user) in the bridge database (`database.UserLogin`). 


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
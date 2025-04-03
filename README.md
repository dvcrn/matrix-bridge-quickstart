# Minimal Mautrix-Go Bridge Quickstart

This project provides a minimal, bare-bones template for creating a Matrix bridge using the [mautrix-go](https://github.com/mautrix/go) library, specifically leveraging the `bridgev2` framework.

It's designed to be a starting point, handling the basic boilerplate of setting up a bridge process so you can focus on implementing the connection to your target network.

## Project Structure

*   **`main.go`**:
    *   The main entry point of the bridge application.
    *   Uses `maunium.net/go/mautrix/bridgev2/matrix/mxmain` to handle command-line flags, configuration loading, logging setup, and the main bridge lifecycle (start/stop).
    *   Instantiates the `SimpleNetworkConnector`.
    *   You generally won't need to modify this file much unless you need custom initialization or lifecycle hooks.

*   **`network_connector.go`**:
    *   Contains the `SimpleNetworkConnector` struct.
    *   This struct implements the `bridgev2.NetworkConnector` interface from mautrix-go.
    *   **This is the primary file you need to modify.** It contains placeholder implementations for methods like `Start`, `Stop`, `GetName`, `GetCapabilities`, `CreateLogin`, `LoadUserLogin`, etc.
    *   You will implement the logic here to connect to your specific network, handle authentication, manage users, and translate messages between Matrix and the remote network.
    *   **Login Flow:**
        *   Implements a basic username/password login (`GetLoginFlows`, `CreateLogin`, `SimpleLogin` struct).
        *   The `SimpleLogin.Start` method prompts the user for a username and password via the Matrix client.
        *   The `SimpleLogin.SubmitUserInput` method receives the input.
        *   **Important:** This implementation does **not** validate the password against any real network. It's a placeholder.
        *   It generates a stable, unique internal ID (`networkid.UserLoginID`) for the remote user based on the provided username using a SHA1 UUID hash with a fixed namespace. This ensures the same username always maps to the same internal ID.
        *   It calls `user.NewLogin` (where `user` is the `bridgev2.User` representing the Matrix user performing the login).
        *   `user.NewLogin` creates and saves a `database.UserLogin` record in the bridge's database (`simple-bridge.db` by default, configured in `config.yaml`). This record links the Matrix user (`UserMXID`) to the generated `UserLoginID`, storing the provided username as `RemoteName` and basic profile info.
        *   Finally, it calls `LoadUserLogin` for the newly created login. In this simple connector, `LoadUserLogin` just logs that the user was loaded; in a real bridge, this is where you would establish the actual connection to the remote network for that user.

*   **`config.yaml`**:
    *   The main configuration file for the bridge core and Matrix connection details (homeserver, database, permissions, etc.).
    *   Follows the standard mautrix-go configuration format.
    *   You **must** configure this file, especially the `homeserver` and `appservice` sections.

*   **`registration.yaml`**:
    *   The Matrix Application Service registration file.
    *   This file **must be generated** (e.g., using `go run . -g`) and placed in your homeserver's configuration directory.
    *   It tells the homeserver how to communicate with your bridge (URL, tokens, user/room namespaces).

*   **`simple-config.yaml` (Placeholder)**:
    *   Defined in `network_connector.go::GetConfig` as the network-specific config file.
    *   Currently empty/unused. You would define and load your network-specific settings (API keys, connection details) here if needed.

*   **`CONFIGURATION_SUMMARY.md`**:
    *   A summary of the configuration changes applied during initial setup (if generated previously).

## Getting Started

1.  **Clone/Copy:** Get a copy of this `v2` directory.
2.  **Implement Connector:**
    *   Open `network_connector.go`.
    *   Fill in the placeholder methods (`GetName`, `GetCapabilities`, `GetLoginFlows`, `CreateLogin`, `LoadUserLogin`, `Start`, `Stop`, etc.) with the logic specific to the network you are bridging.
    *   Decide if you need network-specific configuration and implement `GetConfig` accordingly, creating the corresponding YAML file (like `simple-config.yaml`).
3.  **Generate Registration:**
    *   Run `go run . -g -c config.yaml -r registration.yaml`. This will generate the `registration.yaml` file based on defaults and your `config.yaml`.
4.  **Configure Bridge:**
    *   Edit `config.yaml`.
        *   Update `homeserver.address` and `homeserver.domain`.
        *   Copy the `id`, `as_token`, `hs_token`, `bot.username`, and `username_template` from the generated `registration.yaml` into the `appservice` section of `config.yaml`.
        *   Adjust `permissions`, `database`, `logging`, and other settings as needed.
        *   Configure any network-specific settings in your network config file (e.g., `simple-config.yaml`) if you created one.
5.  **Configure Homeserver:**
    *   Copy the generated `registration.yaml` to your Matrix homeserver's appservice configuration directory (e.g., `/etc/synapse/conf.d/` or similar).
    *   Restart/reload your homeserver to make it aware of the new appservice.
6.  **Build:**
    *   Run `go build` in the `v2` directory.
7.  **Run:**
    *   Execute the compiled binary: `./minibridge -c config.yaml -r registration.yaml` (or `./v2 -c ...` if building from the parent directory).

Now you have a running (though likely basic) Matrix bridge! Monitor the logs and continue implementing features in `network_connector.go`. 
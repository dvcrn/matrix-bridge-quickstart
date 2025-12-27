# Repository Guidelines

## Project Structure & Module Organization
- `main.go` is the entry point; it wires config, logging, and bridge lifecycle (`mxmain`).
- `connector/my_connector.go` implements the core `bridgev2.NetworkConnector` logic.
- `connector/login.go` defines login flows (`bridgev2.LoginProcess`).
- `connector/network_client.go` is the per-user remote network client.
- `go.mod`/`go.sum` manage dependencies; no dedicated `tests/` directory yet.

## Build, Test, and Development Commands
- `go run . -g -c config.yaml -r registration.yaml` generates `registration.yaml`.
- `go run . -c config.yaml -r registration.yaml` runs the bridge from source.
- `go build` builds the binary (default name `minibridge`).
- `go test ./...` runs tests (currently no test files).
- Toolchain: `mise install` and `mise exec go@1.25.5 -- <cmd>`.

## Coding Style & Naming Conventions
- Go code follows `gofmt` (tabs; standard Go formatting).
- Filenames use snake_case (e.g., `network_client.go`).
- Keep new files in the repo root unless you add a new package.

## Bridge Build & Operation Guide
- Implement the network side in `connector/my_connector.go`: `GetName`, `GetCapabilities`, `GetLoginFlows`, `CreateLogin`, `LoadUserLogin`, `Start`, `Stop`.
- If you need network-specific settings, add a config file (e.g., `simple-config.yaml`) and load it in `GetConfig()`.
- Configure `config.yaml` with `homeserver.address`, `homeserver.domain`, database path, and permissions.
- Copy `id`, `as_token`, `hs_token`, and bot settings from `registration.yaml` into `config.yaml`.
- Place `registration.yaml` into your homeserver config directory and restart the homeserver.

## BridgeV2 Lifecycle & Backfill Notes
- `LoadUserLogin` should create your adapter, assign `userLogin.Client`, and avoid blocking network calls; connect in `NetworkAPI.Connect`.
- Create new rooms by queuing `simplevent.ChatResync` with `EventMeta.CreatePortal = true` and `ChatInfo = nil`; the framework will call `GetChatInfo`.
- `GetChatInfo` should return the full desired room state (name, topic, DM type, members, power levels, avatar).
- Use `simplevent.ChatInfoChange` for state-only updates (e.g., lock rooms on unmatch).
- Implement `FetchMessages` (via `BackfillingNetworkAPI`) to backfill history; lock portals, convert remote messages to `BackfillMessage`, and return `HasMore`.
- Send bridge state updates with `userLogin.BridgeState.Send(status.BridgeState{StateEvent: ...})` for connect/backfill/errors.

## Portal Rooms vs Direct Rooms
- **Direct Matrix rooms** (`c.bridge.Bot.CreateRoom`) are not routed through the bridge and will not trigger `HandleMatrixMessage`.
- **Portal rooms** (`portal.CreateMatrixRoom`) are registered in the bridge DB and do trigger routing.
- When `HandleMatrixMessage` is not firing, ensure the room is created as a portal for the remote room ID.

## Testing Guidelines
- Use standard Go tests (`*_test.go`, `TestXxx`).
- Run `go test ./...` before committing.

## Commit & Pull Request Guidelines
- Commit messages are short, imperative summaries (e.g., "Add login flow").
- Keep commits focused; avoid mixing refactors with behavior changes.
- PRs should explain what changed and why, and link issues when applicable.

## Security & Configuration Tips
- libolm is required for crypto features; on macOS: `brew install libolm`.
- Avoid committing credentials; use environment variables or local config files.

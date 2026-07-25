# MineGo Protocol 776 Reliability Review

- Review date: 2026-07-26
- Target: Minecraft 26.2, protocol 776
- Branch: `fix/protocol-776-complete`
- Scope: packet decoding, connection lifecycle, callback isolation, game state handling, malformed input behavior, and protocol compatibility
- Classification: general software reliability and protocol correctness review

## Executive Summary

The review found several places where malformed or truncated packet data, callback panics, blocked writes, or invalid game state could terminate or permanently stall a client process. The main fixes establish bounded wire decoders, remove panic-based packet marshaling, isolate user callbacks, make connection shutdown interrupt blocked I/O, and validate allocation sizes before constructing slices or maps.

The normal test suite and static analysis pass. Login, configuration, and game decoder fuzz smoke tests pass. Game fuzzing found and led to a fix for an oversized `PlayerInfoUpdate` map allocation. The remaining worker termination was isolated to the fuzz harness running every registered game decoder for each input; the harness now targets one packet ID per invocation and completes both the five-second acceptance smoke and a 30-second diagnostic run. Race testing is not available in the current Windows environment because `gcc` is not installed.

No commit has been created yet.

## Completed Review Findings

### Packet Codec Panic Handling

Finding:

- The upstream packet marshal path could panic when a field encoder panicked.
- Decode helpers did not consistently reject trailing bytes.
- A panic in a packet field could escape the codec boundary.

Changes:

- Added `pkg/protocol/packet/codec.go`.
- Added safe `Marshal` and `Scan` implementations.
- Added `CodecPanicError` so codec panics become ordinary errors.
- Added nil-field checks and strict trailing-data validation.

Result:

- Packet codec failures now return errors instead of directly terminating the calling process.

### Frame Encoding and Compression

Finding:

- The project depended on behavior supplied through a local `go.mod replace` of `go-mc`.
- A downstream module does not inherit dependency replacements, so consumers could receive different frame behavior.
- Compressed frames needed stricter length and trailing-data validation.

Changes:

- Added `pkg/protocol/packet/frame.go` with local `ReadFrame` and `WriteFrame` implementations.
- Validated frame length, data length, compression threshold, decompressed size, extra compressed bytes, and trailing compressed input.
- Migrated `pkg/auth/auth.go`, `pkg/client/client.go`, and `pkg/client/connect.go` to the local frame codec.
- Removed the fork replacement from `go.mod` and refreshed `go.sum` with `go mod tidy`.

Result:

- Runtime correctness no longer depends on a transitive `replace` directive.

### Bounded Wire Decoding

Finding:

- Generated and manual codecs previously used decoders that could allocate directly from remote lengths.
- Strings, arrays, bit sets, holder sets, byte arrays, and NBT needed common limits.
- The upstream holder-set decoder could record the encoded byte count instead of the decoded ID value.

Changes:

- Added `pkg/protocol/wire/wire.go` as a compatibility layer over `go-mc/net/packet`.
- Added bounded implementations for `String`, `ByteArray`, `Array`, `BitSet`, `IDSet`, and `NBTField`.
- Added UTF-8 validation and UTF-16 character-count validation for protocol strings.
- Limited normal collections to 32,767 entries on decoding and encoding, including specialized generated array helpers.
- Limited NBT payloads to the packet data limit, nesting to 512 levels, and collection entries to bounded values.
- Added safe `Message`, `JsonMessage`, `Property`, and `Session` codecs.
- Updated generated codec imports to use `pkg/protocol/wire` as `packet`.
- Updated `docs/packet-codegen.md` to document the required generated import and schema ownership boundary.

Result:

- Generated packet codecs receive consistent allocation and nesting limits without duplicating guards in every generated method.

### Player Info Allocation

Finding:

- `PlayerInfoUpdate.ReadFrom` used the decoded player count directly as a map capacity.
- A five-byte fuzz input could request a map with tens of millions of entries before the decoder reached EOF.

Changes:

- Reject negative player counts.
- Reject player counts above `wire.MaxCollectionEntries` before map allocation.
- Reject oversized player maps during encoding.
- Added `TestPlayerInfoUpdateRejectsOversizedPlayerCount`.

Result:

- The original oversized allocation input now fails immediately with an ordinary decode error.

### Item and Component Decoding

Finding:

- `TradeSlot.ReadFrom` ignored important I/O errors.
- Unknown component IDs could result in nil dereferences.
- Component patches accepted invalid counts, duplicate IDs, add/remove conflicts, and excessive recursive slot displays.
- The full-profile branch of `ResolvableProfile` recursively decoded another `ResolvableProfile` instead of a `GameProfile`.
- Unknown recipe display discriminants could return success with a nil or stale display.

Changes:

- Updated `pkg/protocol/slot/trade_item.go` to preserve read errors and reject unknown components.
- Updated `pkg/protocol/slot/item_stack.go` and `pkg/protocol/slot/item_stack_template.go` to validate add/remove counts, negative IDs, duplicates, add/remove conflicts, and write failures.
- Added a maximum slot-display recursion depth of 64.
- Corrected `ResolvableProfile` to decode its full branch as `GameProfile` and reject unknown discriminants.
- Made recipe display codecs reject nil values and unknown discriminants while reporting complete byte counts.
- Added malformed, truncated, unknown-component, oversized-count, and recursion tests.

Result:

- Invalid item/component payloads return decode errors instead of reaching nil dereferences, stale values, or unbounded recursion.

### Chunk Decoding

Finding:

- Chunk payload fields could allocate from unchecked lengths.
- Invalid block-state IDs could later index outside `block.StateList`.

Changes:

- Reworked `LevelChunkWithLight.ReadFrom` in `pkg/protocol/packet/game/client/level_chunk_with_light.go`.
- Bounded heightmaps, section bytes, block entities, NBT, light masks, and light arrays.
- Validated decoded block-state IDs before use by the built-in world handler.
- Added malformed chunk length tests.

Result:

- Malformed chunk packets fail at the decoder or state-validation boundary rather than panicking in world updates.

### Client Lifecycle and I/O Cancellation

Finding:

- `Close` could wait behind a blocked write.
- Context cancellation did not necessarily interrupt an in-flight connection write.
- Concurrent `Connect` calls and duplicate `HandleGame` calls needed serialization.

Changes:

- Updated `pkg/client/client.go` to close the connection without waiting for `writeMu`.
- Used `context.AfterFunc` and connection deadlines to interrupt cancelled writes.
- Serialized `Connect` and rejected duplicate game loops.
- Added handling for server disconnect, cookie, resource-pack, and chunk-batch control packets.
- Added blocked-write, cancellation, and callback lifecycle tests.

Result:

- Shutdown and cancellation can interrupt blocked network operations instead of indefinitely waiting for them.

### Callback Isolation

Finding:

- Panics in typed handlers, raw handlers, generic handlers, or event subscribers could escape into library control flow.
- Raw packet handlers shared one mutable payload slice.
- `bot.AddHandler` and `bot.SubscribeEvent` used unchecked type assertions.

Changes:

- Added callback recovery and `CallbackPanicError` handling in `pkg/client/handler.go` and `pkg/client/event.go`.
- Cloned raw packet bytes per handler.
- Removed unchecked generic type assertions in `pkg/bot/handler.go` and `pkg/bot/event.go`.
- Prevented raw handler shutdown from permanently blocking the client loop.

Result:

- User callback failures are isolated and returned through controlled error paths.

### Configuration Flow

Finding:

- Configuration state lacked a complete timeout and response flow for protocol 776 control packets.
- Decoded configuration packets did not consistently require full packet consumption.

Changes:

- Added a 30-second configuration timeout.
- Sent initial `ConfigClientInformation`.
- Added strict full-packet scanning.
- Added cookie, known-pack, resource-pack, code-of-conduct, keepalive, and ping responses.

Result:

- Configuration state now has bounded lifetime and explicit handling for expected protocol 776 control packets.

### Chat Acknowledgement Validation

Finding:

- Signed chat packets require a fixed 20-bit acknowledgement set.
- Invalid bit-set lengths could produce protocol-invalid packets.

Changes:

- Added validation to `Chat` and `ChatCommandSigned`.
- Initialized the timestamp and fixed acknowledgement bit set in `Player.Chat`.

Result:

- Outbound signed chat data is validated before encoding.

### World, Inventory, Pathfinding, and Movement

Finding:

- Large world query radii and coordinate conversions could overflow.
- Container packets could request invalid sizes or indexes.
- Pathfinding accepted invalid coordinates and unbounded node limits.
- Player movement could commit local state before a packet write succeeded.
- `LookAt` could normalize a zero vector and generate NaN rotations.

Changes:

- Added world query radius, volume, and coordinate overflow checks in `pkg/game/world/world.go`.
- Added container size/index checks and nil-client click handling in `pkg/game/inventory/inventory.go`.
- Added nil-world, finite-coordinate, node-count, coordinate-range, and bounded path-search checks in `pkg/game/player/pathfinding.go`.
- Added finite target validation to `FlyTo`, `WalkTo`, and `LookAt`.
- Avoided a write for `LookAt` when the target equals the current position.
- Updated local movement and rotation state only after a successful packet write.
- Added focused world, inventory, pathfinding, and player tests.

Result:

- Invalid local game inputs return errors and failed network writes no longer falsely commit local movement state.

### Authentication and HTTP Handling

Finding:

- HTTP response bodies needed explicit limits and read-error propagation.
- The local OAuth callback server needed bounded startup/shutdown behavior.
- Unrelated OAuth state callbacks could interrupt a valid login flow.
- Yggdrasil profile JSON uses `id`, while the Go field was named `UUID` without a matching tag.

Changes:

- Limited MSA JSON responses to 1 MiB and propagated body read errors.
- Added OAuth server timeouts and bounded shutdown.
- Ignored unrelated OAuth state callbacks while waiting for the valid callback.
- Added `json:"id"` and `json:"name"` tags to `auth.Profile`.
- Added a Yggdrasil selected-profile decode test.

Result:

- Authentication HTTP operations have bounded response handling and profile UUIDs decode from the expected API shape.

### Proxy Address Handling

Finding:

- A SOCKS5 connection using a server address without an explicit port passed the unmodified address to the proxy dialer.

Changes:

- Normalized proxied addresses with the parsed/default port.
- Added support for a bare IP address using the default Minecraft port.

Result:

- SOCKS5 dialing receives a valid host-and-port destination.

## Verification Status

| Verification | Result | Notes |
| --- | --- | --- |
| `go test ./...` | Passed | Passed after the final codec, component-patch, profile, and fuzz-harness changes. |
| `go vet ./...` | Passed | No output. |
| Configuration decoder fuzz smoke | Passed | Five-second run completed with one worker. |
| Login decoder fuzz smoke | Passed | Five-second run completed with one worker. |
| Game decoder fuzz smoke | Passed | Per-packet-ID harness completed the five-second acceptance run and a 30-second diagnostic run. |
| `go test -race ./...` | Not available | The environment has CGO disabled and no `gcc` executable. |
| Official 26.2 server comparison | Not run | The documented server JAR path does not exist in this workspace. |

## Remaining Verification

### Race Verification

Install a Windows-compatible C compiler or run the repository in an environment with CGO support, then execute:

```powershell
$env:CGO_ENABLED = '1'
go test -race ./...
```

### Vanilla Comparison

If an official Minecraft 26.2 server JAR is supplied, compare the remaining custom packet layouts and state transitions against the protocol 776 implementation. The previously documented path was:

```text
C:\Users\miku0139\AppData\Local\Temp\opencode\minecraft-26.2-server\META-INF\versions\26.2\server-26.2.jar
```

That file was not present during this review.

## Main Files Changed

- `pkg/protocol/packet/codec.go`
- `pkg/protocol/packet/frame.go`
- `pkg/protocol/wire/wire.go`
- `pkg/client/client.go`
- `pkg/client/connect.go`
- `pkg/client/handler.go`
- `pkg/client/event.go`
- `pkg/auth/auth.go`
- `pkg/msa/auth.go`
- `pkg/protocol/packet/game/client/level_chunk_with_light.go`
- `pkg/protocol/packet/game/client/player_info_update.go`
- `pkg/protocol/packet/game/client/fuzz_test.go`
- `pkg/protocol/slot/trade_item.go`
- `pkg/protocol/slot/item_stack.go`
- `pkg/protocol/slot/item_stack_template.go`
- `pkg/protocol/slot/display/recipe/recipe_display.go`
- `pkg/protocol/profile.go`
- `pkg/game/world/world.go`
- `pkg/game/inventory/inventory.go`
- `pkg/game/player/pathfinding.go`
- `pkg/game/player/player.go`
- `docs/packet-codegen.md`
- `go.mod`
- `go.sum`

## Current Repository State

- Changes are present in the working tree and are not committed.
- The regular test suite, static analysis, and all three decoder fuzz smoke suites pass at the recorded checkpoints.
- Race verification and an official server comparison remain unavailable in the current workspace.

# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

MinEGO is a Go-based Minecraft protocol implementation that generates packet codecs and data components for Minecraft Java Edition protocol version 1.21.6 (protocol 771). The project uses code generation to create serialization/deserialization code for Minecraft network packets and item components.

## Architecture

The project follows a modular architecture centered around code generation:

- **`codec/data/component/`** - Contains Go structs for all Minecraft item components (e.g., damage, enchantments, custom_name). Each component implements the `slot.Component` interface with Type() and ID() methods.
- **`codec/data/packet/game/`** - Contains client and server packet definitions following the Minecraft protocol specification
- **`codec/data/slot/`** - Core slot/item stack implementation and component registration system
- **`net/`** - Network layer implementation (currently empty)
- **Protocol Reference** - `codec/data/packet/protocol.wiki` contains the complete Minecraft protocol specification from wiki.vg

## Code Generation System

The project uses a custom code generator called `packetizer` that processes Go structs marked with `//codec:gen` comments. This generates the necessary ReadFrom/WriteTo methods for network serialization.

Key patterns:
- Structs marked with `//codec:gen` will have codecs auto-generated
- Structs marked with `//codec:ignore` are excluded from generation
- Components must implement `Type() slot.ComponentID` and `ID() string` methods
- Components are registered in `codec/data/component/components.go` using `slot.RegisterComponent()`

## Development Commands

### Code Generation
```bash
# Generate packet codecs (run from project root)
./gen-packet.sh
# This runs: packetizer ./codec
```

### Standard Go Commands
```bash
# Build the project
go build

# Run tests (if any exist)
go test ./...

# Format code
go fmt ./...

# Get dependencies
go mod tidy
```

## Dependencies

- **github.com/Tnze/go-mc** - Core Minecraft protocol library (uses forked version)
- **packetizer** - Custom code generation tool (available at /root/go/bin/packetizer)

## Working with Components

When adding new Minecraft item components:
1. Create the struct in `codec/data/component/` following existing patterns
2. Add `//codec:gen` comment above the struct
3. Implement `Type()` and `ID()` methods with correct component ID and namespace
4. Register the component in `components.go` init() function
5. Run `./gen-packet.sh` to generate codecs

## Working with Packets

When adding new packet types:
1. Define the packet struct in appropriate `client/` or `server/` directory
2. Add `//codec:gen` comment for auto-generation
3. Register the packet in the appropriate packets map
4. Run `./gen-packet.sh` to generate codecs

## Module Path

The project uses `git.konjactw.dev/patyhank/minego` as its module path and includes a replace directive for the go-mc dependency pointing to a custom fork.

## Code Generation Tasks

- **Java to Go Packet Conversion**: 
  * Task to implement Game State packets by converting Java code from `~/GoProjects/mc-network-source/network/protocol/game` to Go structs
  * Focus on maintaining the same read/write logic during translation
  * Ensure packet structures match the original Java implementation
  * Use packetizer for automatic codec generation

## Reference Notes

- PacketID Reference:
  * In `/root/go/pkg/mod/git.konjactw.dev/patyhank/go-mc@v1.20.3-0.20250618004758-a3d57fde34e8/data/packetid/packetid.go`, you can find all packet IDs (excluding handshake stage)

## Packet Serialization Techniques

- 對於Optional封包資料 簡易的格式可以使用pk.Option[XXX,*XXX] 的方法 比較複雜的需要自己實作一個ReadFrom&WriteTo 並且移除codec:gen的標記

### Example

```go
// ExampleTypeFull full example for codec generator
//
//codec:gen
type ExampleTypeFull struct {
	PlayerID         int32 `mc:"VarInt"`
	PlayerName       string
	UUID             uuid.UUID      `mc:"UUID"`
	ResourceLocation string         `mc:"Identifier"`
	Data             nbt.RawMessage `mc:"NBT"`
	ByteData         []byte         `mc:"ByteArray"`
	Health           float32
	Balance          float64
	Message          chat.Message
	SentMessages     []chat.Message
}
```
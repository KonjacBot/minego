# Packet Codec Generation

This runtime worktree intentionally does not regenerate packet schema files.

Why:

- `gen-packet.sh` previously pointed at a missing `./codec` helper.
- The pinned upstream `github.com/go-mc/packetizer@v0.0.0-20250619063049-ad94ce2fdd81` is not reproducible for this repository as-is.
- That upstream tool hard-codes `github.com/Tnze/go-mc/net/packet` instead of `github.com/KonjacBot/go-mc/net/packet`.
- It also lacks this repository's `//opt:*` handling, so it cannot reproduce the checked-in `codecs.go` files.

Deterministic runtime-worktree policy:

- Treat packet field order and schema logic in checked-in `pkg/protocol/**/codecs.go` files as generated inputs.
- Generated codecs must import `pkg/protocol/wire` as `packet`; that compatibility layer enforces allocation, string, collection, and NBT limits without duplicating guards in every generated method.
- Specialized generated array helpers that do not call `wire.Array` must reject lengths above `wire.MaxCollectionEntries` in both `ReadFrom` and `WriteTo`.
- Do not hand-edit generated field order or optional/enum logic here. Make schema changes in the schema worktree, then validate the generated output with the runtime safety and fuzz tests.

Schema-owner inputs for protocol 776:

- Official server jar reference:
  `C:\Users\miku0139\AppData\Local\Temp\opencode\minecraft-26.2-server\META-INF\versions\26.2\server-26.2.jar`
- Upstream packetizer reference version:
  `github.com/go-mc/packetizer@v0.0.0-20250619063049-ad94ce2fdd81`

If packet schemas must change, do that in the dedicated schema worktree with its generator,
then bring the generated files back into runtime work after review.

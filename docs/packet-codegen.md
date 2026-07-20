# Packet Codec Generation

This runtime worktree intentionally does not regenerate packet schema files.

Why:

- `gen-packet.sh` previously pointed at a missing `./codec` helper.
- The pinned upstream `github.com/go-mc/packetizer@v0.0.0-20250619063049-ad94ce2fdd81` is not reproducible for this repository as-is.
- That upstream tool hard-codes `github.com/Tnze/go-mc/net/packet` instead of `github.com/KonjacBot/go-mc/net/packet`.
- It also lacks this repository's `//opt:*` handling, so it cannot reproduce the checked-in `codecs.go` files.

Deterministic runtime-worktree policy:

- Treat checked-in `pkg/protocol/**/codecs.go` files as read-only inputs in this worktree.
- Do not regenerate or hand-edit packet schema files here.
- Make runtime-only fixes outside schema-owned files, then validate with tests.

Schema-owner inputs for protocol 776:

- Official server jar reference:
  `C:\Users\miku0139\AppData\Local\Temp\opencode\minecraft-26.2-server\META-INF\versions\26.2\server-26.2.jar`
- Upstream packetizer reference version:
  `github.com/go-mc/packetizer@v0.0.0-20250619063049-ad94ce2fdd81`

If packet schemas must change, do that in the dedicated schema worktree with its generator,
then bring the generated files back into runtime work after review.

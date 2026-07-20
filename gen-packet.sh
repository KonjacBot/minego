#!/usr/bin/env sh
set -eu

cat <<'EOF'
Packet schema regeneration is intentionally disabled in this runtime worktree.

See docs/packet-codegen.md for the pinned upstream generator limitations,
the official protocol 776 server jar reference, and the schema-owner workflow.
EOF

exit 1

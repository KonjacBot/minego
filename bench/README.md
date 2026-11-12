# Benchmarks

`BenchmarkHandleGame` (pkg/client) replays 40 ticks of busy play traffic
(entity movement, rotation, velocity, sounds, time; compression threshold 256)
through the real read loop with the built-in world, inventory and player
modules attached.

- `v0-readloop.txt`: the v0 handler/event system, already on the optimized go-mc.
- `v2-readloop.txt`: the v2 bus with lazy decoding.

```bash
go test -run '^$' -bench BenchmarkHandleGame -benchtime 300ms -count 6 ./pkg/client > bench/new.txt
benchstat bench/v2-readloop.txt bench/new.txt
```

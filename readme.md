# minego

go-mc with command-line-only bot client.

# 建議目錄結構

```
minego/
├─ cmd/
│  ├─ minectl/                 # 範例 CLI：連線、發包、抓封包
│  │  └─ main.go
│  └─ proxy/                   # 範例：簡易協議代理/抓包器
│     └─ main.go
│
├─ pkg/                        # 對外公開 API（庫）
│  ├─ client/                  # 高階 Client SDK（使用者只需要這個）
│  │  ├─ client.go             # Client 對外介面、New(...)、Connect(...)
│  │  ├─ options.go            # 可選項：代理、壓縮、密碼學、登入方式
│  │  ├─ session.go            # 與伺服器的一次連線（狀態機）
│  │  ├─ pipeline.go           # 封包處理管線（decode -> route -> handler）
│  │  ├─ dispatcher.go         # 事件/封包分發與訂閱
│  │  ├─ keepalive.go
│  │  ├─ reconnect.go
│  │  └─ errors.go
│  │
│  ├─ transport/               # 可替換傳輸層（TCP wrapper）
│  │  ├─ tcp/
│  │  │  └─ conn.go
│  │  └─ transport.go          # 抽象介面：Dial(ctx, addr) (Conn, error)
│  │
│  ├─ auth/                    # 登入/加密/密鑰交換（mojang、離線、自訂yggdrasil）
│  │  ├─ offline.go
│  │  ├─ mojang.go
│  │  └─ encrypt.go
│  │
│  ├─ handler/                 # 封包與事件處理（基於協議的 client 方向）
│  │  ├─ login.go
│  │  ├─ play_entities.go
│  │  ├─ play_world.go
│  │  ├─ chat.go
│  │  └─ registry.go           # 封包 -> handler 的綁定註冊
│  │
│  ├─ game/                    # 遊戲狀態（抽象，不與 GUI 綁死）
│  │  ├─ world/
│  │  │  ├─ chunk.go
│  │  │  ├─ palette.go
│  │  │  └─ biome.go
│  │  ├─ entity/
│  │  │  └─ entity.go
│  │  └─ inventory/
│  │     └─ slots.go
│  │
│  ├─ data/                    # 協議資料&對照（版本表、映射、assets）
│  │  ├─ versions/
│  │  │  └─ 1_21.json
│  │  └─ registries/
│  │     └─ packets.json
│  │
│  ├─ protocol/                # 你現有的 codec/packet/metadata 可移到這
│  │  ├─ codec/…
│  │  ├─ packet/…
│  │  └─ nbt/…
│  │
│  └─ util/                    # 小工具：varint、zlib、pool、log
│     └─ …
│
└─ go.mod
```

> 原則：
>
> * `pkg/` 對外公開、穩定 API；`internal/` 只給本專案使用。
> * `protocol` 保持「與傳輸無關」；`transport` 抽象連線；`client` 串起狀態機與 handler。
> * `handler` 專心處理「已解碼封包」到「遊戲狀態/事件」的映射。
> * `game` 做資料模型（世界、實體、物品），不要直接依賴 UI。

---

# 模組邏輯切分（設計重點）

1. **狀態機（State Machine）**

    * `Handshake -> Status/Login -> Play -> (Disconnected)`
    * 在 `session.go` 以 goroutine + channel 管理讀寫，使用 `context.Context` 控制生命週期。
2. **封包管線（Pipeline）**

    * `reader` 取得原始 bytes → `protocol/codec` 解碼 → `dispatcher` 依封包 ID 分派到 handler。
    * 可在管線節點插中介：壓縮、加密、記錄、度量。
3. **事件導向 API**

    * 對外提供：

      ```go
      type Client interface {
        On(event Event, fn any) Unsub
        Send(ctx context.Context, p protocol.Packet) error
        State() client.State
      }
      ```
    * `On(PacketPlayChat, func(*ChatMessage){…})` 這種型別安全的註冊可以用泛型或介面實作。
4. **錯誤與可觀測性**

    * 統一 `errors.go`；在 pipeline/handler 位置加可選的 `WithLogger / WithMetrics`。

---

# 命名與風格小訣竅

* 套件名短小、名詞為主：`client`, `auth`, `transport`, `handler`, `game`, `protocol`.
* 檔名以職責分組，不以每個封包獨立檔案（容易爆量難找）。
* 對外只匯出 `pkg/client` 的型別；其餘盡量小寫封裝。
* 盡量用 `context.Context`、`io.Reader/Writer` 介面做邊界。

---

> MCProtocol(新wiki.vg) https://minecraft.wiki/w/Java_Edition_protocol
> 使用套件: github.com/Tnze/go-mc
package client

import (
	"fmt"
	"io"
	"sort"

	"github.com/google/uuid"

	"git.konjactw.dev/falloutBot/go-mc/chat"
	"git.konjactw.dev/falloutBot/go-mc/chat/sign"
	pk "git.konjactw.dev/falloutBot/go-mc/net/packet"
	"git.konjactw.dev/falloutBot/go-mc/yggdrasil/user"
)

type PlayerInfo interface {
	pk.Field
	playerInfoBitMask() int
}

type PlayerInfoUpdate struct {
	Players map[uuid.UUID][]PlayerInfo
}

func (p PlayerInfoUpdate) WriteTo(w io.Writer) (n int64, err error) {
	actions, err := collectPlayerInfoActions(p.Players)
	if err != nil {
		return 0, err
	}

	bitset := pk.NewFixedBitSet(8)
	for _, action := range actions {
		bitset.Set(actionIndex(action), true)
	}
	n1, err := bitset.WriteTo(w)
	if err != nil {
		return n1, err
	}
	n += n1
	n2, err := pk.VarInt(len(p.Players)).WriteTo(w)
	if err != nil {
		return n1 + n2, err
	}
	n += n2
	for playerUUID, infos := range p.Players {
		n3, err := (*pk.UUID)(&playerUUID).WriteTo(w)
		if err != nil {
			return n1 + n2 + n3, err
		}
		n += n3

		infosByAction, err := indexPlayerInfos(infos)
		if err != nil {
			return n, err
		}

		for _, action := range actions {
			info, ok := infosByAction[action]
			if !ok {
				return n, fmt.Errorf("player %s missing player info action mask %#02x", playerUUID.String(), action)
			}
			n4, err := info.WriteTo(w)
			if err != nil {
				return n1 + n2 + n3 + n4, err
			}
			n += n4
		}
	}
	return n, nil
}

func (p *PlayerInfoUpdate) ReadFrom(r io.Reader) (n int64, err error) {
	bitset := pk.NewFixedBitSet(8)
	n1, err := bitset.ReadFrom(r)
	if err != nil {
		return n1, err
	}
	n += n1

	actions := actionsFromBitSet(bitset)

	var playerCount pk.VarInt
	n2, err := playerCount.ReadFrom(r)
	if err != nil {
		return n + n2, err
	}
	n += n2

	players := make(map[uuid.UUID][]PlayerInfo, int(playerCount))
	for i := 0; i < int(playerCount); i++ {
		var playerUUID uuid.UUID
		n3, err := (*pk.UUID)(&playerUUID).ReadFrom(r)
		if err != nil {
			return n + n3, err
		}
		n += n3

		infos := make([]PlayerInfo, 0, len(actions))
		for _, action := range actions {
			info := newPlayerInfoByAction(action)
			if info == nil {
				return n, fmt.Errorf("unsupported player info action mask %#02x", action)
			}
			n4, err := playerInfoRead(&infos, info, r)
			if err != nil {
				return n + n4, err
			}
			n += n4
		}
		players[playerUUID] = infos
	}

	p.Players = players
	return n, nil
}

func playerInfoRead(infos *[]PlayerInfo, info PlayerInfo, r io.Reader) (int64, error) {
	n, err := info.ReadFrom(r)
	if err != nil {
		return n, err
	}
	*infos = append(*infos, info)
	return n, err
}

func collectPlayerInfoActions(players map[uuid.UUID][]PlayerInfo) ([]int, error) {
	actions := make(map[int]struct{}, 8)
	for playerID, infos := range players {
		seen := make(map[int]struct{}, len(infos))
		for _, info := range infos {
			mask := info.playerInfoBitMask()
			if err := validatePlayerInfoActionMask(mask); err != nil {
				return nil, fmt.Errorf("player %s has invalid action: %w", playerID.String(), err)
			}
			if _, exists := seen[mask]; exists {
				return nil, fmt.Errorf("player %s has duplicated action mask %#02x", playerID.String(), mask)
			}
			seen[mask] = struct{}{}
			actions[mask] = struct{}{}
		}
	}

	sorted := make([]int, 0, len(actions))
	for mask := range actions {
		sorted = append(sorted, mask)
	}
	sort.Ints(sorted)
	return sorted, nil
}

func indexPlayerInfos(infos []PlayerInfo) (map[int]PlayerInfo, error) {
	indexed := make(map[int]PlayerInfo, len(infos))
	for _, info := range infos {
		mask := info.playerInfoBitMask()
		if err := validatePlayerInfoActionMask(mask); err != nil {
			return nil, err
		}
		if _, exists := indexed[mask]; exists {
			return nil, fmt.Errorf("duplicated player info action mask %#02x", mask)
		}
		indexed[mask] = info
	}
	return indexed, nil
}

func validatePlayerInfoActionMask(mask int) error {
	if mask <= 0 || mask > 0x80 {
		return fmt.Errorf("action mask out of range: %#02x", mask)
	}
	// Action mask is an enum bit. It must have exactly one bit set.
	if mask&(mask-1) != 0 {
		return fmt.Errorf("action mask must be a single bit: %#02x", mask)
	}
	return nil
}

func actionIndex(mask int) int {
	index := 0
	for (mask & 0x01) == 0 {
		index++
		mask >>= 1
	}
	return index
}

func actionsFromBitSet(bitset pk.FixedBitSet) []int {
	actions := make([]int, 0, len(playerInfoActionMasks))
	for _, mask := range playerInfoActionMasks {
		if bitset.Get(actionIndex(mask)) {
			actions = append(actions, mask)
		}
	}
	return actions
}

func newPlayerInfoByAction(mask int) PlayerInfo {
	switch mask {
	case 0x01:
		return &PlayerInfoAddPlayer{}
	case 0x02:
		return &PlayerInfoInitializeChat{}
	case 0x04:
		return &PlayerInfoUpdateGameMode{}
	case 0x08:
		return &PlayerInfoUpdateListed{}
	case 0x10:
		return &PlayerInfoUpdateLatency{}
	case 0x20:
		return &PlayerInfoUpdateDisplayName{}
	case 0x40:
		return &PlayerInfoUpdateListPriority{}
	case 0x80:
		return &PlayerInfoUpdateHat{}
	default:
		return nil
	}
}

var playerInfoActionMasks = [...]int{
	0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80,
}

//codec:gen
type PlayerInfoAddPlayer struct {
	Name       string
	Properties []user.Property
}

//codec:gen
type PlayerInfoChatData struct {
	ChatSessionID uuid.UUID `mc:"UUID"`
	Session       sign.Session
}

//codec:gen
type PlayerInfoInitializeChat struct {
	Data pk.Option[PlayerInfoChatData, *PlayerInfoChatData]
}

//codec:gen
type PlayerInfoUpdateGameMode struct {
	GameMode int32 `mc:"VarInt"`
}

//codec:gen
type PlayerInfoUpdateListed struct {
	Listed bool
}

//codec:gen
type PlayerInfoUpdateLatency struct {
	Ping int32 `mc:"VarInt"`
}

//codec:gen
type PlayerInfoUpdateDisplayName struct {
	DisplayName pk.Option[chat.Message, *chat.Message]
}

//codec:gen
type PlayerInfoUpdateListPriority struct {
	Priority int32 `mc:"VarInt"`
}

//codec:gen
type PlayerInfoUpdateHat struct {
	Visible bool
}

func (PlayerInfoAddPlayer) playerInfoBitMask() int {
	return 0x01
}

func (PlayerInfoInitializeChat) playerInfoBitMask() int {
	return 0x02
}

func (PlayerInfoUpdateGameMode) playerInfoBitMask() int {
	return 0x04
}

func (PlayerInfoUpdateListed) playerInfoBitMask() int {
	return 0x08
}

func (PlayerInfoUpdateLatency) playerInfoBitMask() int {
	return 0x10
}

func (PlayerInfoUpdateDisplayName) playerInfoBitMask() int {
	return 0x20
}

func (PlayerInfoUpdateListPriority) playerInfoBitMask() int {
	return 0x40
}

func (PlayerInfoUpdateHat) playerInfoBitMask() int {
	return 0x80
}

package client

import (
	"io"

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
	//bitset := pk.NewFixedBitSet(8)
	//for _, infos := range p.Players {
	//	for _, info := range infos {
	//		bitset.Set(info.playerInfoBitMask(), true)
	//	}
	//}
	//n1, err := bitset.WriteTo(w)
	//if err != nil {
	//	return n1, err
	//}
	//n += n1
	//n2, err := pk.VarInt(len(p.Players)).WriteTo(w)
	//if err != nil {
	//	return n1 + n2, err
	//}
	//n += n2
	//for playerUUID, infos := range p.Players {
	//	n3, err := (*pk.UUID)(&playerUUID).WriteTo(w)
	//	if err != nil {
	//		return n1 + n2 + n3, err
	//	}
	//	n += n3
	//	for _, info := range infos {
	//		n4, err := info.WriteTo(w)
	//		if err != nil {
	//			return n1 + n2 + n3 + n4, err
	//		}
	//		n += n4
	//	}
	//}
	return 0, nil
}

func (p *PlayerInfoUpdate) ReadFrom(r io.Reader) (n int64, err error) {
	//bitset := pk.NewFixedBitSet(256)
	//n1, err := bitset.ReadFrom(r)
	//if err != nil {
	//	return n1, err
	//}
	//m := make(map[uuid.UUID][]PlayerInfo)
	//
	//var playerLens pk.VarInt
	//n2, err := playerLens.ReadFrom(r)
	//if err != nil {
	//	return n1 + n2, err
	//}
	//for i := 0; i < int(playerLens); i++ {
	//	var playerUUID uuid.UUID
	//	n3, err := (*pk.UUID)(&playerUUID).ReadFrom(r)
	//	if err != nil {
	//		return n1 + n2 + n3, err
	//	}
	//	var temp int64
	//	var infos []PlayerInfo
	//	if bitset.Get(0x01) {
	//		n4, err := playerInfoRead(&infos, &PlayerInfoAddPlayer{}, r)
	//		if err != nil {
	//			return n1 + n2 + n3 + n4, err
	//		}
	//		temp += n4
	//	}
	//	if bitset.Get(0x02) {
	//		n4, err := playerInfoRead(&infos, &PlayerInfoInitializeChat{}, r)
	//		if err != nil {
	//			return n1 + n2 + n3 + n4, err
	//		}
	//		temp += n4
	//	}
	//	if bitset.Get(0x04) {
	//		n4, err := playerInfoRead(&infos, &PlayerInfoUpdateGameMode{}, r)
	//		if err != nil {
	//			return n1 + n2 + n3 + n4, err
	//		}
	//		temp += n4
	//	}
	//	if bitset.Get(0x08) {
	//		n4, err := playerInfoRead(&infos, &PlayerInfoUpdateListed{}, r)
	//		if err != nil {
	//			return n1 + n2 + n3 + n4, err
	//		}
	//		temp += n4
	//	}
	//	if bitset.Get(0x10) {
	//		n4, err := playerInfoRead(&infos, &PlayerInfoUpdateLatency{}, r)
	//		if err != nil {
	//			return n1 + n2 + n3 + n4, err
	//		}
	//		temp += n4
	//	}
	//	if bitset.Get(0x20) {
	//		n4, err := playerInfoRead(&infos, &PlayerInfoUpdateDisplayName{}, r)
	//		if err != nil {
	//			return n1 + n2 + n3 + n4, err
	//		}
	//		temp += n4
	//	}
	//	if bitset.Get(0x40) {
	//		n4, err := playerInfoRead(&infos, &PlayerInfoUpdateListPriority{}, r)
	//		if err != nil {
	//			return n1 + n2 + n3 + n4, err
	//		}
	//		temp += n4
	//	}
	//	if bitset.Get(0x80) {
	//		n4, err := playerInfoRead(&infos, &PlayerInfoUpdateHat{}, r)
	//		if err != nil {
	//			return n1 + n2 + n3 + n4, err
	//		}
	//		temp += n4
	//	}
	//
	//	m[playerUUID] = infos
	//}
	//return
	return 0, nil
}

func playerInfoRead(infos *[]PlayerInfo, info PlayerInfo, r io.Reader) (int64, error) {
	n, err := info.ReadFrom(r)
	if err != nil {
		return n, err
	}
	*infos = append(*infos, info)
	return n, err
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

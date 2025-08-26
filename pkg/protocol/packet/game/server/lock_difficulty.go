package server

import "git.konjactw.dev/falloutBot/go-mc/data/packetid"

//codec:gen
type LockDifficulty struct {
	Locked bool
}

func (*LockDifficulty) PacketID() packetid.ServerboundPacketID {
	return packetid.ServerboundLockDifficulty
}

func init() {
	registerPacket(packetid.ServerboundLockDifficulty, func() ServerboundPacket {
		return &LockDifficulty{}
	})
}

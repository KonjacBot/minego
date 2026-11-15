package keeper

import (
	"bytes"
	"slices"

	"github.com/KonjacBot/go-mc/data/packetid"
	"github.com/KonjacBot/go-mc/level/item"
	pk "github.com/KonjacBot/go-mc/net/packet"

	"github.com/KonjacBot/minego/v2/pkg/protocol/packet/game/client"
	"github.com/KonjacBot/minego/v2/pkg/protocol/packet/game/server"
	"github.com/KonjacBot/minego/v2/pkg/protocol/slot"
)

// playerSlots is the size of the player's own inventory menu (window 0):
// crafting result and grid, armor, main inventory, hotbar, offhand.
const playerSlots = 46

// emptySlot is an empty item stack on the wire (count 0).
var emptySlot = []byte{0}

// inventory follows the player's inventory and the open container. Every slot
// is the server's encoding of that stack (see raw.go).
//
// The server confirms clicks without echoing the result (the client predicts
// it), so the keeper also applies the core's clicks: they carry the changed
// slots as hashes, which are resolved against the stacks the keeper already
// knows (an item can only move between slots it can see).
type inventory struct {
	player  [playerSlots][]byte
	stateID int32
	cursor  []byte
	open    *screen
	// unresolved counts click slots that matched no known stack and were kept
	// as bare item id and count.
	unresolved int
}

// screen is an open container. Its last 36 slots are the player's main
// inventory and hotbar, shared with window 0.
type screen struct {
	id       int32
	open     []byte
	own      [][]byte
	stateID  int32
	data     map[int16]int16
	dataSeq  []int16
	merchant []byte
	content  bool
}

func newInventory() inventory {
	inv := inventory{cursor: emptySlot}
	for i := range inv.player {
		inv.player[i] = emptySlot
	}
	return inv
}

// view returns the slot behind index in window id, or nil if unknown.
func (inv *inventory) view(id int32, index int) *[]byte {
	switch {
	case id == 0:
		if index >= 0 && index < playerSlots {
			return &inv.player[index]
		}
	case inv.open != nil && id == inv.open.id:
		if index >= 0 && index < len(inv.open.own) {
			return &inv.open.own[index]
		}
		if k := index - len(inv.open.own); k >= 0 && k < 36 {
			return &inv.player[9+k]
		}
	}
	return nil
}

func (inv *inventory) fromServer(id packetid.ClientboundPacketID, p pk.Packet) (handled bool, err error) {
	r := bytes.NewReader(p.Data)
	switch id {
	case packetid.ClientboundContainerSetContent:
		var window, state, count pk.VarInt
		for _, f := range []pk.FieldDecoder{&window, &state, &count} {
			if _, err := f.ReadFrom(r); err != nil {
				return true, err
			}
		}
		slots := make([][]byte, 0, int(count))
		for range int(count) {
			raw, err := span(p.Data, r, &slot.Slot{})
			if err != nil {
				return true, err
			}
			slots = append(slots, raw)
		}
		carried, err := span(p.Data, r, &slot.Slot{})
		if err != nil {
			return true, err
		}
		inv.cursor = carried
		switch {
		case window == 0:
			copy(inv.player[:], slots)
			inv.stateID = int32(state)
		case inv.open != nil && int32(window) == inv.open.id:
			n := max(len(slots)-36, 0)
			inv.open.own = slices.Clone(slots[:n])
			copy(inv.player[9:45], slots[n:])
			inv.open.stateID, inv.open.content = int32(state), true
		}
	case packetid.ClientboundContainerSetSlot:
		var window, state pk.VarInt
		var index pk.Short
		for _, f := range []pk.FieldDecoder{&window, &state, &index} {
			if _, err := f.ReadFrom(r); err != nil {
				return true, err
			}
		}
		raw, err := span(p.Data, r, &slot.Slot{})
		if err != nil {
			return true, err
		}
		if window == -1 && index == -1 {
			inv.cursor = raw
			return true, nil
		}
		if s := inv.view(int32(window), int(index)); s != nil {
			*s = raw
		}
		if window == 0 {
			inv.stateID = int32(state)
		} else if inv.open != nil && int32(window) == inv.open.id {
			inv.open.stateID = int32(state)
		}
	case packetid.ClientboundSetCursorItem:
		raw, err := span(p.Data, r, &slot.Slot{})
		if err != nil {
			return true, err
		}
		inv.cursor = raw
	case packetid.ClientboundSetPlayerInventory:
		var index pk.VarInt
		if _, err := index.ReadFrom(r); err != nil {
			return true, err
		}
		raw, err := span(p.Data, r, &slot.Slot{})
		if err != nil {
			return true, err
		}
		if i, ok := menuSlot(int32(index)); ok {
			inv.player[i] = raw
		}
	case packetid.ClientboundOpenScreen:
		inv.open = &screen{id: leadingVarInt(p.Data), open: p.Data, data: map[int16]int16{}}
	case packetid.ClientboundMountScreenOpen:
		// Horse inventories are opened by the server too, but their layout depends
		// on the animal; the keeper does not replay them.
		inv.open = nil
	case packetid.ClientboundContainerClose:
		inv.open = nil
	case packetid.ClientboundContainerSetData:
		var c client.ContainerSetData
		if err := decode(p, &c); err != nil {
			return true, err
		}
		if inv.open != nil && int32(c.ContainerID) == inv.open.id {
			if _, seen := inv.open.data[c.ID]; !seen {
				inv.open.dataSeq = append(inv.open.dataSeq, c.ID)
			}
			inv.open.data[c.ID] = c.Value
		}
	case packetid.ClientboundMerchantOffers:
		if inv.open != nil && leadingVarInt(p.Data) == inv.open.id {
			inv.open.merchant = p.Data
		}
	default:
		return false, nil
	}
	return true, nil
}

// menuSlot maps an index of the player's Inventory (hotbar 0-8, main 9-35,
// armor 36-39 from feet, offhand 40) to its slot in window 0.
func menuSlot(i int32) (int, bool) {
	switch {
	case i >= 0 && i < 9:
		return 36 + int(i), true
	case i >= 9 && i < 36:
		return int(i), true
	case i >= 36 && i < 40:
		return 8 - int(i-36), true
	case i == 40:
		return 45, true
	}
	return 0, false
}

func (inv *inventory) fromCore(id packetid.ServerboundPacketID, p pk.Packet) error {
	switch id {
	case packetid.ServerboundContainerClose:
		var c server.ContainerClose
		if err := decode(p, &c); err != nil {
			return err
		}
		if inv.open != nil && c.WindowID == inv.open.id {
			inv.open = nil
		}
	case packetid.ServerboundContainerClick:
		var c server.ContainerClick
		if err := decode(p, &c); err != nil {
			return err
		}
		inv.click(c)
	}
	return nil
}

// click applies the core's prediction of a click. Every stack the keeper can
// see is a candidate source for the hashed slots; they are hashed before any
// slot changes so a stack that moves can still be found.
func (inv *inventory) click(c server.ContainerClick) {
	if c.WindowID != 0 && (inv.open == nil || c.WindowID != inv.open.id) {
		return
	}
	known := inv.candidates()
	for _, changed := range c.ChangedSlots {
		if s := inv.view(c.WindowID, int(changed.Slot)); s != nil {
			*s = inv.resolve(known, changed.SlotData)
		}
	}
	inv.cursor = inv.resolve(known, c.CarriedSlot)
}

type candidate struct {
	raw  []byte
	hash slot.HashedSlot
}

func (inv *inventory) candidates() []candidate {
	var out []candidate
	add := func(raw []byte) {
		var s slot.Slot
		if _, err := s.ReadFrom(bytes.NewReader(raw)); err != nil || s.Count <= 0 {
			return
		}
		h, err := slot.HashSlot(s)
		if err != nil {
			return
		}
		out = append(out, candidate{raw: raw, hash: h})
	}
	add(inv.cursor)
	for _, s := range inv.player {
		add(s)
	}
	if inv.open != nil {
		for _, s := range inv.open.own {
			add(s)
		}
	}
	return out
}

func (inv *inventory) resolve(known []candidate, want slot.HashedSlot) []byte {
	if !want.HasItem || want.ItemCount <= 0 {
		return emptySlot
	}
	for _, c := range known {
		if sameItem(c.hash, want) {
			return withCount(c.raw, want.ItemCount)
		}
	}
	inv.unresolved++
	var buf bytes.Buffer
	_, _ = slot.Slot{ItemID: item.ID(want.ItemID), Count: want.ItemCount}.WriteTo(&buf)
	return buf.Bytes()
}

// withCount replaces the leading count of an encoded stack.
func withCount(raw []byte, count int32) []byte {
	r := bytes.NewReader(raw)
	var old pk.VarInt
	if _, err := old.ReadFrom(r); err != nil {
		return raw
	}
	return append(varint(count), raw[len(raw)-r.Len():]...)
}

func sameItem(a, b slot.HashedSlot) bool {
	return a.ItemID == b.ItemID &&
		slices.Equal(a.AddComponents, b.AddComponents) &&
		slices.Equal(a.RemovedComponents, b.RemovedComponents)
}

// content encodes a SetContainerContent from raw slots.
func content(window, state int32, slots [][]byte, cursor []byte) pk.Packet {
	data := append(varint(window), varint(state)...)
	data = append(data, varint(int32(len(slots)))...)
	for _, s := range slots {
		data = append(data, s...)
	}
	data = append(data, cursor...)
	return pk.Packet{ID: int32(packetid.ClientboundContainerSetContent), Data: data}
}

// packets returns window 0 and, if one is open, the container with its data
// and trades.
func (inv *inventory) packets() []pk.Packet {
	out := []pk.Packet{content(0, inv.stateID, inv.player[:], inv.cursor)}
	if o := inv.open; o != nil {
		out = append(out, pk.Packet{ID: int32(packetid.ClientboundOpenScreen), Data: o.open})
		if o.content {
			out = append(out, content(o.id, o.stateID, append(slices.Clone(o.own), inv.player[9:45]...), inv.cursor))
		}
		for _, id := range o.dataSeq {
			var buf bytes.Buffer
			_, _ = client.ContainerSetData{ContainerID: int8(o.id), ID: id, Value: o.data[id]}.WriteTo(&buf)
			out = append(out, pk.Packet{ID: int32(packetid.ClientboundContainerSetData), Data: buf.Bytes()})
		}
		if o.merchant != nil {
			out = append(out, pk.Packet{ID: int32(packetid.ClientboundMerchantOffers), Data: o.merchant})
		}
	}
	return out
}

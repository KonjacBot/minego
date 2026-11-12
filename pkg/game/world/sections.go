package world

import (
	"sync"

	"github.com/KonjacBot/go-mc/level"
)

// Sections that hold a single value (all air above the terrain, all stone
// deep down, one biome) carry no per-block data. Every column of every World
// in the process points such sections at one shared, read-only container per
// value instead of keeping its own; World gives a section a private container
// before writing to it (writableSection).
var (
	sharedStates sync.Map // level.BlocksState -> *level.PaletteContainer[level.BlocksState]
	sharedBiomes sync.Map // level.BiomesState -> *level.PaletteContainer[level.BiomesState]
)

func sharedStatesOf(v level.BlocksState) *level.PaletteContainer[level.BlocksState] {
	if c, ok := sharedStates.Load(v); ok {
		return c.(*level.PaletteContainer[level.BlocksState])
	}
	c, _ := sharedStates.LoadOrStore(v, level.NewStatesPaletteContainer(16*16*16, v))
	return c.(*level.PaletteContainer[level.BlocksState])
}

func sharedBiomesOf(v level.BiomesState) *level.PaletteContainer[level.BiomesState] {
	if c, ok := sharedBiomes.Load(v); ok {
		return c.(*level.PaletteContainer[level.BiomesState])
	}
	c, _ := sharedBiomes.LoadOrStore(v, level.NewBiomesPaletteContainer(4*4*4, v))
	return c.(*level.PaletteContainer[level.BiomesState])
}

// shareSingleValueSections points single-value containers at the shared ones.
func shareSingleValueSections(sections []level.Section) {
	for i := range sections {
		s := &sections[i]
		if s.States != nil {
			if v, ok := s.States.SingleValue(); ok {
				s.States = sharedStatesOf(v)
			}
		}
		if s.Biomes != nil {
			if v, ok := s.Biomes.SingleValue(); ok {
				s.Biomes = sharedBiomesOf(v)
			}
		}
	}
}

// writableSection returns a section of chunk ready for SetBlock: a section
// still on a single-value container (possibly shared) gets a private copy
// first. World never writes biomes, so those stay shared.
func writableSection(chunk *level.Chunk, sectionY int32) *level.Section {
	s := &chunk.Sections[sectionY]
	if v, ok := s.States.SingleValue(); ok {
		s.States = level.NewStatesPaletteContainer(16*16*16, v)
	}
	return s
}

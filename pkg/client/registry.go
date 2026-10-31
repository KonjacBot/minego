package client

// RegistryEntry returns the name of entry id in a registry the server sent during configuration, for example
// registry "minecraft:damage_type" with a DamageEvent's SourceType gives "minecraft:mace_smash".
// It reports false before configuration finished and for unknown registries or ids.
func (b *botClient) RegistryEntry(registry string, id int32) (string, bool) {
	b.registryMu.RLock()
	defer b.registryMu.RUnlock()
	name, ok := b.registryEntries[registry][id]
	return name, ok
}

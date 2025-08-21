package component

import (
	"git.konjactw.dev/patyhank/minego/pkg/protocol/slot"
)

func init() {
	slot.RegisterComponent(func() slot.Component {
		return &CustomData{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &MaxStackSize{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &MaxDamage{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &Damage{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &Unbreakable{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &CustomName{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &ItemName{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &ItemModel{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &Lore{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &Rarity{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &Enchantments{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &CanPlaceOn{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &CanBreak{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &AttributeModifiers{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &CustomModelData{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &TooltipDisplay{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &RepairCost{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &CreativeSlotLock{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &EnchantmentGlintOverride{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &IntangibleProjectile{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &Food{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &Consumable{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &UseRemainder{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &UseCooldown{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &DamageResistant{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &Tool{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &Weapon{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &Enchantable{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &Equippable{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &Repairable{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &Glider{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &TooltipStyle{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &DeathProtection{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &BlocksAttacks{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &StoredEnchantments{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &DyedColor{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &MapColor{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &MapID{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &MapDecorations{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &MapPostProcessing{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &ChargedProjectiles{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &BundleContents{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &PotionContents{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &PotionDurationScale{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &SuspiciousStewEffects{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &WritableBookContent{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &WrittenBookContent{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &Trim{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &DebugStickState{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &EntityData{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &BucketEntityData{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &BlockEntityData{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &Instrument{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &ProvidesTrimMaterial{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &OminousBottleAmplifier{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &JukeboxPlayable{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &ProvidesBannerPatterns{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &Recipes{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &LodestoneTracker{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &FireworkExplosion{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &Fireworks{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &Profile{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &NoteBlockSound{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &BannerPatterns{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &BaseColor{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &PotDecorations{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &Container{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &BlockState{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &Bees{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &Lock{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &ContainerLoot{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &BreakSound{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &VillagerVariant{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &WolfVariant{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &WolfSoundVariant{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &WolfCollar{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &FoxVariant{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &SalmonSize{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &ParrotVariant{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &TropicalFishPattern{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &TropicalFishBaseColor{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &TropicalFishPatternColor{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &MooshroomVariant{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &RabbitVariant{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &PigVariant{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &CowVariant{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &ChickenVariant{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &FrogVariant{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &HorseVariant{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &PaintingVariant{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &LlamaVariant{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &AxolotlVariant{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &CatVariant{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &CatCollar{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &SheepColor{}
	})
	slot.RegisterComponent(func() slot.Component {
		return &ShulkerColor{}
	})
}

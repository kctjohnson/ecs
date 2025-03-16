package templates

import "ecs/internal/game/components"

type Rarity string

const (
	Common    Rarity = "common"
	Uncommon  Rarity = "uncommon"
	Rare      Rarity = "rare"
	Epic      Rarity = "epic"
	Legendary Rarity = "legendary"
)

type AttributeRange struct {
	Base     int `json:"base"`
	Variance int `json:"variance"`
}

type LootTableEntry struct {
	ItemID      string  `json:"item_id"`
	Probability float64 `json:"probability"`
}

type EnemyTemplate struct {
	ID         string                    `json:"id"`
	Name       string                    `json:"name"`
	Sprite     rune                      `json:"sprite"`
	Rarity     Rarity                    `json:"rarity"`
	Attributes map[string]AttributeRange `json:"attributes"`
	LootTable  []LootTableEntry          `json:"loot_table"`
}

type ItemType string

const (
	Weapon     ItemType = "weapon"
	Armor      ItemType = "armor"
	Consumable ItemType = "consumable"
)

type ItemTemplate struct {
	ID         string                     `json:"id"`
	Name       string                     `json:"name"`
	Type       ItemType                   `json:"type"`
	Rarity     Rarity                     `json:"rarity"`
	Sprite     rune                       `json:"sprite"`
	Attributes map[string]any             `json:"attributes"`
	Slots      []components.EquipmentSlot `json:"slots,omitempty"`
}

type TemplateManager struct {
	EnemyTemplates map[string]EnemyTemplate
	ItemTemplates  map[string]ItemTemplate
}

func NewTemplateManager() *TemplateManager {
	return &TemplateManager{
		EnemyTemplates: make(map[string]EnemyTemplate),
		ItemTemplates:  make(map[string]ItemTemplate),
	}
}

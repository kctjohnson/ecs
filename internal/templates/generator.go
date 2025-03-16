package templates

import (
	"math/rand"
	"time"
)

// Generator handles the creation of entities from templates
type Generator struct {
	rng *rand.Rand
}

// NewGenerator creates a new Generator with a seeded random number generator
func NewGenerator() *Generator {
	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	rng := rand.New(source)
	return &Generator{rng: rng}
}

// GenerateValueFromRange generates a random value within the given range
func (g *Generator) GenerateValueFromRange(attrRange AttributeRange) int {
	minValue := attrRange.Base - attrRange.Variance
	maxValue := attrRange.Base + attrRange.Variance

	if minValue < 1 {
		minValue = 1
	}

	return minValue + g.rng.Intn(maxValue-minValue+1)
}

// GenerateFloatValue generates a random float value
func (g *Generator) GenerateFloatValue(min, max float64) float64 {
	return min + g.rng.Float64()*(max-min)
}

// DetermineDrops determines which items from the loot table should be dropped
func (g *Generator) DetermineDrops(lootTable []LootTableEntry) []string {
	droppedItems := []string{}

	// Check each item in the loot table
	for _, entry := range lootTable {
		// Generate a random number between 0 and 1
		roll := g.rng.Float64()

		// If the roll is less than the probability, drop the item
		if roll <= entry.Probability {
			droppedItems = append(droppedItems, entry.ItemID)
		}
	}

	return droppedItems
}

package ui

import (
	"ecs/internal/game/components"
	"ecs/pkg/ecs"
)

var equipmentSlotDisplayOrder = []components.EquipmentSlot{
	components.Head,
	components.Torso,
	components.Legs,
	components.Feet,
	components.RightHand,
	components.LeftHand,
}

var equipmentValueToLabel = map[components.EquipmentSlot]string{
	components.Head:      "Head",
	components.Torso:     "Torso",
	components.Legs:      "Legs",
	components.Feet:      "Feet",
	components.RightHand: "Right Hand",
	components.LeftHand:  "Left Hand",
}

type displayEquipmentSlot struct {
	Label string
	Slot  components.EquipmentSlot
	Item  ecs.Entity
}

func makeOreredEquipmentSlice(
	equipment map[components.EquipmentSlot]ecs.Entity,
) []displayEquipmentSlot {
	var ordered []displayEquipmentSlot
	for _, slot := range equipmentSlotDisplayOrder {
		ordered = append(ordered, displayEquipmentSlot{
			Label: equipmentValueToLabel[slot],
			Slot:  slot,
			Item:  equipment[slot],
		})
	}
	return ordered
}

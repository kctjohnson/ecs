package events

import "ecs/pkg/ecs"

const (
	EntityMoved    ecs.EventType = "entity_moved"
	EntityAttacked ecs.EventType = "entity_attacked"
	EntityDefeated ecs.EventType = "entity_defeated"
	HealthChanged  ecs.EventType = "health_changed"
	TurnEnded      ecs.EventType = "turn_ended"
	ItemPickedUp   ecs.EventType = "item_picked_up"
	ItemUsed       ecs.EventType = "item_used"
)

// TODO: We need a way for events to be a bit more typed and have a consistent structure?
type EntityMovedEventData struct {
}

type EntityAttackedEventData struct {
}

type EntityDefeatedEventData struct {
}

type HealthChangedEventData struct {
}

type TurnEndedEventData struct {
}

type ItemPickedUpEventData struct {
}

type ItemUsedEventData struct {
	Item   ecs.Entity
	Target ecs.Entity
}

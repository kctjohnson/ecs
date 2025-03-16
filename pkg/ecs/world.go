package ecs

import "log"

// World is the main ECS container that holds all entities, components, and systems
type World struct {
	EntityManager    *EntityManager
	ComponentManager *ComponentManager
	systems          []System
	eventQueue       []Event // Simple event queue for communication
	eventHandlers    map[EventType][]func(Event)
	Logger           *log.Logger
}

func NewWorld(logger *log.Logger) *World {
	return &World{
		EntityManager:    NewEntityManager(),
		ComponentManager: NewComponentManager(),
		systems:          []System{},
		eventQueue:       []Event{},
		eventHandlers:    make(map[EventType][]func(Event)),
		Logger:           logger,
	}
}

func (w *World) AddSystem(system System) {
	w.systems = append(w.systems, system)
}

func (w *World) RemoveEntity(entity Entity) {
	w.EntityManager.RemoveEntity(entity)
	w.ComponentManager.RemoveAllComponents(entity)
}

func (w *World) Update() {
	for _, system := range w.systems {
		system.Update(w)
	}

	// Process events after all systems have updated
	w.processEvents()
}

// Simple event system for communication between ECS and external systems
type EventType string

type Event struct {
	Type   EventType
	Entity Entity
	Data   map[string]any
}

func (w *World) RegisterEventHandler(eventType EventType, handler func(Event)) {
	w.eventHandlers[eventType] = append(w.eventHandlers[eventType], handler)
}

func (w *World) QueueEvent(eventType EventType, entity Entity, data map[string]any) {
	w.eventQueue = append(w.eventQueue, Event{
		Type:   eventType,
		Entity: entity,
		Data:   data,
	})
}

func (w *World) processEvents() {
	// Process all events in the queue
	for _, event := range w.eventQueue {
		if handlers, exists := w.eventHandlers[event.Type]; exists {
			for _, handler := range handlers {
				handler(event)
			}
		}
	}

	// Clear queue
	w.eventQueue = w.eventQueue[:0]
}

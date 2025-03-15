package ecs

// Component is a marker interface for all component types
type Component any

type ComponentType string

// ComponentManager handles storage and retrieval of components
type ComponentManager struct {
	components map[ComponentType]map[Entity]Component
}

func NewComponentManager() *ComponentManager {
	return &ComponentManager{
		components: make(map[ComponentType]map[Entity]Component),
	}
}

func (cm *ComponentManager) RegisterComponentType(componentType ComponentType) {
	if _, exists := cm.components[componentType]; !exists {
		cm.components[componentType] = make(map[Entity]Component)
	}
}

func (cm *ComponentManager) AddComponent(
	entity Entity,
	componentType ComponentType,
	component Component,
) {
	if _, exists := cm.components[componentType]; !exists {
		cm.RegisterComponentType(componentType)
	}
	cm.components[componentType][entity] = component
}

func (cm *ComponentManager) RemoveComponent(entity Entity, componentType ComponentType) {
	if componentMap, exists := cm.components[componentType]; exists {
		delete(componentMap, entity)
	}
}

func (cm *ComponentManager) GetComponent(
	entity Entity,
	componentType ComponentType,
) (Component, bool) {
	if componentMap, exists := cm.components[componentType]; exists {
		component, found := componentMap[entity]
		return component, found
	}
	return nil, false
}

func (cm *ComponentManager) HasComponent(entity Entity, componentType ComponentType) bool {
	if componentMap, exists := cm.components[componentType]; exists {
		_, found := componentMap[entity]
		return found
	}
	return false
}

func (cm *ComponentManager) GetAllEntitiesWithComponent(componentType ComponentType) []Entity {
	if componentMap, exists := cm.components[componentType]; exists {
		entities := make([]Entity, 0, len(componentMap))
		for e := range componentMap {
			entities = append(entities, e)
		}
		return entities
	}
	return []Entity{}
}

func (cm *ComponentManager) RemoveAllComponents(entity Entity) {
	for _, componentMap := range cm.components {
		delete(componentMap, entity)
	}
}

package ecs

// Entity is just an identifier for game objects
type Entity int

// EntityManager handles entity creation and removal
type EntityManager struct {
	nextEntityID Entity
	entities     map[Entity]struct{}
}

func NewEntityManager() *EntityManager {
	return &EntityManager{
		nextEntityID: 1,
		entities:     make(map[Entity]struct{}),
	}
}

func (em *EntityManager) CreateEntity() Entity {
	e := em.nextEntityID
	em.entities[e] = struct{}{}
	em.nextEntityID++
	return e
}

func (em *EntityManager) RemoveEntity(entity Entity) {
	delete(em.entities, entity)
}

func (em *EntityManager) HasEntity(entity Entity) bool {
	_, exists := em.entities[entity]
	return exists
}

func (em *EntityManager) GetAllEntities() []Entity {
	entities := make([]Entity, 0, len(em.entities))
	for e := range em.entities {
		entities = append(entities, e)
	}
	return entities
}

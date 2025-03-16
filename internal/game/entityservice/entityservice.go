package entityservice

import (
	"log"

	"ecs/pkg/ecs"
)

type EntityService struct {
	world  *ecs.World
	logger *log.Logger
}

func NewEntityService(world *ecs.World, logger *log.Logger) *EntityService {
	return &EntityService{
		world:  world,
		logger: logger,
	}
}

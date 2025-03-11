package main

type TurnManager struct {
	turnOrder []Entity
	current   int
}

func NewTurnManager() *TurnManager {
	return &TurnManager{turnOrder: []Entity{}, current: 0}
}

func (tm *TurnManager) AddEntity(entity Entity) {
	tm.turnOrder = append(tm.turnOrder, entity)
}

func (tm *TurnManager) RemoveEntity(entity Entity) {
	for i, e := range tm.turnOrder {
		if e == entity {
			tm.turnOrder = append(tm.turnOrder[:i], tm.turnOrder[i+1:]...)
			if tm.current > i {
				tm.current--
			}
			return
		}
	}
}

func (tm *TurnManager) NextTurn() Entity {
	if len(tm.turnOrder) == 0 {
		return -1
	}
	tm.current = (tm.current + 1) % len(tm.turnOrder)
	return tm.turnOrder[tm.current]
}

func (tm *TurnManager) GetCurrentEntity() Entity {
	if len(tm.turnOrder) == 0 {
		return -1
	}
	return tm.turnOrder[tm.current]
}

package ecs

import "reflect"

type World struct {
	Entities map[reflect.Type]map[EntityID]interface{}
	LastID   EntityID
}

type Repository interface {
	GetEntities(types ...reflect.Type) []EntityID
	CreateEntity() EntityID
	SetComponent(entity EntityID, component interface{})
	RemoveComponent(entity EntityID, component interface{})
	DestroyEntity(entity EntityID)
}

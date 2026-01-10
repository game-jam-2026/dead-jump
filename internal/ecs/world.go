package ecs

import (
	"fmt"
	"reflect"
)

type World struct {
	Entities  map[reflect.Type]map[EntityID]interface{}
	Resources map[reflect.Type]interface{}
	LastID    EntityID
}

func NewWorld() *World {
	return &World{
		Entities:  make(map[reflect.Type]map[EntityID]interface{}),
		Resources: make(map[reflect.Type]interface{}),
		LastID:    0,
	}
}

func (w *World) CreateEntity() EntityID {
	w.LastID++
	return w.LastID
}

func (w *World) SetComponent(entity EntityID, component interface{}) {
	t := reflect.TypeOf(component)
	if w.Entities[t] == nil {
		w.Entities[t] = make(map[EntityID]interface{})
	}
	w.Entities[t][entity] = component
}

func (w *World) RemoveComponent(entity EntityID, component interface{}) error {
	t := reflect.TypeOf(component)
	if w.Entities[t] == nil {
		return fmt.Errorf("component %v not found", t)
	}

	delete(w.Entities[t], entity)
	return nil
}

func (w *World) GetEntities(types ...reflect.Type) []EntityID {
	if len(types) == 0 {
		return nil
	}

	firstType := types[0]
	if w.Entities[firstType] == nil {
		return nil
	}

	var result []EntityID
	for entityID := range w.Entities[firstType] {
		hasAll := true
		for _, t := range types[1:] {
			if w.Entities[t] == nil || w.Entities[t][entityID] == nil {
				hasAll = false
				break
			}
		}
		if hasAll {
			result = append(result, entityID)
		}
	}
	return result
}

func (w *World) DestroyEntity(entity EntityID) {
	for _, entities := range w.Entities {
		delete(entities, entity)
	}
}

func (w *World) SetResource(resource interface{}) {
	t := reflect.TypeOf(resource)
	w.Resources[t] = resource
}

func GetComponent[C any](w *World, entity EntityID) (*C, error) {
	t := reflect.TypeOf((*C)(nil)).Elem()
	if w.Entities[t] == nil || w.Entities[t][entity] == nil {
		return nil, fmt.Errorf("component %v not found", t)
	}
	component, ok := w.Entities[t][entity].(C)
	if !ok {
		return nil, fmt.Errorf("component has the wrong type")
	}

	return &component, nil
}

func GetResource[R any](w *World) (*R, error) {
	t := reflect.TypeOf((*R)(nil)).Elem()
	if w.Resources[t] == nil {
		return nil, fmt.Errorf("resource %v not found", t)
	}
	resource, ok := w.Resources[t].(R)
	if !ok {
		return nil, fmt.Errorf("resource has the wrong type")
	}

	return &resource, nil
}

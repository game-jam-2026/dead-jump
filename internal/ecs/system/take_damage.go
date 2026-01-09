package system

import "game/internal/ecs/components"

func TakeDamage() {
	// ...type.Type
	entities.Get(components.Health{})
}

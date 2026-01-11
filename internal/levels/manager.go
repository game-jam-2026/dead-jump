package levels

import (
	"github.com/game-jam-2026/dead-jump/internal/assets"
	"github.com/game-jam-2026/dead-jump/internal/ecs"
)

type LevelLoader func() *ecs.World

var LevelSequence = []LevelLoader{
	assets.LoadLoreDumpLevel,
	assets.LoadLevel1,
	assets.LoadPhysicsTestLevel,
}

type Manager struct {
	currentLevel int
}

func NewManager() *Manager {
	return &Manager{
		currentLevel: -1,
	}
}

func (m *Manager) StartGame() *ecs.World {
	m.currentLevel = 0
	return LevelSequence[0]()
}

func (m *Manager) NextLevel() *ecs.World {
	m.currentLevel++
	if m.currentLevel >= len(LevelSequence) {
		return nil
	}
	return LevelSequence[m.currentLevel]()
}

func (m *Manager) RestartLevel() *ecs.World {
	if m.currentLevel < 0 || m.currentLevel >= len(LevelSequence) {
		return nil
	}
	return LevelSequence[m.currentLevel]()
}

func (m *Manager) HasNextLevel() bool {
	return m.currentLevel < len(LevelSequence)-1
}

func (m *Manager) Reset() {
	m.currentLevel = -1
}

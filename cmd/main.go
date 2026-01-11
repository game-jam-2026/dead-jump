package main

import (
	"log"
	"os"
	"reflect"

	"github.com/game-jam-2026/dead-jump/internal/assets"
	"github.com/game-jam-2026/dead-jump/internal/ecs"
	"github.com/game-jam-2026/dead-jump/internal/ecs/components"
	"github.com/game-jam-2026/dead-jump/internal/ecs/systems"
	"github.com/game-jam-2026/dead-jump/internal/levels"
	"github.com/game-jam-2026/dead-jump/internal/menu"
	"github.com/game-jam-2026/dead-jump/internal/physics"
	"github.com/game-jam-2026/dead-jump/internal/utils"
	"github.com/game-jam-2026/dead-jump/internal/utils/audio"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	w            *ecs.World
	menu         *menu.Menu
	levelManager *levels.Manager
}

func NewGame() *Game {
	g := &Game{}

	// Initialize audio context and register all sounds
	assets.InitAudio()

	g.levelManager = levels.NewManager()

	g.menu = menu.NewMenu()
	g.menu.OnStartGame = func() {
		g.w = g.levelManager.StartGame()
		g.menu.SetState(menu.StatePlaying)
	}
	g.menu.OnRestart = func() {
		g.w = g.levelManager.RestartLevel()
	}
	g.menu.OnNextLevel = func() {
		g.w = g.levelManager.NextLevel()
		if g.w == nil {
			g.menu.SetState(menu.StateMenu)
			g.levelManager.Reset()
		}
	}
	g.menu.OnResume = func() {
		g.menu.SetState(menu.StatePlaying)
	}
	g.menu.OnMainMenu = func() {
		g.w = nil
		g.levelManager.Reset()
		g.menu.SetState(menu.StateMenu)
	}
	g.menu.OnQuit = func() {
		os.Exit(0)
	}

	return g
}

func (g *Game) Update() error {
	state := g.menu.GetState()

	systems.UpdateLevelMusic(state)

	// Handle ESC key
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		switch state {
		case menu.StatePlaying:
			g.menu.SetState(menu.StatePaused)
			return nil
		case menu.StatePaused:
			g.menu.SetState(menu.StatePlaying)
			return nil
		}
	}

	// Update based on state
	switch state {
	case menu.StateMenu, menu.StatePaused, menu.StateConfirmRestart, menu.StateSettings, menu.StateLevelComplete, menu.StateGameOver:
		g.menu.Update()
	case menu.StatePlaying:
		if g.w != nil {
			g.updateGame()
		}
	}

	return nil
}

func (g *Game) updateGame() {
	systems.MoveCharacter(g.w)
	systems.UpdateCannons(g.w)

	cfg, _ := ecs.GetResource[physics.Config](g.w)

	systems.ApplyGravity(g.w, cfg)
	systems.ApplyAccumulated(g.w, cfg)
	collisions := systems.ApplyVelocityWithCollisions(g.w, cfg)
	systems.HandleProjectileCollisions(g.w, collisions)
	systems.ApplySpikes(g.w)
	systems.ApplyAnimation(g.w)
	systems.ApplySlopeGravity(g.w, cfg)
	systems.ApplyFriction(g.w, cfg)
	systems.ApplyConveyorBelt(g.w)
	systems.UpdateProjectileLifetime(g.w)
	systems.CleanupOffscreenProjectiles(g.w, assets.WorldWidth, assets.WorldHeight)
	systems.DrawLifeCounter(g.w)

	if systems.ApplyLevelFinish(g.w) {
		g.menu.ShowLevelComplete()
		return
	}

	g.checkGameOver()

	g.updateCameraTarget()
	systems.UpdateCameraSystem(g.w)
}

func (g *Game) checkGameOver() {
	if g.menu.GetState() == menu.StateGameOver {
		return
	}

	lifeEntities := g.w.GetEntities(
		reflect.TypeOf((*components.Life)(nil)).Elem(),
	)

	if len(lifeEntities) == 0 {
		audio.Play(audio.SoundGameOver)
		g.menu.SetState(menu.StateGameOver)
		return
	}

	life, err := ecs.GetComponent[components.Life](g.w, lifeEntities[0])
	if err != nil {
		return
	}

	if life.Count <= 0 {
		audio.Play(audio.SoundGameOver)
		g.menu.SetState(menu.StateGameOver)
	}
}

func (g *Game) updateCameraTarget() {
	camera, err := ecs.GetResource[components.Camera](g.w)
	if err != nil {
		return
	}

	entities := g.w.GetEntities(
		reflect.TypeOf((*components.Character)(nil)).Elem(),
	)
	if len(entities) > 0 {
		camera.Target = int64(entities[0])
		g.w.SetResource(*camera)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	state := g.menu.GetState()

	switch state {
	case menu.StateMenu:
		g.menu.Draw(screen)
	case menu.StatePlaying:
		if g.w != nil {
			camera, _ := ecs.GetResource[components.Camera](g.w)
			systems.DrawSpritesWithCamera(g.w, screen, camera)
			utils.DrawLoreText(g.w, screen)
		}
	case menu.StatePaused, menu.StateConfirmRestart, menu.StateSettings, menu.StateLevelComplete, menu.StateGameOver:
		// Draw game underneath if exists
		if g.w != nil {
			camera, _ := ecs.GetResource[components.Camera](g.w)
			systems.DrawSpritesWithCamera(g.w, screen, camera)
			utils.DrawLoreText(g.w, screen)
		}
		// Draw menu overlay
		g.menu.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return menu.ScreenWidth, menu.ScreenHeight
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Dead Jump")

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}

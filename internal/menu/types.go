package menu

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Screen dimensions
const (
	ScreenWidth  = 320
	ScreenHeight = 240
)

// Animation timing constants
const (
	GlitchMinInterval    = 120  // Minimum frames between random glitches
	GlitchChance         = 0.02 // Probability per frame after interval
	EasterEggDuration    = 180  // Frames to show easter egg
	LetterAppearDelay    = 8    // Frames between each letter appearing
	LetterReturnDuration = 60   // Frames for glitched letters to return
)

// Visual constants
const (
	TitleSpacing        = 8.0  // Pixels between title images
	MenuItemSpacing     = 20.0 // Pixels between menu items
	SettingsItemSpacing = 18.0 // Pixels between settings items
	SkullMargin         = 18.0 // Distance from text to skull indicators
	CharWidthPixels     = 8.0  // Width of one character in pixels
)

// Font sizes
const (
	FontSizeSmall  = 8
	FontSizeMedium = 10
)

// Volume bar constants
const (
	VolumeSteps     = 10
	VolumeStepValue = 0.1
)

// GameState represents the current state of the menu system
type GameState int

const (
	StateMenu GameState = iota
	StatePlaying
	StatePaused
	StateConfirmRestart
	StateSettings
	StateLevelComplete
	StateGameOver
	StateEpilogueEnding
	StateDifficultySelect
)

type SubtitlePhase int

const (
	SubtitleHidden    SubtitlePhase = iota // Not shown
	SubtitleAppearing                      // Letters gathering
	SubtitleVisible                        // Gentle wave animation
	SubtitleGlitching                      // Scatter and return
)

type MenuItem struct {
	Text     string
	Selected bool
	Action   func()
}

type Dialog struct {
	Title    string
	Message  string
	Buttons  []MenuItem
	Selected int
}

type FallingObject struct {
	X, Y     float64
	VelY     float64
	ObjType  int
	Rotation float64
	RotSpeed float64
	Alpha    float64
}

type SubtitleLetter struct {
	OffsetX   float64
	OffsetY   float64
	TargetX   float64
	TargetY   float64
	Alpha     float64
	Rotation  float64
	Scale     float64
	Visible   bool
	EchoAlpha float64
}

type Menu struct {
	state           GameState
	items           []MenuItem
	pauseItems      []MenuItem
	confirmItems    []MenuItem
	settingsItems   []MenuItem
	selectedIndex   int
	previousState   GameState // For returning from settings
	titleOffset     float64
	frameCount      int
	objects         []FallingObject
	dejavuGlitch    float64
	dejavuTimer     int
	secretCode      []ebiten.Key
	secretIndex     int
	easterEggActive bool
	easterEggTimer  int
	screenShake     float64

	// Subtitle letter animation
	subtitleLetters   []SubtitleLetter
	subtitleText      string
	subtitleAnimTimer int
	subtitleShowPhase SubtitlePhase

	// Loaded images
	titleDeadImg *ebiten.Image
	titleJumpImg *ebiten.Image
	skullImg     *ebiten.Image
	objectImages []*ebiten.Image

	// Font
	fontFace   *text.GoTextFaceSource
	fontSmall  *text.GoTextFace
	fontMedium *text.GoTextFace

	// Audio state
	musicPlaying bool

	// Active dialog (for universal dialog system)
	activeDialog *Dialog

	// Level complete/game over data
	levelCompleteItems []MenuItem
	gameOverItems      []MenuItem

	epilogueItems   []MenuItem
	epilogueTimer   int
	difficultyItems []MenuItem

	// Callbacks
	OnStartGame        func()
	OnRestart          func()
	OnQuit             func()
	OnResume           func()
	OnNextLevel        func()
	OnLevelComplete    func()
	OnGameOver         func()
	OnMainMenu         func()
	OnEpilogueComplete func()
}

var (
	colorBgDark       = color.RGBA{R: 10, G: 8, B: 15, A: 255}
	colorBloodRed     = color.RGBA{R: 120, G: 20, B: 30, A: 255}
	colorDeathRed     = color.RGBA{R: 180, G: 40, B: 50, A: 255}
	colorGhostWhite   = color.RGBA{R: 200, G: 190, B: 200, A: 255}
	colorDimGray      = color.RGBA{R: 60, G: 55, B: 70, A: 255}
	colorSelectedGlow = color.RGBA{R: 220, G: 180, B: 160, A: 255}
	colorWarning      = color.RGBA{R: 200, G: 100, B: 80, A: 255}
	colorDeadPurple   = color.RGBA{R: 100, G: 60, B: 120, A: 255}
)

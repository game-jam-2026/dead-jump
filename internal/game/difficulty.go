package game

// Difficulty represents game difficulty level
type Difficulty int

const (
	DifficultyEasy Difficulty = iota
	DifficultyHard
)

var currentDifficulty = DifficultyHard // По умолчанию hard

func GetDifficulty() Difficulty {
	return currentDifficulty
}

func SetDifficulty(d Difficulty) {
	currentDifficulty = d
}

func IsEasyMode() bool {
	return currentDifficulty == DifficultyEasy
}

func IsHardMode() bool {
	return currentDifficulty == DifficultyHard
}

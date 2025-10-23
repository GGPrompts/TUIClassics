package stats

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const (
	statsFileName = "scores.json"
	currentVersion = 1
)

// GetStatsFilePath returns the path to the scores file
func GetStatsFilePath() (string, error) {
	// Use XDG config directory or fall back to home directory
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configDir = filepath.Join(home, ".config")
	}

	tuiDir := filepath.Join(configDir, "tuiclassics")

	// Ensure directory exists
	if err := os.MkdirAll(tuiDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(tuiDir, statsFileName), nil
}

// Load reads the stats file from disk
func Load() (*ScoreData, error) {
	filePath, err := GetStatsFilePath()
	if err != nil {
		return newScoreData(), nil // Return empty data on error
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return newScoreData(), nil // File doesn't exist yet, return empty
		}
		return nil, err
	}

	var scores ScoreData
	if err := json.Unmarshal(data, &scores); err != nil {
		return nil, err
	}

	return &scores, nil
}

// Save writes the stats file to disk
func Save(scores *ScoreData) error {
	filePath, err := GetStatsFilePath()
	if err != nil {
		return err
	}

	scores.Version = currentVersion

	data, err := json.MarshalIndent(scores, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// newScoreData creates an empty score database
func newScoreData() *ScoreData {
	return &ScoreData{
		Version: currentVersion,
	}
}

// UpdateMinesweeperScore updates a minesweeper score and returns achievements
func UpdateMinesweeperScore(scores *ScoreData, difficulty string, timeSeconds int) []Achievement {
	var diffStats *DifficultyStats

	switch difficulty {
	case "Easy":
		diffStats = &scores.Minesweeper.Easy
	case "Medium":
		diffStats = &scores.Minesweeper.Medium
	case "Hard":
		diffStats = &scores.Minesweeper.Hard
	default:
		return nil
	}

	now := time.Now()
	month, week := GetCurrentPeriod()
	achievements := []Achievement{}

	// Check all-time best
	if diffStats.AllTime.TimeSeconds == 0 || timeSeconds < diffStats.AllTime.TimeSeconds {
		diffStats.AllTime = TimeScore{
			TimeSeconds: timeSeconds,
			Date:        now,
		}
		achievements = append(achievements, AchievementAllTime)
	}

	// Check monthly best (reset if new month)
	if ShouldResetPeriod(diffStats.Monthly.Period, "month") || diffStats.Monthly.TimeSeconds == 0 {
		diffStats.Monthly = TimeScore{
			TimeSeconds: timeSeconds,
			Date:        now,
			Period:      month,
		}
		if len(achievements) == 0 { // Only if not already all-time best
			achievements = append(achievements, AchievementMonthly)
		}
	} else if timeSeconds < diffStats.Monthly.TimeSeconds {
		diffStats.Monthly = TimeScore{
			TimeSeconds: timeSeconds,
			Date:        now,
			Period:      month,
		}
		if len(achievements) == 0 {
			achievements = append(achievements, AchievementMonthly)
		}
	}

	// Check weekly best (reset if new week)
	if ShouldResetPeriod(diffStats.Weekly.Period, "week") || diffStats.Weekly.TimeSeconds == 0 {
		diffStats.Weekly = TimeScore{
			TimeSeconds: timeSeconds,
			Date:        now,
			Period:      week,
		}
		if len(achievements) == 0 {
			achievements = append(achievements, AchievementWeekly)
		}
	} else if timeSeconds < diffStats.Weekly.TimeSeconds {
		diffStats.Weekly = TimeScore{
			TimeSeconds: timeSeconds,
			Date:        now,
			Period:      week,
		}
		if len(achievements) == 0 {
			achievements = append(achievements, AchievementWeekly)
		}
	}

	return achievements
}

// Update2048Score updates a 2048 score and returns achievements
func Update2048Score(scores *ScoreData, score int) []Achievement {
	now := time.Now()
	month, week := GetCurrentPeriod()
	achievements := []Achievement{}

	// Check all-time best
	if score > scores.Game2048.AllTime.Score {
		scores.Game2048.AllTime = PointsScore{
			Score: score,
			Date:  now,
		}
		achievements = append(achievements, AchievementAllTime)
	}

	// Check monthly best
	if ShouldResetPeriod(scores.Game2048.Monthly.Period, "month") {
		scores.Game2048.Monthly = PointsScore{
			Score:  score,
			Date:   now,
			Period: month,
		}
		if len(achievements) == 0 {
			achievements = append(achievements, AchievementMonthly)
		}
	} else if score > scores.Game2048.Monthly.Score {
		scores.Game2048.Monthly = PointsScore{
			Score:  score,
			Date:   now,
			Period: month,
		}
		if len(achievements) == 0 {
			achievements = append(achievements, AchievementMonthly)
		}
	}

	// Check weekly best
	if ShouldResetPeriod(scores.Game2048.Weekly.Period, "week") {
		scores.Game2048.Weekly = PointsScore{
			Score:  score,
			Date:   now,
			Period: week,
		}
		if len(achievements) == 0 {
			achievements = append(achievements, AchievementWeekly)
		}
	} else if score > scores.Game2048.Weekly.Score {
		scores.Game2048.Weekly = PointsScore{
			Score:  score,
			Date:   now,
			Period: week,
		}
		if len(achievements) == 0 {
			achievements = append(achievements, AchievementWeekly)
		}
	}

	return achievements
}

// UpdateSolitaireScore updates a solitaire score and returns achievements
func UpdateSolitaireScore(scores *ScoreData, score, timeSeconds, moves int) []Achievement {
	now := time.Now()
	month, week := GetCurrentPeriod()
	achievements := []Achievement{}

	// Check all-time best (highest score)
	if score > scores.Solitaire.AllTime.Score {
		scores.Solitaire.AllTime = SolitaireScore{
			Score:       score,
			TimeSeconds: timeSeconds,
			Moves:       moves,
			Date:        now,
		}
		achievements = append(achievements, AchievementAllTime)
	}

	// Check monthly best
	if ShouldResetPeriod(scores.Solitaire.Monthly.Period, "month") {
		scores.Solitaire.Monthly = SolitaireScore{
			Score:       score,
			TimeSeconds: timeSeconds,
			Moves:       moves,
			Date:        now,
			Period:      month,
		}
		if len(achievements) == 0 {
			achievements = append(achievements, AchievementMonthly)
		}
	} else if score > scores.Solitaire.Monthly.Score {
		scores.Solitaire.Monthly = SolitaireScore{
			Score:       score,
			TimeSeconds: timeSeconds,
			Moves:       moves,
			Date:        now,
			Period:      month,
		}
		if len(achievements) == 0 {
			achievements = append(achievements, AchievementMonthly)
		}
	}

	// Check weekly best
	if ShouldResetPeriod(scores.Solitaire.Weekly.Period, "week") {
		scores.Solitaire.Weekly = SolitaireScore{
			Score:       score,
			TimeSeconds: timeSeconds,
			Moves:       moves,
			Date:        now,
			Period:      week,
		}
		if len(achievements) == 0 {
			achievements = append(achievements, AchievementWeekly)
		}
	} else if score > scores.Solitaire.Weekly.Score {
		scores.Solitaire.Weekly = SolitaireScore{
			Score:       score,
			TimeSeconds: timeSeconds,
			Moves:       moves,
			Date:        now,
			Period:      week,
		}
		if len(achievements) == 0 {
			achievements = append(achievements, AchievementWeekly)
		}
	}

	return achievements
}

// UpdateSnakeScore updates a snake score and returns achievements
func UpdateSnakeScore(scores *ScoreData, score int) []Achievement {
	now := time.Now()
	month, week := GetCurrentPeriod()
	achievements := []Achievement{}

	// Check all-time best
	if score > scores.Snake.AllTime.Score {
		scores.Snake.AllTime = PointsScore{
			Score: score,
			Date:  now,
		}
		achievements = append(achievements, AchievementAllTime)
	}

	// Check monthly best
	if ShouldResetPeriod(scores.Snake.Monthly.Period, "month") {
		scores.Snake.Monthly = PointsScore{
			Score:  score,
			Date:   now,
			Period: month,
		}
		if len(achievements) == 0 {
			achievements = append(achievements, AchievementMonthly)
		}
	} else if score > scores.Snake.Monthly.Score {
		scores.Snake.Monthly = PointsScore{
			Score:  score,
			Date:   now,
			Period: month,
		}
		if len(achievements) == 0 {
			achievements = append(achievements, AchievementMonthly)
		}
	}

	// Check weekly best
	if ShouldResetPeriod(scores.Snake.Weekly.Period, "week") {
		scores.Snake.Weekly = PointsScore{
			Score:  score,
			Date:   now,
			Period: week,
		}
		if len(achievements) == 0 {
			achievements = append(achievements, AchievementWeekly)
		}
	} else if score > scores.Snake.Weekly.Score {
		scores.Snake.Weekly = PointsScore{
			Score:  score,
			Date:   now,
			Period: week,
		}
		if len(achievements) == 0 {
			achievements = append(achievements, AchievementWeekly)
		}
	}

	return achievements
}

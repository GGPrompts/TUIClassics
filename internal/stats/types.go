package stats

import "time"

// ScoreData holds the complete scoring database
type ScoreData struct {
	Version     int                        `json:"version"`
	Minesweeper MinesweeperStats           `json:"minesweeper"`
	Game2048    Game2048Stats              `json:"2048"`
	Solitaire   SolitaireStats             `json:"solitaire"`
	Snake       SnakeStats                 `json:"snake"`
}

// MinesweeperStats holds minesweeper scores per difficulty
type MinesweeperStats struct {
	Easy   DifficultyStats `json:"easy"`
	Medium DifficultyStats `json:"medium"`
	Hard   DifficultyStats `json:"hard"`
}

// DifficultyStats holds time-based scores for a difficulty level
type DifficultyStats struct {
	AllTime TimeScore `json:"all_time"`
	Monthly TimeScore `json:"monthly"`
	Weekly  TimeScore `json:"weekly"`
}

// TimeScore represents a time-based score (lower is better)
type TimeScore struct {
	TimeSeconds int       `json:"time_seconds"`
	Date        time.Time `json:"date"`
	Period      string    `json:"period,omitempty"` // e.g., "2025-01" or "2025-W04"
}

// Game2048Stats holds 2048 game scores
type Game2048Stats struct {
	AllTime PointsScore `json:"all_time"`
	Monthly PointsScore `json:"monthly"`
	Weekly  PointsScore `json:"weekly"`
}

// SolitaireStats holds solitaire game scores
type SolitaireStats struct {
	AllTime SolitaireScore `json:"all_time"`
	Monthly SolitaireScore `json:"monthly"`
	Weekly  SolitaireScore `json:"weekly"`
}

// SnakeStats holds snake game scores
type SnakeStats struct {
	AllTime PointsScore `json:"all_time"`
	Monthly PointsScore `json:"monthly"`
	Weekly  PointsScore `json:"weekly"`
}

// PointsScore represents a points-based score (higher is better)
type PointsScore struct {
	Score  int       `json:"score"`
	Date   time.Time `json:"date"`
	Period string    `json:"period,omitempty"`
}

// SolitaireScore includes both points and time
type SolitaireScore struct {
	Score       int       `json:"score"`
	TimeSeconds int       `json:"time_seconds"`
	Moves       int       `json:"moves"`
	Date        time.Time `json:"date"`
	Period      string    `json:"period,omitempty"`
}

// Achievement represents a newly broken record
type Achievement int

const (
	AchievementNone Achievement = iota
	AchievementWeekly
	AchievementMonthly
	AchievementAllTime
)

func (a Achievement) String() string {
	switch a {
	case AchievementWeekly:
		return "🏆 New Weekly Best!"
	case AchievementMonthly:
		return "🌟 New Monthly Best!"
	case AchievementAllTime:
		return "👑 New All-Time Best!"
	default:
		return ""
	}
}

package balatro

// blinds.go
// Purpose: Blind system for roguelike progression
// Implements three blind types per ante with scaling difficulty

// BlindType represents the three stages of each ante
type BlindType int

const (
	SmallBlind BlindType = iota
	BigBlind
	BossBlind
)

// Blind represents a challenge that must be beaten to progress
type Blind struct {
	Type        BlindType
	Name        string
	TargetScore int
	Reward      int // Money earned for beating this blind
	Ante        int // Which ante this belongs to
}

// Boss blind names from Balatro
var bossBlindNames = []string{
	"The Hook",
	"The Ox",
	"The House",
	"The Wall",
	"The Wheel",
	"The Eye",
	"The Mouth",
	"The Serpent",
	"The Pillar",
	"The Needle",
	"The Head",
	"The Club",
	"The Window",
	"The Manacle",
	"The Mark",
}

// GetBlind calculates and returns the appropriate blind for the given ante and type
func GetBlind(ante int, blindType BlindType) Blind {
	blind := Blind{
		Type: blindType,
		Ante: ante,
	}

	// Calculate base score for small blind: 300 × ante
	smallBlindScore := 300 * ante

	switch blindType {
	case SmallBlind:
		blind.Name = "Small Blind"
		blind.TargetScore = smallBlindScore
		blind.Reward = 3 + ante

	case BigBlind:
		blind.Name = "Big Blind"
		blind.TargetScore = int(float64(smallBlindScore) * 1.5)
		blind.Reward = 4 + ante

	case BossBlind:
		// Pick a boss name based on ante (cycle through the list)
		bossIndex := (ante - 1) % len(bossBlindNames)
		blind.Name = bossBlindNames[bossIndex]
		blind.TargetScore = smallBlindScore * 2
		blind.Reward = 5 + ante
	}

	return blind
}

// String returns a formatted string representation of the blind
func (b Blind) String() string {
	return b.Name
}

// GetBlindTypeName returns the type name as a string
func (bt BlindType) String() string {
	switch bt {
	case SmallBlind:
		return "Small Blind"
	case BigBlind:
		return "Big Blind"
	case BossBlind:
		return "Boss Blind"
	default:
		return "Unknown"
	}
}

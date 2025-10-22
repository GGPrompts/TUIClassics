package hero

import "time"

// New creates a new Hero game model
func New() Model {
	return Model{
		state:       StateMenu,
		scrollSpeed: 10.0, // 10 rows per second
		laneWidth:   12,
		notes:       []Note{},
		multiplier:  1,
		lastHit:     HitMiss,
		songIndex:   0,
	}
}

// startSong initializes the game with a song
func (m *Model) startSong(song *Song) {
	m.currentSong = song
	m.state = StatePlaying
	m.notes = []Note{}
	m.score = 0
	m.combo = 0
	m.multiplier = 1
	m.lastHit = HitMiss
	m.startTime = time.Now()
	m.currentTime = m.startTime
	m.showHitFeedback = false

	// Convert chart notes to game notes
	for _, chartNote := range song.Chart {
		m.notes = append(m.notes, Note{
			Lane:    chartNote.Lane,
			HitTime: m.startTime.Add(time.Duration(chartNote.Time * float64(time.Second))),
			Y:       -10, // Start off screen
			Hit:     false,
		})
	}
}

// checkSongEnd checks if the song has finished
func (m *Model) checkSongEnd() {
	if m.state != StatePlaying {
		return
	}

	// Check if all notes have been processed
	allNotesProcessed := true
	for _, note := range m.notes {
		if !note.Hit && note.Y <= NoteAreaHeight+GracePeriod {
			allNotesProcessed = false
			break
		}
	}

	// Also check if we've passed the song duration
	elapsed := m.currentTime.Sub(m.startTime).Seconds()
	if allNotesProcessed && elapsed > m.currentSong.Duration+2.0 {
		m.state = StateFinished
	}
}

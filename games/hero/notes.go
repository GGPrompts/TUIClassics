package hero

import "math"

// updateNotePositions updates all note Y positions based on current time
func (m *Model) updateNotePositions() {
	elapsed := m.currentTime.Sub(m.startTime).Seconds()

	for i := range m.notes {
		if !m.notes[i].Hit {
			// Calculate Y position based on time until hit
			timeToHit := m.notes[i].HitTime.Sub(m.startTime).Seconds()
			distanceFromHit := (timeToHit - elapsed) * m.scrollSpeed

			// Note appears at the top (Y=0) when it's scrollSpeed/noteRows seconds away from hit time
			// Hit zone is at Y=NoteAreaHeight
			// Map the time distance to screen position
			m.notes[i].Y = NoteAreaHeight - int(distanceFromHit)
		}
	}
}

// checkMissedNotes removes notes that have passed the hit zone
func (m *Model) checkMissedNotes() {
	newNotes := []Note{}

	for _, note := range m.notes {
		// Keep notes that haven't passed the grace period yet
		if note.Y <= NoteAreaHeight+GracePeriod {
			newNotes = append(newNotes, note)
		} else if !note.Hit {
			// Missed! Reset combo
			m.combo = 0
			m.multiplier = 1
		}
	}

	m.notes = newNotes
}

// getNotesInLane returns all unhit notes in a specific lane
func (m *Model) getNotesInLane(lane int) []*Note {
	var laneNotes []*Note
	for i := range m.notes {
		if m.notes[i].Lane == lane && !m.notes[i].Hit {
			laneNotes = append(laneNotes, &m.notes[i])
		}
	}
	return laneNotes
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// absFloat returns the absolute value of a float64
func absFloat(x float64) float64 {
	return math.Abs(x)
}

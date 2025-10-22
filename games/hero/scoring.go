package hero

import "time"

// checkHit checks if a key press successfully hit a note in the given lane
func (m *Model) checkHit(lane int) HitResult {
	laneNotes := m.getNotesInLane(lane)

	if len(laneNotes) == 0 {
		// No notes in this lane, it's a miss
		m.resetCombo()
		return HitMiss
	}

	// Find the closest note to the hit zone
	var closest *Note
	minDistance := 999

	for _, note := range laneNotes {
		distance := abs(note.Y - NoteAreaHeight)
		if distance < minDistance {
			minDistance = distance
			closest = note
		}
	}

	// Check timing windows
	// Perfect: within 1 row (50ms equivalent)
	// Good: within 2 rows (100ms equivalent)
	// OK: within 3 rows (150ms equivalent)
	// Miss: beyond 3 rows

	var result HitResult

	if minDistance <= 1 {
		result = HitPerfect
		closest.Hit = true
		m.addScore(100)
	} else if minDistance <= 2 {
		result = HitGood
		closest.Hit = true
		m.addScore(50)
	} else if minDistance <= 3 {
		result = HitOK
		closest.Hit = true
		m.addScore(25)
	} else {
		result = HitMiss
		m.resetCombo()
	}

	// Show visual feedback
	m.lastHit = result
	m.showHitFeedback = true
	m.hitFeedbackTime = time.Now()

	return result
}

// addScore adds points to the score and updates combo/multiplier
func (m *Model) addScore(basePoints int) {
	m.combo++

	// Update multiplier based on combo milestones
	if m.combo >= 30 {
		m.multiplier = 4
	} else if m.combo >= 20 {
		m.multiplier = 3
	} else if m.combo >= 10 {
		m.multiplier = 2
	} else {
		m.multiplier = 1
	}

	// Apply multiplier to score
	m.score += basePoints * m.multiplier
}

// resetCombo resets the combo counter and multiplier
func (m *Model) resetCombo() {
	m.combo = 0
	m.multiplier = 1
}

package solitaire

// CanMoveToTableau checks if a card (or sequence) can be moved to a tableau pile
func (m *Model) CanMoveToTableau(cards []Card, tableauIndex int) bool {
	if tableauIndex < 0 || tableauIndex >= 7 {
		return false
	}

	pile := m.tableau[tableauIndex]
	firstCard := cards[0]

	// If pile is empty, only Kings can be placed
	if len(pile.Cards) == 0 {
		return firstCard.IsKing()
	}

	// Check if the first card can stack on the top card of the tableau
	topCard := pile.Cards[len(pile.Cards)-1]
	return firstCard.CanStackOn(topCard)
}

// CanMoveToFoundation checks if a card can be moved to a foundation pile
func (m *Model) CanMoveToFoundation(card Card, foundationIndex int) bool {
	if foundationIndex < 0 || foundationIndex >= 4 {
		return false
	}

	pile := m.foundation[foundationIndex]

	// If pile is empty, only Aces can be placed
	if len(pile.Cards) == 0 {
		return card.Rank == Ace
	}

	// Check if card can be placed on foundation
	return card.CanMoveToFoundation(pile.Cards)
}

// GetMovableSequence returns the sequence of cards that can be moved from a tableau pile
// starting at the given index. Returns nil if the sequence is invalid.
func (m *Model) GetMovableSequence(tableauIndex, startIndex int) []Card {
	if tableauIndex < 0 || tableauIndex >= 7 {
		return nil
	}

	pile := m.tableau[tableauIndex]
	if startIndex < 0 || startIndex >= len(pile.Cards) {
		return nil
	}

	// Check if the card at startIndex is face up
	if !pile.Cards[startIndex].FaceUp {
		return nil
	}

	// Build the sequence from startIndex to the end
	sequence := pile.Cards[startIndex:]

	// Validate the sequence (must be alternating colors, descending ranks)
	for i := 0; i < len(sequence)-1; i++ {
		if !sequence[i+1].CanStackOn(sequence[i]) {
			return nil
		}
	}

	return sequence
}

// MoveCards moves cards from one pile to another
func (m *Model) MoveCards(from, to CursorLocation, cards []Card) bool {
	// Validate the move
	valid := false

	if to.PileType == TableauPile {
		valid = m.CanMoveToTableau(cards, to.PileIndex)
	} else if to.PileType == FoundationPile {
		// Can only move single cards to foundation
		if len(cards) == 1 {
			valid = m.CanMoveToFoundation(cards[0], to.PileIndex)
		}
	}

	if !valid {
		return false
	}

	// Remove cards from source pile
	switch from.PileType {
	case TableauPile:
		pile := &m.tableau[from.PileIndex]
		pile.Cards = pile.Cards[:len(pile.Cards)-len(cards)]

		// Flip the top card if it's face down
		if len(pile.Cards) > 0 && !pile.Cards[len(pile.Cards)-1].FaceUp {
			pile.Cards[len(pile.Cards)-1].FaceUp = true
			m.AddScore(5) // +5 for revealing a card
		}

	case WastePile:
		// Remove the top card from waste
		m.waste.Cards = m.waste.Cards[:len(m.waste.Cards)-1]

	case FoundationPile:
		// Remove cards from foundation
		pile := &m.foundation[from.PileIndex]
		pile.Cards = pile.Cards[:len(pile.Cards)-len(cards)]
		m.AddScore(-15) // -15 for moving from foundation
	}

	// Add cards to destination pile
	switch to.PileType {
	case TableauPile:
		pile := &m.tableau[to.PileIndex]
		pile.Cards = append(pile.Cards, cards...)

		// Score for moving from waste to tableau
		if from.PileType == WastePile {
			m.AddScore(5)
		}

	case FoundationPile:
		pile := &m.foundation[to.PileIndex]
		pile.Cards = append(pile.Cards, cards...)
		m.AddScore(10) // +10 for moving to foundation
	}

	m.moves++
	return true
}

// AutoMoveToFoundation attempts to automatically move a card to the appropriate foundation
func (m *Model) AutoMoveToFoundation(card Card, from CursorLocation) bool {
	// Try each foundation pile
	for i := 0; i < 4; i++ {
		if m.CanMoveToFoundation(card, i) {
			to := CursorLocation{PileType: FoundationPile, PileIndex: i}
			return m.MoveCards(from, to, []Card{card})
		}
	}
	return false
}

# Next Session: Implement Phase 3 - Jokers & Effects

**Date**: 2025-10-22 (Updated after Phase 2 completion)
**Current Status**: ✅ Phase 2 Complete - Full roguelike progression!
**Next Task**: Implement Phase 3 (Jokers, Card Enhancements, Enhanced Scoring)

---

## 🎉 What's Already Done

### Phase 1: Core Game ✅
- ✅ Complete card system (cards.go - 473 lines)
- ✅ Poker hand evaluation (hands.go - 378 lines)
- ✅ Scoring engine (scoring.go - 191 lines)
- ✅ Beautiful centered UI with neon colors
- ✅ Card selection, playing hands, discarding
- ✅ Compact layouts that work on small terminals
- ✅ Trippy lava lamp landing page
- ✅ All controls working (1-8, arrows, space, enter, D)
- ✅ **BUG FIX**: Ten cards now display correctly (was showing ":" instead of "10")

### Phase 2: Game Loop & Progression ✅ (JUST COMPLETED!)
- ✅ **blinds.go** (96 lines) - Blind system with 3 types
  - Small Blind: 300×ante chips, $3+ante reward
  - Big Blind: 1.5×Small (450 chips), $4+ante reward
  - Boss Blind: 2×Small (600 chips), $5+ante reward, 15 unique names
- ✅ **round.go** (119 lines) - Round state management
  - Ante progression
  - Win/loss detection
  - Hands and discards tracking
- ✅ **Game phases** (types.go + view.go + 219 lines)
  - PhaseSelectCards - Normal gameplay
  - PhaseBlindComplete - Victory screen with earnings
  - PhaseShop - Shop placeholder (basic)
  - PhaseGameOver - Loss screen with restart option
- ✅ **Phase transitions** (update_keyboard.go + 161 lines)
  - Enter on victory → Shop
  - Enter on shop → Next blind
  - R on game over → Restart
  - Q anywhere → Quit to menu
- ✅ **Integration** (model.go updated)
  - newGame() initializes with Ante 1, Small Blind
  - playHand() handles win/loss transitions
  - Full game loop works: Play → Win → Shop → Next Blind

**Binary**: `~/projects/balatro-tui/balatro-tui` (5.0 MB, builds cleanly)

---

## 🚀 Phase 3: Jokers & Effects (~800 LOC)

**Goal**: Add the magic! Jokers that modify scoring, card enhancements, and a proper scoring breakdown.

### Implementation Plan (from `docs/IMPLEMENTATION_PLAN.md`)

#### **Step 1: Joker System** (~300 LOC)
**File**: Create `jokers.go`

**Core Types**:
```go
type JokerEffect int
const (
    NoEffect JokerEffect = iota
    AddChips       // +N chips
    AddMult        // +N mult
    MultIfPair     // xN mult if hand contains pair
    ChipsPerCard   // +N chips per card of rank
    MultIfSuit     // +N mult per card of suit
    RetriggerFirst // Retrigger first card
)

type Joker struct {
    ID          string
    Name        string
    Description string
    Rarity      Rarity
    Effect      JokerEffect
    Value       int   // Effect magnitude
    Counter     int   // State for stateful jokers
}

type Rarity int
const (
    Common Rarity = iota
    Uncommon
    Rare
    Legendary
)
```

**Starting Jokers** (implement 5 simple ones):
1. **Joker** - Common - "+4 Mult"
2. **Greedy Joker** - Common - "Played Diamond cards give +3 Mult"
3. **Lusty Joker** - Common - "Played Heart cards give +3 Mult"
4. **Wrathful Joker** - Common - "Played Spade cards give +3 Mult"
5. **Gluttonous Joker** - Common - "Played Club cards give +3 Mult"

**Key Function**:
```go
func (j Joker) Apply(calc *ScoreCalculation, hand HandInfo, played []Card)
```

---

#### **Step 2: Card Enhancements** (~200 LOC)
**File**: Extend `cards.go`

**Already have the enums** - Just need to implement effects:
- **Bonus**: +30 chips
- **Mult**: +4 mult
- **Steel**: +50 chips
- **Stone**: 50 chips (fixed, ignores rank)
- **Glass**: x2 mult (then destroys after scoring)

**Already have editions** - Implement effects:
- **Foil**: +50 chips
- **Holographic**: +10 mult
- **Polychrome**: x1.5 mult

**Update scoring.go**:
```go
// Current: CalculateScore(handInfo, cards)
// New:     CalculateScore(handInfo, cards, jokers)
```

---

#### **Step 3: Enhanced Scoring Engine** (~300 LOC)
**File**: Rewrite `scoring.go`

**New calculation order**:
1. Base hand chips + mult (Pair = 10 + 2×)
2. Add card chip values (face cards, enhancements)
3. Apply card mult effects (Mult enhancement, Glass, editions)
4. **Apply joker effects IN ORDER** (order matters!)
5. Calculate final: chips × mult

**Track all modifiers** for display:
```go
type ScoreCalculation struct {
    BaseChips      int
    BaseMult       int
    TotalChips     int
    TotalMult      int
    FinalScore     int
    ChipModifiers  []string  // "+10 (5 of Hearts)", "+30 (Bonus)"
    MultModifiers  []string  // "+4 (Mult card)", "x2 (Glass)"
}
```

**Glass card destruction**:
- After scoring, remove Glass cards from hand
- Already have `ApplyGlassDestruction()` in model.go!

---

## 📝 Implementation Checklist

### Part 1: Joker System
- [ ] Create `jokers.go`
- [ ] Implement `JokerEffect` enum and `Joker` struct
- [ ] Implement `Rarity` enum
- [ ] Create array of 5 starting jokers
- [ ] Write `Apply()` method for each effect type
- [ ] Add `jokers []Joker` field to model
- [ ] Test: Create jokers, verify effects apply correctly

### Part 2: Card Enhancements
- [ ] Update `GetChipValue()` in `cards.go` to handle enhancements
- [ ] Create `ApplyMultEffects()` method in `cards.go`
- [ ] Handle edition effects (Foil, Holo, Poly)
- [ ] Test: Create enhanced cards, verify chip/mult values

### Part 3: Enhanced Scoring
- [ ] Update `ScoreCalculation` struct with modifier arrays
- [ ] Rewrite `CalculateScore()` to accept jokers parameter
- [ ] Implement 4-step calculation (chips → card mults → jokers → final)
- [ ] Track all modifiers for display
- [ ] Update `renderScoreBreakdown()` to show detailed modifiers
- [ ] Test: Score hands with jokers, verify order matters

### Part 4: UI Integration
- [ ] Add joker display to game view (show owned jokers)
- [ ] Update shop to sell jokers (basic - just show price)
- [ ] Add "buy joker" handler in shop phase
- [ ] Render joker cards (compact style)
- [ ] Test: Buy joker → see it in game → verify scoring changes

### Part 5: Testing
- [ ] Test all 5 joker effects work correctly
- [ ] Test card enhancements modify score
- [ ] Test Glass cards destroy after scoring
- [ ] Test edition effects (Foil, Holo, Poly)
- [ ] Test scoring breakdown shows all steps
- [ ] Test joker order matters (multiplicative effects)

---

## 🎯 Success Criteria for Phase 3

When complete, you should be able to:
1. ✅ Own up to 5 jokers (joker slots)
2. ✅ See jokers displayed in game view
3. ✅ Buy jokers from shop
4. ✅ Jokers modify score when playing hands
5. ✅ Card enhancements work (Bonus, Mult, Steel, Stone, Glass)
6. ✅ Card editions work (Foil, Holographic, Polychrome)
7. ✅ Glass cards destroy after scoring
8. ✅ Scoring breakdown shows all modifiers
9. ✅ Joker order matters (test multiplicative effects)

---

## 📚 Reference Documents

- `docs/balatro/IMPLEMENTATION_PLAN.md` - Complete implementation roadmap
- `docs/balatro/REFERENCES.md` - **NEW! External resources & data**
  - **Balatro Calculator**: https://github.com/EFHIII/balatro-calculator
  - **150 Joker definitions** with full mechanics
  - **All card types**: Enhancements, Editions, Seals, Tarots, Planets, Spectrals
  - **Boss Blind effects** reference
  - **Scoring mechanics** documentation
- `games/balatro/cards.go` - Enhancement and Edition enums (already defined!)
- `games/balatro/scoring.go` - Current scoring logic (will be extended)
- `games/balatro/jokers.go` - Current joker implementations (5 done, 145 to go!)

---

## 🎨 UI Design Notes

**Joker Display** (in game view):
- Show owned jokers in a row above the hand
- Compact cards (same style as playing cards)
- Show joker name and effect on hover/select
- Max 5 jokers visible

**Shop Joker Display**:
- Show 2 jokers for sale
- Display: Name, Description, Rarity, Cost
- Highlight selected with gold border
- Press 1-2 to select, Enter to buy

**Scoring Breakdown** (after playing hand):
- Show step-by-step calculation
- List all chip modifiers
- List all mult modifiers
- Make it clear how final score was calculated
- Example:
  ```
  SCORE: 450
  ─────────────
  Base: 10 chips × 2 mult
  +10 (5♥), +10 (6♥), +10 (7♥)
  +4 mult (Joker)
  +3 mult (Greedy Joker - Diamond)
  ─────────────
  30 chips × 9 mult = 270
  ```

---

## 💡 Pro Tips

1. **Start simple** - Implement AddMult and AddChips effects first
2. **Test incrementally** - Add one joker at a time
3. **Reuse patterns** - Copy card rendering style for jokers
4. **Order matters** - Apply chip effects before mult effects
5. **Glass is special** - Handle destruction in model.go, not scoring.go

---

## 🚀 Quick Start Command

```bash
cd ~/projects/balatro-tui

# 1. Create jokers.go
touch jokers.go

# 2. Start with the Joker type and 5 basic jokers
# 3. Update scoring.go to accept jokers parameter
# 4. Update model.go to store jokers and pass to scoring
# 5. Test with simple AddMult joker

# Build and test
go build && ./balatro-tui
```

---

## 📊 Current Stats

**Lines of Code**:
- Phase 1: ~1,400 LOC (cards, hands, scoring, UI)
- Phase 2: ~595 LOC (blinds, rounds, phases, transitions)
- **Phase 3 Target**: ~800 LOC (jokers, enhancements, scoring v2)
- **Total after Phase 3**: ~2,800 LOC

**Binary Size**: 5.0 MB (will grow slightly with jokers)

---

## 🔜 Future Phases (After Phase 3)

- **Phase 4**: Shop & Economy (~500 LOC)
  - Proper shop with rerolls
  - Booster packs (Tarot, Planet, Spectral)
  - Vouchers
  - Money management

- **Phase 5**: Polish & Animations (~800 LOC)
  - Scoring animations (numbers flying)
  - Card flip animations
  - Shop animations
  - Sound effects (optional)

- **Phase 6**: Advanced Features (~1000 LOC)
  - Boss blind special effects
  - Tarot and Planet cards
  - Spectral cards
  - Vouchers
  - Deck types

---

**Ready to add the magic! Let's implement Phase 3! 🃏✨**

# Balatro TUI - Implementation Plan

**Status**: Phase 3 Complete + UI Polish Added ✅
**Timeline**: 4-5 weeks to MVP + polish
**Estimated Size**: 3,000-4,000 LOC
**Current LOC**: ~2,800 LOC

---

## 🎉 Latest Updates (Current Session)

### Phase 3: COMPLETE ✅
- ✅ Joker system fully implemented
- ✅ Enhanced scoring engine with joker effects
- ✅ Card enhancements working (Bonus, Mult, Glass, Steel, Stone, Gold, Lucky)
- ✅ Card editions working (Foil, Holographic, Polychrome)
- ✅ Card seals system implemented
- ✅ Shop with joker purchasing
- ✅ Score breakdown showing all modifiers

### UI Polish: COMPLETE ✅
- ✅ **Color-coded game info header** (Hands in blue, Discards in red) - view.go:371-415
- ✅ **Reserved space for hand info** - Prevents UI jumping when selecting cards - view.go:684-685
- ✅ **Full mouse support for card selection** - Click cards to toggle for play - update_mouse.go:28-54
- ✅ **Full mouse support for shop** - Click jokers to select/buy - update_mouse.go:57-86
- ✅ **Coordinate calculation patterns** - Following CLAUDE.md patterns from solitaire

### Mouse Implementation Details
- Uses click vs drag detection (< 2 pixel movement = click)
- Coordinate calculations match rendering logic exactly
- Phase-aware handling (game phase vs shop phase)
- Mouse state tracking in Model (mousePressX, mousePressY)
- Shop joker selection and purchasing with visual feedback

---

## UI Design Decision: Card Display Strategy

### The Challenge
Balatro cards need to display MORE information than solitaire:
- Rank + Suit (standard)
- Enhancement (Bonus/Mult/Glass/Steel/Stone/Gold)
- Edition (Foil/Holographic/Polychrome)
- Seals (Red/Blue/Gold/Purple)

### Solution: Dual Display System

**Option A: Large Cards Only** (7-9 chars wide)
```
╭─────────╮
│ A ♠     │
│         │
│  +30    │  ← Enhancement
│  STEEL  │  ← Type
│  🌟HOLO │  ← Edition
╰─────────╯
```
❌ Problem: 5-8 cards won't fit horizontally on narrow terminals

**Option B: Compact Hand + Detail Panel** ✅ RECOMMENDED
```
HAND (Compact 5-wide):               SELECTED CARD (Detail):
╭─────╮ ╭─────╮ ╭─────╮             ╭───────────────────╮
│ A ♠ │ │ K ♠ │ │ Q ♠ │             │   ACE OF SPADES   │
│ [*] │ │     │ │     │             │                   │
╰─────╯ ╰─────╯ ╰─────╯             │ Enhancement: +30  │
  [1]     [2]     [3]                │ Type: Steel       │
                                     │ Edition: Holo     │
                                     │ Seal: Red         │
                                     │                   │
                                     │ +50 chips while   │
                                     │ this card is in   │
                                     │ hand (Steel)      │
                                     ╰───────────────────╯
```

**Benefits**:
- Hand stays compact (5 chars wide like solitaire)
- Arrow keys select cards
- Right panel shows ALL details of selected card
- Works on narrow terminals
- Similar to how Balatro shows card details on hover

---

## Phase 1: Core Foundation (MVP) - ~800 LOC

### Week 1: Get Basic Poker Working

**Goal**: Play 5 cards, detect poker hand, show score

#### 1.1 Card System (~150 LOC)
**File**: `cards.go`

```go
type Suit int
const (
    Clubs Suit = iota
    Diamonds
    Hearts
    Spades
)

type Rank int
const (
    Ace Rank = 1
    Two Rank = 2
    // ... through King = 13
)

type Enhancement int
const (
    NoEnhancement Enhancement = iota
    Bonus      // +30 chips
    Mult       // +4 mult
    Glass      // x2 mult, destroys when scored
    Steel      // +50 chips while in hand
    Stone      // +50 chips, no rank
    Gold       // $3 at end of round
    Lucky      // 1 in 5 chance +20 mult
)

type Edition int
const (
    NoEdition Edition = iota
    Foil         // +50 chips
    Holographic  // +10 mult
    Polychrome   // x1.5 mult
)

type Card struct {
    Suit        Suit
    Rank        Rank
    Enhancement Enhancement
    Edition     Edition
    Seal        Seal
    Selected    bool  // For UI
}

func NewDeck() []Card {
    // Standard 52-card deck
}

func Shuffle(deck []Card) {
    // Fisher-Yates shuffle
}
```

#### 1.2 Poker Hand Evaluation (~250 LOC)
**File**: `hands.go`

```go
type PokerHand int
const (
    HighCard PokerHand = iota
    Pair
    TwoPair
    ThreeOfKind
    Straight
    Flush
    FullHouse
    FourOfKind
    StraightFlush
    RoyalFlush
)

// HandInfo stores detected hand details
type HandInfo struct {
    Type       PokerHand
    Cards      []Card     // The 5 cards that make the hand
    BaseChips  int        // Base chips for this hand type
    BaseMult   int        // Base multiplier
    Level      int        // Hand level (from Planet cards)
}

func EvaluateHand(cards []Card) HandInfo {
    // Detect best poker hand from 5 cards
    // Order: RoyalFlush > StraightFlush > ... > HighCard
}

func FindBestPlay(hand []Card) ([]Card, HandInfo) {
    // From 5-8 cards, find best 5-card poker hand
    // Iterate all combinations
}

// Helper functions
func isFlush(cards []Card) bool
func isStraight(cards []Card) bool
func rankCounts(cards []Card) map[Rank]int
```

**Base Hand Values** (from Balatro):
```go
var baseHandValues = map[PokerHand]HandInfo{
    HighCard:      {BaseChips: 5,   BaseMult: 1},
    Pair:          {BaseChips: 10,  BaseMult: 2},
    TwoPair:       {BaseChips: 20,  BaseMult: 2},
    ThreeOfKind:   {BaseChips: 30,  BaseMult: 3},
    Straight:      {BaseChips: 30,  BaseMult: 4},
    Flush:         {BaseChips: 35,  BaseMult: 4},
    FullHouse:     {BaseChips: 40,  BaseMult: 4},
    FourOfKind:    {BaseChips: 60,  BaseMult: 7},
    StraightFlush: {BaseChips: 100, BaseMult: 8},
    RoyalFlush:    {BaseChips: 100, BaseMult: 8},
}
```

#### 1.3 Basic Scoring (~150 LOC)
**File**: `scoring.go`

```go
type ScoreCalculation struct {
    BaseChips      int
    ChipModifiers  []string  // "+30 (Bonus)", "+10 (Ace)"
    TotalChips     int

    BaseMult       int
    MultModifiers  []string  // "+4 (Mult card)", "x2 (Glass)"
    TotalMult      int

    FinalScore     int       // Chips * Mult
}

func CalculateScore(hand HandInfo, playedCards []Card) ScoreCalculation {
    calc := ScoreCalculation{
        BaseChips: hand.BaseChips,
        BaseMult:  hand.BaseMult,
    }

    // Add chip bonuses from played cards
    for _, card := range playedCards {
        calc.TotalChips += card.GetChipValue()
        calc.ChipModifiers = append(calc.ChipModifiers,
            fmt.Sprintf("+%d (%s)", card.GetChipValue(), card.Rank))
    }

    // Add mult bonuses
    // (Jokers will go here later)

    calc.FinalScore = calc.TotalChips * calc.TotalMult
    return calc
}
```

#### 1.4 Simple UI (~250 LOC)
**File**: `view.go`

```go
func (m Model) View() string {
    switch m.state {
    case StateSelectCards:
        return m.viewSelectCards()
    case StateScoring:
        return m.viewScoring()
    case StateGameOver:
        return m.viewGameOver()
    }
}

func (m Model) viewSelectCards() string {
    var b strings.Builder

    // Title
    b.WriteString(m.renderTitle())
    b.WriteString("\n\n")

    // Game info (round, blind, target score)
    b.WriteString(m.renderGameInfo())
    b.WriteString("\n\n")

    // Hand (compact 5-wide cards)
    b.WriteString(m.renderHand())
    b.WriteString("\n\n")

    // Selected card detail panel (if card selected)
    if m.selectedCardIndex >= 0 {
        b.WriteString(m.renderCardDetail(m.hand[m.selectedCardIndex]))
    }

    // Controls hint
    b.WriteString("1-8: Select | Space: Toggle | Enter: Play | D: Discard")

    return b.String()
}

func (m Model) renderCard(card Card, selected bool) string {
    // Compact 5-wide card
    style := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        Width(5).
        Height(3)

    if selected {
        style = style.BorderForeground(lipgloss.Color("226"))  // Gold
    }

    // Top line: Rank + Suit
    top := fmt.Sprintf("%s %s", card.Rank.String(), card.Suit.Symbol())

    // Middle: Enhancement indicator
    mid := card.Enhancement.ShortString()  // "+30", "STL", "GLS", etc.

    return style.Render(
        fmt.Sprintf("%s\n%s",
            centerText(top, 3),
            centerText(mid, 3)))
}

func (m Model) renderCardDetail(card Card) string {
    // Large detail panel (18 chars wide)
    var b strings.Builder

    b.WriteString(fmt.Sprintf("   %s OF %s\n", card.Rank.LongString(), card.Suit.LongString()))
    b.WriteString("\n")

    if card.Enhancement != NoEnhancement {
        b.WriteString(fmt.Sprintf("Enhancement: %s\n", card.Enhancement))
        b.WriteString(fmt.Sprintf("  %s\n", card.Enhancement.Description()))
    }

    if card.Edition != NoEdition {
        b.WriteString(fmt.Sprintf("Edition: %s\n", card.Edition))
    }

    // Wrap in panel
    return lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        Padding(1).
        Render(b.String())
}
```

**Deliverable**: Can play 5 cards, see poker hand detected, see basic score.

---

## Phase 2: Game Loop & Progression - ~600 LOC

### Week 2: Make It a Game

#### 2.1 Blind System (~200 LOC)
**File**: `blind.go`

```go
type BlindType int
const (
    SmallBlind BlindType = iota
    BigBlind
    BossBlind
)

type Blind struct {
    Type       BlindType
    Name       string
    TargetScore int
    Reward     int  // Money earned

    // Boss blind effects (later)
    Effect     string
}

func GetBlind(ante int, blindType BlindType) Blind {
    // Calculate target score based on ante
    // Formula: baseScore * (ante multiplier)
    baseScore := 300
    multiplier := 1.0 + (float64(ante-1) * 0.5)  // Scales with ante

    switch blindType {
    case SmallBlind:
        return Blind{
            Type: SmallBlind,
            Name: "Small Blind",
            TargetScore: int(float64(baseScore) * multiplier),
            Reward: 3 + ante,
        }
    case BigBlind:
        return Blind{
            Type: BigBlind,
            Name: "Big Blind",
            TargetScore: int(float64(baseScore) * multiplier * 1.5),
            Reward: 4 + ante,
        }
    case BossBlind:
        // Boss blinds have special effects
        return Blind{
            Type: BossBlind,
            Name: getBossBlindName(ante),
            TargetScore: int(float64(baseScore) * multiplier * 2.0),
            Reward: 5 + ante,
        }
    }
}
```

#### 2.2 Round Management (~200 LOC)
**File**: `round.go`

```go
type RoundState struct {
    Ante           int
    CurrentBlind   Blind
    BlindProgress  int  // 0=Small, 1=Big, 2=Boss

    HandsRemaining int  // Decrements each play
    DiscardsRemaining int

    CurrentScore   int
    TotalMoney     int
}

func NewRound(ante int) RoundState {
    return RoundState{
        Ante: ante,
        CurrentBlind: GetBlind(ante, SmallBlind),
        BlindProgress: 0,
        HandsRemaining: 4,      // Default: 4 hands per blind
        DiscardsRemaining: 3,   // Default: 3 discards per blind
        CurrentScore: 0,
        TotalMoney: 4,          // Starting cash
    }
}

func (r *RoundState) PlayHand(score int) (won bool, gameOver bool) {
    r.HandsRemaining--
    r.CurrentScore += score

    // Check win condition
    if r.CurrentScore >= r.CurrentBlind.TargetScore {
        r.TotalMoney += r.CurrentBlind.Reward
        return true, false
    }

    // Check loss condition
    if r.HandsRemaining <= 0 {
        return false, true  // Game over!
    }

    return false, false
}

func (r *RoundState) AdvanceBlind() {
    r.BlindProgress++
    r.CurrentScore = 0
    r.HandsRemaining = 4
    r.DiscardsRemaining = 3

    switch r.BlindProgress {
    case 1:
        r.CurrentBlind = GetBlind(r.Ante, BigBlind)
    case 2:
        r.CurrentBlind = GetBlind(r.Ante, BossBlind)
    case 3:
        // Advance to next ante
        r.Ante++
        r.BlindProgress = 0
        r.CurrentBlind = GetBlind(r.Ante, SmallBlind)
    }
}
```

#### 2.3 Game State Machine (~200 LOC)
**File**: `types.go` + `model.go`

```go
type GameState int
const (
    StateMainMenu GameState = iota
    StateSelectCards
    StateScoring        // Show scoring animation
    StateBlindComplete  // Show blind complete screen
    StateShop
    StateGameOver
)

type Model struct {
    // Terminal dimensions
    width, height int

    // Game state
    state      GameState
    round      RoundState

    // Deck management
    deck       []Card
    hand       []Card
    playedCards []Card

    // UI state
    selectedCardIndex int
    cursorPosition    int

    // Animation
    animating     bool
    animationFrame int
    scoreCalc     ScoreCalculation
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch m.state {
    case StateSelectCards:
        return m.updateSelectCards(msg)
    case StateScoring:
        return m.updateScoring(msg)
    case StateShop:
        return m.updateShop(msg)
    }
    return m, nil
}
```

**Deliverable**: Playable game loop - beat small blind, advance to big blind, beat big blind, go to shop, next ante.

---

## Phase 3: Jokers & Effects - ~800 LOC ✅ COMPLETE

### Week 3: Add the Magic ✅

#### 3.1 Joker System (~300 LOC)
**File**: `jokers.go`

```go
type JokerEffect int
const (
    NoEffect JokerEffect = iota
    AddChips       // +N chips
    AddMult        // +N mult
    MultIfPair     // xN mult if hand contains pair
    ChipsPerCard   // +N chips per card of rank
    RetriggerFirst // Retrigger first card
    // ... more effects
)

type Joker struct {
    ID          string
    Name        string
    Description string
    Rarity      Rarity
    Effect      JokerEffect
    Value       int   // Effect magnitude

    // State (for stateful jokers)
    Counter     int
}

type Rarity int
const (
    Common Rarity = iota
    Uncommon
    Rare
    Legendary
)

// Example jokers (start with 5 simple ones)
var startingJokers = []Joker{
    {
        ID: "joker_1",
        Name: "Joker",
        Description: "+4 Mult",
        Rarity: Common,
        Effect: AddMult,
        Value: 4,
    },
    {
        ID: "greedy_joker",
        Name: "Greedy Joker",
        Description: "Played cards with Diamond suit give +3 Mult when scored",
        Rarity: Common,
        Effect: MultIfSuit,
        Value: 3,
    },
    {
        ID: "lusty_joker",
        Name: "Lusty Joker",
        Description: "Played cards with Heart suit give +3 Mult when scored",
        Rarity: Common,
        Effect: MultIfSuit,
        Value: 3,
    },
    // ... more
}

func (j Joker) Apply(calc *ScoreCalculation, hand HandInfo, played []Card) {
    switch j.Effect {
    case AddMult:
        calc.TotalMult += j.Value
        calc.MultModifiers = append(calc.MultModifiers,
            fmt.Sprintf("+%d (%s)", j.Value, j.Name))

    case MultIfPair:
        if hand.Type >= Pair {
            calc.TotalMult *= j.Value
            calc.MultModifiers = append(calc.MultModifiers,
                fmt.Sprintf("x%d (%s)", j.Value, j.Name))
        }

    // ... more effect types
    }
}
```

#### 3.2 Card Enhancements (~200 LOC)
**File**: `cards.go` (extend)

```go
func (c Card) GetChipValue() int {
    baseValue := c.Rank.ChipValue()  // Face cards = 10, Ace = 11, etc.

    switch c.Enhancement {
    case Bonus:
        return baseValue + 30
    case Steel:
        return baseValue + 50
    case Stone:
        return 50  // Fixed value, no rank
    default:
        return baseValue
    }
}

func (c Card) ApplyMultEffects(calc *ScoreCalculation) {
    switch c.Enhancement {
    case Mult:
        calc.TotalMult += 4
        calc.MultModifiers = append(calc.MultModifiers, "+4 (Mult card)")

    case Glass:
        calc.TotalMult *= 2
        calc.MultModifiers = append(calc.MultModifiers, "x2 (Glass)")
        // NOTE: Card destroys after scoring (handle in game logic)
    }

    // Edition effects
    switch c.Edition {
    case Foil:
        calc.TotalChips += 50
        calc.ChipModifiers = append(calc.ChipModifiers, "+50 (Foil)")

    case Holographic:
        calc.TotalMult += 10
        calc.MultModifiers = append(calc.MultModifiers, "+10 (Holo)")

    case Polychrome:
        calc.TotalMult = int(float64(calc.TotalMult) * 1.5)
        calc.MultModifiers = append(calc.MultModifiers, "x1.5 (Poly)")
    }
}
```

#### 3.3 Enhanced Scoring Engine (~300 LOC)
**File**: `scoring.go` (rewrite)

```go
func CalculateScore(hand HandInfo, playedCards []Card, jokers []Joker) ScoreCalculation {
    calc := ScoreCalculation{
        BaseChips: hand.BaseChips,
        BaseMult:  hand.BaseMult,
        TotalChips: hand.BaseChips,
        TotalMult:  hand.BaseMult,
    }

    // Step 1: Add card chip values
    for _, card := range playedCards {
        chips := card.GetChipValue()
        calc.TotalChips += chips
        calc.ChipModifiers = append(calc.ChipModifiers,
            fmt.Sprintf("+%d (%s)", chips, card.Rank))
    }

    // Step 2: Apply card mult effects (Mult enhancement, editions)
    for _, card := range playedCards {
        card.ApplyMultEffects(&calc)
    }

    // Step 3: Apply joker effects IN ORDER
    for _, joker := range jokers {
        joker.Apply(&calc, hand, playedCards)
    }

    // Step 4: Calculate final score
    calc.FinalScore = calc.TotalChips * calc.TotalMult

    return calc
}
```

**Deliverable**: Jokers modify scores, card enhancements work, scoring breakdown shows all steps.

---

## Phase 4: Shop & Economy - ~500 LOC

### Week 4: Progression System

#### 4.1 Shop System (~300 LOC)
**File**: `shop.go`

```go
type ShopItem interface {
    GetName() string
    GetDescription() string
    GetCost() int
    GetType() ShopItemType
}

type ShopItemType int
const (
    ItemJoker ShopItemType = iota
    ItemPack
    ItemVoucher
)

type ShopState struct {
    Jokers    []Joker       // 2 jokers for sale
    Packs     []Pack        // 2 booster packs
    Vouchers  []Voucher     // 1-2 vouchers
    RerollCost int          // Increases each reroll
}

func GenerateShop(ante int, money int) ShopState {
    shop := ShopState{
        RerollCost: 5,
    }

    // Generate 2 random jokers
    shop.Jokers = []Joker{
        getRandomJoker(Common, Uncommon),
        getRandomJoker(Common, Rare),
    }

    // Generate 2 packs
    shop.Packs = []Pack{
        {Type: ArcanaPack, Cost: 4},  // Tarot cards
        {Type: CelestialPack, Cost: 4}, // Planet cards
    }

    return shop
}

func (m *Model) BuyJoker(joker Joker) error {
    if m.round.TotalMoney < joker.GetCost() {
        return errors.New("not enough money")
    }

    if len(m.jokers) >= 5 {
        return errors.New("joker slots full")
    }

    m.jokers = append(m.jokers, joker)
    m.round.TotalMoney -= joker.GetCost()
    return nil
}
```

#### 4.2 Shop UI (~200 LOC)
**File**: `view.go` (extend)

```go
func (m Model) viewShop() string {
    var b strings.Builder

    b.WriteString("╔════════════════════════════════════════╗\n")
    b.WriteString("║              SHOP                      ║\n")
    b.WriteString(fmt.Sprintf("║  Money: $%d                           ║\n", m.round.TotalMoney))
    b.WriteString("╚════════════════════════════════════════╝\n\n")

    // Jokers for sale
    b.WriteString("JOKERS:\n")
    for i, joker := range m.shop.Jokers {
        selected := (m.cursorPosition == i)
        b.WriteString(m.renderShopJoker(joker, selected))
    }

    b.WriteString("\n\nPACKS:\n")
    for i, pack := range m.shop.Packs {
        selected := (m.cursorPosition == 2+i)
        b.WriteString(m.renderShopPack(pack, selected))
    }

    b.WriteString("\n\n[Enter] Buy | [R] Reroll ($5) | [Space] Continue")

    return b.String()
}

func (m Model) renderShopJoker(joker Joker, selected bool) string {
    style := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        Padding(1).
        Width(20)

    if selected {
        style = style.BorderForeground(lipgloss.Color("226"))
    }

    content := fmt.Sprintf("%s\n%s\n\nCost: $%d",
        joker.Name,
        joker.Description,
        joker.GetCost())

    return style.Render(content)
}
```

**Deliverable**: Can buy jokers, packs, reroll shop, manage money.

---

## Phase 5: Polish & Animations - ~800 LOC

### Week 5: Make It Feel Good

#### 5.1 Scoring Animation (~300 LOC)
**File**: `animations.go`

```go
type ScoringAnimation struct {
    Step            int
    MaxSteps        int
    CurrentChips    int
    TargetChips     int
    CurrentMult     int
    TargetMult      int
    CurrentScore    int
    TargetScore     int
    ModifierIndex   int  // Which modifier to show
}

func (m *Model) StartScoringAnimation(calc ScoreCalculation) {
    m.animating = true
    m.scoreAnim = ScoringAnimation{
        Step: 0,
        MaxSteps: 120,  // 2 seconds at 60fps
        TargetChips: calc.TotalChips,
        TargetMult: calc.TotalMult,
        TargetScore: calc.FinalScore,
    }
}

func (m *Model) ProgressScoringAnimation() {
    m.scoreAnim.Step++

    // Animate chip counter
    if m.scoreAnim.CurrentChips < m.scoreAnim.TargetChips {
        increment := max(1, (m.scoreAnim.TargetChips - m.scoreAnim.CurrentChips) / 10)
        m.scoreAnim.CurrentChips = min(
            m.scoreAnim.CurrentChips + increment,
            m.scoreAnim.TargetChips,
        )
    }

    // Animate mult counter
    if m.scoreAnim.CurrentMult < m.scoreAnim.TargetMult {
        increment := max(1, (m.scoreAnim.TargetMult - m.scoreAnim.CurrentMult) / 10)
        m.scoreAnim.CurrentMult = min(
            m.scoreAnim.CurrentMult + increment,
            m.scoreAnim.TargetMult,
        )
    }

    // Show modifiers one by one
    if m.scoreAnim.Step % 10 == 0 {
        m.scoreAnim.ModifierIndex++
    }

    // Calculate score
    m.scoreAnim.CurrentScore = m.scoreAnim.CurrentChips * m.scoreAnim.CurrentMult

    // Done?
    if m.scoreAnim.Step >= m.scoreAnim.MaxSteps {
        m.animating = false
        m.state = StateSelectCards  // Return to game
    }
}

func (m Model) viewScoring() string {
    var b strings.Builder

    b.WriteString("SCORING...\n\n")

    // Show poker hand
    b.WriteString(fmt.Sprintf("Hand: %s\n\n", m.lastHand.Type))

    // Show chip calculation
    b.WriteString(fmt.Sprintf("Chips: %d", m.scoreAnim.CurrentChips))
    for i := 0; i < m.scoreAnim.ModifierIndex && i < len(m.scoreCalc.ChipModifiers); i++ {
        b.WriteString(fmt.Sprintf(" %s", m.scoreCalc.ChipModifiers[i]))
    }
    b.WriteString("\n\n")

    // Show mult calculation
    b.WriteString(fmt.Sprintf("Mult: x%d", m.scoreAnim.CurrentMult))
    for i := 0; i < m.scoreAnim.ModifierIndex && i < len(m.scoreCalc.MultModifiers); i++ {
        b.WriteString(fmt.Sprintf(" %s", m.scoreCalc.MultModifiers[i]))
    }
    b.WriteString("\n\n")

    // Show final score (BIG)
    b.WriteString(fmt.Sprintf("SCORE: %d\n", m.scoreAnim.CurrentScore))

    return b.String()
}
```

#### 5.2 UI Polish (~300 LOC)
- Card selection highlighting (gold border from solitaire)
- Joker panel (3-5 visible at bottom)
- Info panel updates
- Responsive layout (vertical stack on narrow terminals - reuse solitaire pattern!)
- Smooth transitions between states

#### 5.3 Enhanced Input (~200 LOC)
- Mouse click to select cards (from solitaire)
- Drag to select multiple cards (maybe)
- Number keys 1-8 to toggle cards
- D to discard, P to play
- Shop navigation with mouse

**Deliverable**: Polished, smooth gameplay experience.

---

## Phase 6: Content Expansion - ~500 LOC

### Beyond Week 5: Add Depth

#### 6.1 More Jokers (~200 LOC)
Implement 20-30 jokers total with varied effects:
- Simple: +chips, +mult
- Conditional: based on hand type, suit, rank
- Multiplicative: x2 if condition
- Retrigger: cards score multiple times
- Money: earn $$ on triggers
- Card creation: add cards to deck

#### 6.2 Planet Cards (~150 LOC)
Level up specific poker hands:
```go
type PlanetCard struct {
    Name        string
    HandType    PokerHand
    ChipsBonus  int  // +25 chips per level
    MultBonus   int  // +2 mult per level
}

func (p PlanetCard) Apply(hand *HandInfo) {
    if hand.Type == p.HandType {
        hand.BaseChips += p.ChipsBonus * hand.Level
        hand.BaseMult += p.MultBonus * hand.Level
        hand.Level++
    }
}
```

#### 6.3 Tarot Cards (~150 LOC)
Instant effects on cards:
- The Magician: Add random enhancement to card
- The Empress: Add random edition to card
- The Hanged Man: Destroy up to 2 selected cards
- The Devil: Add enhancement to selected card

**Deliverable**: Deep, replayable game with tons of variety.

---

## Technical Debt & Future Work

### Known Limitations (MVP)
- No save/load system
- No statistics tracking
- No achievements
- No seed system for runs
- Limited boss blind effects

### Post-MVP Features
- Save/load via JSON
- Statistics (best run, total hands played, etc.)
- Daily challenge with fixed seed
- More boss blinds (20+ with unique effects)
- Vouchers (permanent upgrades)
- Deck variants (different starting conditions)
- Stakes (difficulty modifiers)

---

## Success Criteria

### MVP (End of Week 2) ✅ COMPLETE
- ✅ Can play poker hands
- ✅ Blinds and progression work
- ✅ Scoring is correct
- ✅ Game over on loss
- ✅ Basic UI

### Playable (End of Week 4) ✅ COMPLETE
- ✅ Jokers modify gameplay (15+ jokers implemented)
- ✅ Shop works (buy jokers, manage money)
- ✅ Can beat multiple antes
- ✅ Economy system functional

### Polished (End of Week 5) 🚧 IN PROGRESS
- ⏳ Animations smooth (scoring animation needed)
- ✅ UI responsive and beautiful
- ✅ Mouse + keyboard input (COMPLETE - both card selection and shop)
- ✅ 15+ jokers implemented
- ✅ Color-coded UI elements
- ✅ Stable layout (no UI jumping)

### Feature Complete (Week 6+)
- ✅ 20-30 jokers
- ✅ Planet and tarot cards
- ✅ Boss blind effects
- ✅ Multiple deck types
- ✅ Save/load

---

## External Resources

**NEW!** Complete reference document with all game data:
- **`docs/balatro/REFERENCES.md`** - Comprehensive resource guide
  - Links to Balatro Calculator (150 joker definitions!)
  - All card types, effects, and mechanics
  - Boss blind effects
  - Tarot, Planet, and Spectral cards
  - Implementation roadmap

**Balatro Calculator Repository**: https://github.com/EFHIII/balatro-calculator
- 150 joker definitions with full mechanics
- Complete scoring engine implementation
- Card enhancement/edition/seal formulas
- Use for reference when implementing new jokers!

---

## Next Steps

### Immediate Priorities (Next Session)
1. **Scoring Animation** (~300 LOC) - Add smooth scoring animations like Phase 5.1
   - Chip counter animation
   - Mult counter animation
   - Progressive modifier reveal
   - Final score calculation display

2. **Boss Blind Effects** (~200 LOC) - Implement special boss blind mechanics
   - "The Hook" - Disables 2 random Jokers
   - "The Ox" - Sets money to $0 after beating blind
   - "The House" - First hand drawn face down
   - "The Plant" - All face cards debuffed

3. **More Jokers** (~300 LOC) - Expand to 30+ jokers
   - Retrigger effects
   - Money-generating jokers
   - Card-modifying jokers
   - Conditional multipliers

### Medium Term (Future Sessions)
4. **Planet Cards** (~150 LOC) - Level up poker hands
5. **Tarot Cards** (~150 LOC) - Card modification system
6. **Booster Packs** (~200 LOC) - Pack opening system
7. **Vouchers** (~150 LOC) - Permanent upgrades

### Long Term (Post-MVP)
8. **Save/Load System** - Persist game state
9. **Statistics Tracking** - Best runs, total hands played
10. **Daily Challenge** - Seeded runs

---

## Current Status Summary

**What's Working**: ✅
- Complete poker hand evaluation
- Full blind progression system (Small → Big → Boss)
- 15+ jokers with varied effects
- Card enhancements (Bonus, Mult, Glass, Steel, Stone, Gold, Lucky)
- Card editions (Foil, Holographic, Polychrome)
- Card seals system
- Shop with purchasing
- Mouse + keyboard input for all interactions
- Responsive, color-coded UI

**What's Next**: 🚧
- Scoring animations
- Boss blind special effects
- More jokers (30+ target)
- Planet and Tarot cards
- Polish and balance

**Game is Playable**: You can play through multiple antes, buy jokers, and experience the core Balatro gameplay loop! 🎉

Let's keep building! 🚀

# Balatro TUI - External References & Resources

This document catalogs external resources used for implementing accurate Balatro game mechanics.

---

## Primary Reference: Balatro Calculator

**Repository**: https://github.com/EFHIII/balatro-calculator
**License**: MIT
**Stars**: 239+ ⭐
**Language**: JavaScript (85.1%), HTML (9.6%), CSS (5.3%)

### What It Contains

A complete Balatro scoring calculator with:
- **150 Joker definitions** with full mechanics (`cards.js`)
- **Scoring simulation engine** (`balatro-sim.js`)
- **Card enhancements, editions, and seals** with exact formulas
- **Hand evaluation logic**
- **Score breakdown calculations**

### Key Files for Reference

| File | Purpose | Use For |
|------|---------|---------|
| `cards.js` | All 150 joker definitions | Implementing joker effects, costs, rarities |
| `balatro-sim.js` | Core scoring engine | Score calculation algorithms, card mechanics |
| `breakdown.js` | Score display logic | Formatting score breakdowns |

### Limitations

Per the README: *"Currently, a couple of the jokers may not work correctly, but most of them do."*

Most jokers (~95%+) are accurately implemented, making this an excellent reference.

---

## Data Structures

### Joker Properties (from cards.js)

Each joker has:
```javascript
{
  order: 1-150,              // Sequential ID
  name: "Joker Name",        // Display name
  rarity: 0-3,               // 0=Common, 1=Uncommon, 2=Rare, 3=Legendary
  cost: 5-20,                // Shop price ($)
  pos: {x, y},               // Sprite grid position
  config: {...},             // Effect values (mult, chips, extra, etc.)
  blueprint_compat: true,    // Can be copied by Blueprint joker
  perishable_compat: true,   // Can have Perishable edition
  eternal_compat: true,      // Can have Eternal edition
  unlock_condition: {...}    // Unlock requirements
}
```

### Joker Rarities & Costs

| Rarity | Cost | Count (approx) |
|--------|------|----------------|
| Common (0) | $5 | ~60 jokers |
| Uncommon (1) | $7 | ~50 jokers |
| Rare (2) | $10 | ~30 jokers |
| Legendary (3) | $20 | ~10 jokers |

---

## Card Mechanics Reference

### Enhancements

| Enhancement | Effect | Implementation Status |
|-------------|--------|----------------------|
| **Bonus** | +30 chips | ✅ Implemented |
| **Mult** | +4 mult | ✅ Implemented |
| **Glass** | ×2 mult (destroys after scoring) | ✅ Implemented |
| **Steel** | ×1.5 mult when in hand (unplayed cards) | ✅ Implemented |
| **Stone** | +50 chips, no rank bonus | ✅ Implemented |
| **Gold** | $3 at end of round | ✅ Implemented |
| **Lucky** | 1 in 5 chance: +20 mult and $1 | ✅ Implemented |

### Editions

| Edition | Effect | Implementation Status |
|---------|--------|----------------------|
| **Foil** | +50 chips | ✅ Implemented |
| **Holographic** | +10 mult | ✅ Implemented |
| **Polychrome** | ×1.5 mult | ✅ Implemented |

### Seals

| Seal | Effect | Implementation Status |
|------|--------|----------------------|
| **Red Seal** | Retriggers card (scores twice) | ✅ Implemented |
| **Blue Seal** | Creates Planet card when held at end of round | ✅ Implemented |
| **Gold Seal** | $3 when played | ✅ Implemented |
| **Purple Seal** | Creates Tarot card when discarded | ✅ Implemented |

---

## Joker Categories (150 Total)

### Basic Stat Boost (~15 jokers)
- **Joker** (#1): +4 mult
- **Greedy/Lusty/Wrathful/Gluttonous Jokers** (#2-5): +3 mult per suit card
- **Jolly Joker** (#6): +8 mult if hand contains Pair
- **Zany Joker** (#7): +12 mult if hand contains Three of a Kind
- **Mad Joker** (#8): +10 mult if hand contains Two Pair
- **Crazy Joker** (#9): +12 mult if hand contains Straight
- **Droll Joker** (#10): +10 mult if hand contains Flush

### Hand Size & Discards (~5 jokers)
- **Juggler** (#87): +1 hand size
- **Drunkard** (#88): +1 discard
- **Troubadour** (#89): +2 hand size, -1 discard
- **Certificate** (#90): +1 hand size during boss blinds
- **Smeared Joker** (#91): Hearts/Diamonds count as same suit, Spades/Clubs count as same suit

### Scaling Jokers (~30 jokers)
Examples:
- **Gros Michel** (#43): +15 mult, 1 in 6 chance to destroy after round
- **Cavendish** (#99): +3 mult, 1 in 1000 chance to destroy, appears after Gros Michel destroyed
- **Steel Joker** (#45): +0.25 mult per Steel card in deck
- **Fibonacci** (#42): Each Ace, 2, 3, 5, or 8 gives +8 mult
- **Scary Face** (#41): Each face card gives +30 chips

### Money Jokers (~20 jokers)
Examples:
- **Golden Ticket** (#100): Gives $4 per $Gold Seal$ card played
- **Mr. Bones** (#101): Prevents death once, destroys self
- **Acrobat** (#102): +3 mult, final hand gives ×3 mult
- **Sock and Buskin** (#103): Retriggers all face cards

### Retrigger Jokers (~10 jokers)
- **Sock and Buskin** (#103): Retriggers face cards
- **Hanging Chad** (#89): Retriggers first played card 2 additional times
- **Dusk** (#95): Retriggers all played cards in final hand

### Conditional Jokers (~40 jokers)
Based on:
- Played suits
- Card ranks
- Hand types
- Poker hand levels
- Cards remaining in deck
- Round number
- Money amount

### Special Mechanic Jokers (~30 jokers)
Examples:
- **Blueprint** (#135): Copies adjacent joker
- **Brainstorm** (#136): Copies leftmost joker
- **Satellite** (#137): $1 per unique Planet card used
- **Shoot the Moon** (#138): +13 mult per Queen in hand
- **Driver's License** (#139): ×3 mult if deck has ≥16 enhanced cards

---

## Tarot Cards (22 Total)

From standard tarot deck - modify cards and deck:
- **The Fool**: Creates last played hand as Planet card
- **The Magician**: Enhances selected cards with random enhancement
- **The High Priestess**: Creates up to 2 random Planet cards
- **The Empress**: Enhances selected cards with random edition
- **The Emperor**: Creates up to 2 random Tarot cards
- **The Hierophant**: Creates up to 2 random Planet cards
- **The Lovers**: Enhances 1 selected card to Gold
- **The Chariot**: Enhances 1 selected card to Steel
- **Justice**: Enhances 1 selected card to Glass
- **The Hermit**: Doubles money (max $20)
- **Wheel of Fortune**: 1 in 4 chance to add Foil/Holographic/Polychrome to random card
- **Strength**: Increases rank of up to 2 selected cards by 1
- **The Hanged Man**: Destroys up to 2 selected cards
- **Death**: Converts 2 selected cards to random different cards
- **Temperance**: Sell value of jokers increased by $50 total
- **The Devil**: Enhances 1 selected card to Gold
- **The Tower**: Enhances 1 selected card to Stone
- **The Star**: Converts up to 3 selected cards to Diamonds
- **The Moon**: Converts up to 3 selected cards to Clubs
- **The Sun**: Converts up to 3 selected cards to Hearts
- **Judgement**: Creates random Joker card (free)
- **The World**: Converts up to 3 selected cards to Spades

---

## Planet Cards (11 Total)

Level up specific poker hands (increase base chips & mult):

| Planet | Poker Hand | Base Upgrade |
|--------|------------|--------------|
| **Mercury** | Pair | +15 chips, +1 mult |
| **Venus** | Three of a Kind | +20 chips, +2 mult |
| **Earth** | Full House | +25 chips, +2 mult |
| **Mars** | Flush | +30 chips, +3 mult |
| **Jupiter** | Flush Five | +40 chips, +3 mult |
| **Saturn** | Straight Flush | +35 chips, +3 mult |
| **Uranus** | Royal Flush | +40 chips, +4 mult |
| **Neptune** | Straight | +30 chips, +3 mult |
| **Pluto** | High Card | +10 chips, +1 mult |
| **Planet X** | Four of a Kind | +30 chips, +3 mult |
| **Ceres** | Two Pair | +20 chips, +2 mult |

Each use increases hand level, compounding upgrades.

---

## Spectral Cards (19 Total)

High-risk, high-reward cards with unique effects:
- **Familiar**: Destroy 1 random card, add 3 random Enhanced cards to deck
- **Grim**: Destroy 1 random card, add 2 random Enhanced cards to hand
- **Incantation**: Destroy 1 random card, add 4 random Enhanced cards to deck
- **Talisman**: Add Gold Seal to 1 random card in hand
- **Aura**: Add Foil/Holographic/Polychrome to 1 random card
- **Wraith**: Creates random Rare Joker, sets money to $0
- **Sigil**: Converts all cards in hand to single random suit
- **Ouija**: Converts all cards in hand to single random rank
- **Ectoplasm**: Add Negative to random Joker, sets hand size to -1
- **Immolate**: Destroys 5 random cards, gains $20
- **Ankh**: Create copy of random Joker, destroy all other Jokers
- **Deja Vu**: Add Red Seal to 1 random card
- **Hex**: Add Polychrome to random Joker, destroy all other Jokers
- **Trance**: Add Blue Seal to 1 random card
- **Medium**: Add Purple Seal to 1 random card
- **Cryptid**: Create 2 copies of 1 random card
- **The Soul**: Creates random Legendary Joker
- **Black Hole**: Upgrade every poker hand by 1 level
- **White Hole**: Select 1 hand type, reset all levels, gain ×3 to selected hand's multiplier

---

## Boss Blinds (15+ Unique)

Special conditions that modify gameplay:

| Boss Blind | Effect |
|------------|--------|
| **The Hook** | 2 random jokers disabled |
| **The Ox** | Sets money to $0 after defeating |
| **The House** | First hand drawn face down |
| **The Wall** | Extra large blind |
| **The Wheel** | 1 in 7 cards drawn face down |
| **The Arm** | Decrease level of played hand by 1 |
| **The Club** | All Clubs debuffed |
| **The Fish** | All cards drawn face down until 1 hand played |
| **The Psychic** | Must play 5 cards |
| **The Goad** | All Spades debuffed |
| **The Water** | Start with 0 discards |
| **The Window** | All Diamonds debuffed |
| **The Manacle** | Start with 1 hand |
| **The Eye** | No repeating hand types allowed |
| **The Mouth** | Play only 1 hand type this blind |
| **The Plant** | All face cards debuffed |
| **The Serpent** | After playing hand, always draw 3 cards |
| **The Pillar** | Cards played previously this ante are debuffed |
| **The Needle** | Play only 1 hand this blind |
| **The Head** | All Hearts debuffed |
| **The Tooth** | Lose $1 per card played |
| **The Flint** | Base chips and mult halved |
| **The Mark** | All face cards drawn face down |

---

## Implementation Roadmap

### Phase 3: ✅ Complete
- 5 basic jokers implemented
- Card enhancements working
- Card editions working
- Card seals working

### Phase 4: Next Steps (~800 LOC)
**Priority 1 - More Jokers** (~300 LOC):
- Implement 10-15 additional common/uncommon jokers
- Focus on simple stat boosts and suit-based bonuses
- Target: 20+ total jokers

**Priority 2 - Boss Blinds** (~200 LOC):
- Implement 5-10 boss blind special effects
- Start with simpler ones (The Hook, The Ox, The House)

**Priority 3 - Planet Cards** (~150 LOC):
- Implement hand leveling system
- Create 11 planet cards
- Level-up UI and persistence

**Priority 4 - Tarot Cards** (~150 LOC):
- Implement card modification system
- Create 10-15 most common tarots
- Card selection UI

### Phase 5: Advanced (~500 LOC)
- Spectral cards
- Remaining jokers (scale to 50+)
- Blueprint/Brainstorm copy mechanics
- Negative/Eternal/Perishable editions

---

## Using This Reference

### For Implementing Jokers
1. Browse `cards.js` in the calculator repo for joker definitions
2. Check `balatro-sim.js` for exact effect implementations
3. Cross-reference with your `games/balatro/jokers.go`
4. Implement in order of complexity (simple stat boosts → conditional → scaling)

### For Scoring Mechanics
1. Reference `balatro-sim.js` for calculation order:
   - Base hand chips/mult
   - Card chip values
   - Card mult effects (enhancements, editions)
   - Joker effects (in order owned)
   - Final calculation: chips × mult

2. Use `breakdown.js` for formatting score displays

### Testing Accuracy
Use the live calculator at the repo's GitHub Pages to verify:
- Hand evaluations match
- Joker effects calculate correctly
- Edge cases behave as expected

---

## Additional Resources

### Official Game
- **Balatro** on Steam: https://store.steampowered.com/app/2379780/Balatro/
- Developer: LocalThunk
- Publisher: Playstack

### Community Resources
- Balatro Wiki (community-maintained)
- Balatro Discord server
- r/balatro subreddit

---

## Credits

**Balatro Calculator**: Created by EFHIII (https://github.com/EFHIII)
**Balatro Game**: Created by LocalThunk
**This Reference**: Compiled for TUI Classics Balatro implementation

---

**Last Updated**: 2025-01-25
**Current Implementation**: Phase 3 Complete (5 jokers, all card mechanics)
**Next Target**: Phase 4 (20+ jokers, boss blinds, planets)

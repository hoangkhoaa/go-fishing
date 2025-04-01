# Terminal Fishing Game

A relaxing terminal-based fishing game written in Go using the Bubble Tea TUI framework.

## Features

- Catch 50 different fish species with unique characteristics
- Fish have varying rarity, weight, color, and pattern
- Weather conditions affect fishing success
- Time of day affects fishing success and fish availability
- Auto-fishing feature that works in the background (every 10 seconds)
- Attractive ASCII art and animations
- Simple keyboard-driven interface

## Installation

### Prerequisites

- Go 1.18 or higher

### Building from source

```bash
# Clone the repository
git clone https://github.com/yourusername/fishing-game.git
cd fishing-game

# Build the game
make build

# Run the game
make run
```

## How to Play

- Use arrow keys or j/k to navigate menus
- Press Enter or Space to select menu options
- Press 'a' to toggle auto-fishing from anywhere
- Press q or Esc to go back or exit

### Game Options

- **Go Fishing**: Cast your line and try to catch fish
- **View Inventory**: See what fish you've caught so far
- **Quit Game**: Exit the game

### Command-Line Options

- **Normal Mode**: `./fishing-game` - Standard gameplay with fishing times between 5 seconds and 5 minutes
- **Test Mode**: `./fishing-game -test` - Shortened fishing durations (5-10 seconds) for easier testing

### Auto-Fishing

You can toggle auto-fishing by pressing 'a' at any time. When enabled:
- The game will automatically attempt to catch fish every 10 seconds
- Each fish has different catch probabilities based on rarity
- Auto-fishing works even when you're in other menus
- You can see your auto-fishing status in the stats bar

### Fish Collection

The game features 60 different catchable items with varying:
- Rarity (common to extremely rare)
- Weight (small to massive)
- Colors and patterns
- Habitats
- Time preferences (each fish is more active at certain times of day)

#### Legendary and Mythical Creatures

Beyond regular fish, you may occasionally encounter legendary and mythical creatures:
- **Extremely Rare**: Only a 0.5-2% chance depending on conditions
- **Highly Valuable**: Worth much more than regular fish
- **Time-Specific**: Some legendary creatures appear only at certain times
- Special creatures include:
  - The Kraken (night)
  - Loch Ness Monster (night)
  - Megalodon (night)
  - Golden Carp (morning)
  - Phoenix Fish (afternoon)
  - And more!

#### Trash Items

Not every catch will be a fish! Sometimes you'll pull up various trash items:
- About 12% chance to catch trash instead of a fish
- Most trash items have no value, but some may be worth a small amount
- Examples include old boots, plastic bottles, shopping bags, and more
- The rare Treasure Chest is technically trash but quite valuable!

### Time of Day System

The game uses real-time to simulate different fishing periods throughout the day:
- **Morning (5am-11am)** - Perfect for surface feeders, 20% boost to catch rates
- **Afternoon (11am-5pm)** - Deep water fish are more active, but overall slower fishing
- **Evening (5pm-9pm)** - Prime fishing time with largest catch rate boost (30%)
- **Night (9pm-5am)** - Good for nocturnal species, especially dark-colored fish

Each time period has different effects:
- Changes which fish are more likely to appear
- Adjusts overall catch success rates
- Provides visual indicators in the game interface
- Shows current time period and its effect on fishing

### Inventory Management

The inventory screen provides a complete overview of your catches:
- See all fish types you've caught with quantities, weights, and values
- Total catch count, total weight, and total value are displayed
- Sort your inventory in different ways:
  - Press `1` to sort by rarity (rarest fish first)
  - Press `2` to sort by weight (heaviest fish first)
  - Press `3` to sort by value (most valuable fish first)
- The currently sorted column is highlighted for easy reference

## Development

The code is organized in a modular structure:
- `main.go` - Program entry point and initialization
- `model.go` - Core data structures and UI styling
- `fishing.go` - Fishing mechanics and logic
- `graphics.go` - Visual elements and fish patterns
- `views.go` - UI rendering for different game states
- `background.go` - Background processes for idle fishing

## License

MIT License 

## Future Roadmap: Real-World News Integration

In the future, we plan to integrate real-world news into the game as "Random Encounters" that will dynamically affect gameplay. This system will connect fishing mechanics to external events for a more engaging and varied experience.

### News-Based Random Encounters

The fishing game will fetch real-world data from various sources to create in-game events:

#### Resource Markets
- **Oil Price Fluctuations**: Oil spills during price crashes could reduce fishing success rates in certain areas
- **Gold Market**: Gold rushes during price spikes might introduce gold nugget fishing opportunities
- **Crypto Markets**: Crypto booms might allow discovery of waterproof hard drives containing digital currency

#### Weather Systems
- **Real Weather Data**: Local or global weather patterns will affect fishing conditions
- **Storm Warnings**: Major storms could create opportunities to catch rare storm-driven species
- **Temperature Changes**: Seasonal fish migrations based on real climate data

#### Global Events
- **Conflict Zones**: War news might temporarily block certain fishing areas or introduce military debris
- **Trade Agreements**: New fishing territories could open based on international relations
- **Conservation News**: Endangered species events based on environmental news

### Technical Implementation Plan

1. **Phase 1: News API Integration**
   - Connect to financial, weather, and news APIs
   - Develop a news processing system that extracts relevant game-affecting data

2. **Phase 2: Event System Development**
   - Create an event generation engine that converts news into gameplay effects
   - Implement a quest/encounter system for special time-limited events

3. **Phase 3: Dynamic World Adaptation**
   - Develop changing fish populations based on real-world conditions
   - Implement seasonal and news-reactive fishing locations

4. **Phase 4: Player Economy**
   - Add trading system affected by resource market news
   - Implement fish market price fluctuations based on real economic data

This roadmap represents our vision for evolving the game into a dynamically changing world that reflects real-world events while maintaining the relaxing core fishing experience. 
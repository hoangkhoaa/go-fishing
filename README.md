# 🎣 Terminal Fishing Game

Hey there! Welcome to my little passion project - a cozy terminal-based fishing game that lets you relax and catch some digital fish while multitasking with your other work. I built this as a solo developer using Go and the Bubble Tea framework, hoping to create something that brings a bit of joy to your terminal.

## ✨ What Makes This Special

- 🐟 Catch 50 different fish species, each with their own personality (rarity, weight, colors)
- 🌦️ Real-time weather affects your fishing luck
- 🕒 Time of day changes which fish are active (morning, afternoon, evening, night)
- 🤖 Auto-fishing lets you catch fish in the background while you do other things
- 🎨 Charming ASCII art and animations to brighten your terminal
- 🎮 Super simple keyboard controls - nothing complicated here!

## 🚀 Getting Started

### What You'll Need

- Go 1.18 or newer

### Quick Setup

```bash
# Grab a copy of the code
git clone https://github.com/hoangkhoaa/go-fishing.git
cd go-fishing

# Build it
make build

# Let's go fishing!
make run
```

## 🎮 How to Play

Nothing complicated here - just a few simple controls:

- Move around with arrow keys (or j/k if you're a keyboard wizard)
- Select stuff with Enter or Space
- Press 'a' anytime to toggle auto-fishing (this is the best part!)
- Need to escape? Press q or Esc

### Game Options

- **Go Fishing**: Throw in your line and see what bites
- **View Inventory**: Check out your fishy collection
- **View History**: See how your fishing has gone over time
- **Quit Game**: Take a break (but come back soon!)

### Command-Line Options

- **Chill Mode**: `./fishing-game` - Normal fishing times (10-120 seconds)
- **Impatient Mode**: `./fishing-game -test` - Quick fishing (5-10 seconds) for when you just want to catch 'em all

### 🤖 Auto-Fishing - Fish While You Work!

This is my favorite feature! Press 'a' anytime to let the game fish for you:
- It'll automatically try for a catch every 10 seconds
- Works while you're doing other things in the game
- Even works when you're in another tab doing actual work!
- The status bar shows if it's active

### 🐟 Fish Collection

I've added 60 different catchable items to discover:
- From common minnows to ultra-rare legendary creatures
- Weights ranging from tiny to massive
- Various colors and patterns to collect
- Different habitats and time preferences

#### 🌟 Legendary and Mythical Creatures

If you're lucky, you might encounter something extraordinary:
- Super rare (0.5-2% chance)
- Worth a small fortune
- Some only appear at specific times of day
- Keep an eye out for The Kraken, Loch Ness Monster, Golden Carp, and more!

#### 🗑️ Not-So-Treasures

Sometimes you'll catch... well, junk:
- About a 12% chance to reel in something that's not a fish
- Old boots, plastic bottles, the usual suspects
- But don't ignore them - there's a rare Treasure Chest hiding among the trash!

### ⏰ Time of Day Affects Your Fishing

I added a time system that uses your computer's real time:
- **Morning (5am-11am)** - Early birds get 20% better catch rates
- **Afternoon (11am-5pm)** - Things slow down in the midday heat
- **Evening (5pm-9pm)** - Prime fishing time! 30% boost to catch rates
- **Night (9pm-5am)** - Perfect for night owls hunting nocturnal species

The game shows you which time period you're in and how it affects fishing.

### 🧾 Inventory and History Features

Keep track of everything you've caught:
- Sort your catches by rarity (1), weight (2), value (3), or quantity (4)
- Press 'h' to see your fishing history organized by date
- The game automatically saves your progress - no "save game" needed!

## 💾 Save System

- Your catches are saved by date in the `saves` directory
- Game automatically saves when you catch something or every 5 minutes
- No need to worry about losing your progress!

## 🔮 What's Coming Next: News Integration

I'm working on something cool - integrating real-world news into the game to create dynamic events:

- Oil spills, gold rushes, and crypto booms affecting fishing
- Real weather patterns influencing what you can catch
- Global events creating special limited-time opportunities

This is still in development, but I'm excited to bring the real world into our little fishing game!

## 🤝 Want to Help?

This is an open source passion project. If you'd like to contribute or just say hi, check out the CONTRIBUTING.md file or open an issue. I'd love to hear from you!

## 📄 License

MIT License - Feel free to play with the code! 
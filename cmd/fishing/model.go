package main

import (
	"math/rand"
	"time" // Add time package for auto-continue

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/user/fishing-game/game"
)

// Model for bubbletea
type model struct {
	state           string
	menuItems       []string
	selectedItem    int
	fishingState    int
	message         string
	catchSuccess    bool
	caughtFish      game.Fish
	fishShape       string
	width           int
	height          int
	autoFishMsg     string
	autoFishTick    int
	resultTimer     int     // Track time in fish result screen
	fishingDuration int64   // Total duration for the fishing attempt in milliseconds
	fishingStarted  int64   // When fishing started (unix timestamp in milliseconds)
	fishingProgress float64 // Progress from 0.0 to 1.0
	inventorySort   string  // Tracks inventory sort mode: "name", "weight", "value"
	inventoryPage   int     // Current page of inventory when viewing
	itemsPerPage    int     // Number of items to show per page
}

// Custom message type for auto-continuing
type autoContinueMsg time.Time

func autoContinue() tea.Cmd {
	return tea.Tick(time.Second*3, func(t time.Time) tea.Msg {
		return autoContinueMsg(t)
	})
}

func initialModel() model {
	// Initialize UI state
	updateCurrentUIState("menu")

	return model{
		state:           "menu",
		menuItems:       []string{"Go Fishing", "View Inventory", "Quit Game"},
		selectedItem:    0,
		fishingState:    0,
		width:           80,
		height:          24,
		autoFishTick:    0,
		resultTimer:     0,
		fishingDuration: 0,
		fishingStarted:  0,
		fishingProgress: 0.0,
		inventorySort:   "name", // Default sorting by name
		inventoryPage:   0,      // Start at first page
		itemsPerPage:    8,      // Show 8 items per page
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case "menu":
			// Track UI state for background processes
			updateCurrentUIState("menu")
			return m.updateMenu(msg)
		case "fishing":
			// Track UI state for background processes
			updateCurrentUIState("fishing")
			if msg.String() == "q" || msg.String() == "esc" {
				m.state = "menu"
				updateCurrentUIState("menu")
				return m, nil
			} else if msg.String() == "a" { // Toggle auto-fishing while fishing
				autoFishing = !autoFishing
				if autoFishing {
					m.message = "Auto-fishing enabled. Next catches will happen automatically."
				} else {
					m.message = "Auto-fishing disabled."
				}
			} else if msg.String() == "s" { // Save progress
				saveGameProgress()
				m.message = "Game progress saved to /tmp."
			}
			return m, nil
		case "autoFishing":
			// Track UI state for background processes
			updateCurrentUIState("autoFishing")
			if msg.String() == "q" || msg.String() == "esc" {
				autoFishing = false
				m.state = "menu"
				updateCurrentUIState("menu")
				return m, nil
			} else if msg.String() == "s" { // Save progress
				saveGameProgress()
				m.message = "Game progress saved to /tmp."
			}
			// Auto-fishing continues even when user presses other keys
			return m, nil
		case "inventory":
			// Track UI state for background processes
			updateCurrentUIState("inventory")
			if msg.String() == "q" || msg.String() == "esc" || msg.String() == "enter" {
				m.state = "menu"
				updateCurrentUIState("menu")
				return m, nil
			} else if msg.String() == "w" || msg.String() == "2" {
				// Toggle to sort by weight
				m.inventorySort = "weight"
				m.inventoryPage = 0 // Reset to first page when changing sort
				return m, nil
			} else if msg.String() == "v" || msg.String() == "3" {
				// Toggle to sort by value
				m.inventorySort = "value"
				m.inventoryPage = 0 // Reset to first page when changing sort
				return m, nil
			} else if msg.String() == "r" || msg.String() == "1" {
				// Toggle to sort by rarity (renamed from name)
				m.inventorySort = "name" // Keep using "name" as the key for backward compatibility
				m.inventoryPage = 0      // Reset to first page when changing sort
				return m, nil
			} else if msg.String() == "q" || msg.String() == "4" {
				// Toggle to sort by quantity
				m.inventorySort = "quantity"
				m.inventoryPage = 0 // Reset to first page when changing sort
				return m, nil
			} else if msg.String() == "s" { // Save progress
				saveGameProgress()
				m.message = "Game progress saved to /tmp."
			} else if msg.String() == "down" || msg.String() == "j" || msg.String() == "n" {
				// Next page
				m.inventoryPage++
				return m, nil
			} else if msg.String() == "up" || msg.String() == "k" || msg.String() == "p" {
				// Previous page
				if m.inventoryPage > 0 {
					m.inventoryPage--
				}
				return m, nil
			}
		case "fishResult":
			// Track UI state for background processes
			updateCurrentUIState("fishResult")
			// Any key press immediately continues to next step
			if autoFishing {
				// Instead of going directly to auto-fishing state, go to fishing state
				// to show the fishing animation for next catch
				m.state = "fishing"
				m.fishingState = 0
				// Set up a new random fishing duration
				duration := getRandomFishingDuration()
				m.fishingDuration = duration.Milliseconds()
				m.fishingStarted = time.Now().UnixNano() / 1e6
				m.fishingProgress = 0.0
				updateCurrentUIState("fishing")
				return m, tick()
			} else {
				m.state = "menu"
				updateCurrentUIState("menu")
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Update terminal size check for responsive graphics
		updateTerminalSizeCheck(m.width)

	// Custom message types
	case tickMsg:
		if m.state == "fishing" {
			// Calculate fishing progress
			now := time.Now().UnixNano() / 1e6
			elapsed := now - m.fishingStarted
			m.fishingProgress = float64(elapsed) / float64(m.fishingDuration)
			m.fishingState++ // Increment for animation frames

			// Check if fishing is complete
			if m.fishingProgress >= 1.0 {
				return m.completeFishing()
			}

			// Continue animation
			return m, tick()
		}
	case autoTickMsg:
		if m.state == "autoFishing" {
			// When the auto-fishing timer fires, transition to fishing state
			// to show the fishing animation
			m.state = "fishing"
			m.fishingState = 0
			// Set up a new random fishing duration
			duration := getRandomFishingDuration()
			m.fishingDuration = duration.Milliseconds()
			m.fishingStarted = time.Now().UnixNano() / 1e6
			m.fishingProgress = 0.0
			updateCurrentUIState("fishing")
			return m, tick()
		}
	case autoContinueMsg:
		// Auto-continue after showing the result for a moment
		if m.state == "fishResult" && autoFishing {
			// Show fishing animation after catch result
			m.state = "fishing"
			m.fishingState = 0
			// Set up a new random fishing duration
			duration := getRandomFishingDuration()
			m.fishingDuration = duration.Milliseconds()
			m.fishingStarted = time.Now().UnixNano() / 1e6
			m.fishingProgress = 0.0
			updateCurrentUIState("fishing")
			return m, tick()
		}
		return m, nil
	case catchResultMsg:
		if msg.success {
			m.catchSuccess = true
			m.caughtFish = msg.fish
			m.fishShape = fishShapes[rand.Intn(len(fishShapes))]
		} else {
			m.catchSuccess = false
		}
		m.state = "fishResult"
		updateCurrentUIState("fishResult")
		m.resultTimer = 0

		// If auto-fishing is enabled, automatically continue after showing results
		if autoFishing {
			return m, autoContinue()
		}
		return m, nil
	}

	return m, nil
}

// Handle completion of fishing and determine catch
func (m model) completeFishing() (tea.Model, tea.Cmd) {
	// Calculate catch chance with weather and time of day factors
	catchChance := (rand.Float64() * 10) + float64(player.RodStrength) + float64(player.BaitStrength)
	catchChance *= weatherFactor
	catchChance *= timeFactor // Apply time of day factor

	success := catchChance >= 5

	var fish game.Fish
	if success {
		// Choose a fish based on rarity and time of day
		fish = chooseFish()
		// Add to inventory
		player.AddFish(fish)
		// Auto-save when a fish is caught
		saveGameProgress()
	}

	// Return catch result
	return m, func() tea.Msg {
		return catchResultMsg{
			success: success,
			fish:    fish,
		}
	}
}

func (m model) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		close(stopIdle) // Safely close channel
		return m, tea.Quit
	case "up", "k":
		if m.selectedItem > 0 {
			m.selectedItem--
		}
	case "down", "j":
		if m.selectedItem < len(m.menuItems)-1 {
			m.selectedItem++
		}
	case "enter", " ":
		switch m.selectedItem {
		case 0: // Go Fishing
			m.state = "fishing"
			m.fishingState = 0
			m.message = ""
			// Set up a new random fishing duration
			duration := getRandomFishingDuration()
			m.fishingDuration = duration.Milliseconds()
			m.fishingStarted = time.Now().UnixNano() / 1e6
			m.fishingProgress = 0.0
			return m, tick()
		case 1: // View Inventory
			m.state = "inventory"
			m.message = ""
		case 2: // Quit
			close(stopIdle) // Safely close channel
			return m, tea.Quit
		}
	case "a": // Toggle auto-fishing with 'a' key from anywhere in the menu
		autoFishing = !autoFishing
		if autoFishing {
			m.message = "Auto-fishing enabled. Press 'a' again to disable."
		} else {
			m.message = "Auto-fishing disabled."
		}
	}
	return m, nil
}

// Style definitions using lipgloss
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(1, 0, 1, 0)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#61AFEF")).
			Padding(1, 2)

	accentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#61AFEF")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E06C75")).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#98C379")).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#56B6C2"))

	menuItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1)

	highlightedMenuItemStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color("#000000")).
					Background(lipgloss.Color("#61AFEF")).
					Padding(0, 1)

	fishStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#56B6C2")).
			Bold(true)
)

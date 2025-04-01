package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/user/fishing-game/game"
)

// Global variables
var (
	player         game.Player
	availableFish  []game.Fish
	weatherFactor  float64 = 1.0 // Default weather factor
	idleCatchRate  float64 = 0.3 // Fish caught per minute while idle
	lastActiveTime time.Time
	mu             sync.Mutex // Mutex to prevent race conditions
	stopIdle       chan bool  // Channel to stop idle routine
	autoFishing    bool       // Auto fishing enabled
	testMode       bool       // Test mode for faster fishing
	saveFile       string     // Path to save file

	// Time of day variables
	timeOfDay  string        // Current time period (morning, afternoon, evening, night)
	timeFactor float64 = 1.0 // How time of day affects fishing success
)

// TimeOfDay represents different fishing periods with unique characteristics
type TimeOfDay struct {
	Name        string  // Name of the period
	StartHour   int     // Hour when this period starts (24h format)
	EndHour     int     // Hour when this period ends (24h format)
	CatchFactor float64 // Multiplier for catch rates
	Icon        string  // Visual representation
	Description string  // Short description of fishing conditions
}

// Define time periods
var timePeriods = []TimeOfDay{
	{
		Name:        "Morning",
		StartHour:   5,
		EndHour:     11,
		CatchFactor: 1.2,
		Icon:        "🌅",
		Description: "Perfect for early biters",
	},
	{
		Name:        "Afternoon",
		StartHour:   11,
		EndHour:     17,
		CatchFactor: 0.8,
		Icon:        "☀️",
		Description: "Slower fishing during midday heat",
	},
	{
		Name:        "Evening",
		StartHour:   17,
		EndHour:     21,
		CatchFactor: 1.3,
		Icon:        "🌇",
		Description: "Prime fishing hours, increased activity",
	},
	{
		Name:        "Night",
		StartHour:   21,
		EndHour:     5,
		CatchFactor: 0.9,
		Icon:        "🌙",
		Description: "Good for nocturnal species",
	},
}

// GameSave represents the saveable game state
type GameSave struct {
	Player         game.Player
	WeatherFactor  float64
	LastActiveTime time.Time
	AutoFishing    bool
	SaveTime       time.Time // When the game was saved
}

func main() {
	// Parse command line flags
	flag.BoolVar(&testMode, "test", false, "Run in test mode with shorter fishing times (5-10 seconds)")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	// Set save file path in /tmp directory
	saveFile = filepath.Join(os.TempDir(), "fishing-game-save.json")

	// Initialize game
	initializeGame()

	// Start background routines
	startBackgroundRoutines()

	// Start auto-save routine
	go startAutoSaveRoutine()

	// Start the Bubble Tea program
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}

	// Save progress on exit
	saveGameProgress()
}

// Initialize game state and variables
func initializeGame() {
	// Try to load saved game first
	if !loadGameProgress() {
		// If loading fails, start a new game
		player = game.NewPlayer()
		lastActiveTime = time.Now()
		autoFishing = false
		weatherFactor = 1.0
	}

	// These are always initialized fresh
	availableFish = game.GetAllFish()
	stopIdle = make(chan bool)

	// Initialize time of day
	updateTimeOfDay()
}

// updateTimeOfDay checks the current system time and updates time-related variables
func updateTimeOfDay() {
	currentHour := time.Now().Hour()

	// Find which time period we're in
	for _, period := range timePeriods {
		// Handle periods that cross midnight
		if period.StartHour > period.EndHour {
			// Night period (e.g., 21:00 to 05:00)
			if currentHour >= period.StartHour || currentHour < period.EndHour {
				timeOfDay = period.Name
				timeFactor = period.CatchFactor
				return
			}
		} else {
			// Normal period within same day
			if currentHour >= period.StartHour && currentHour < period.EndHour {
				timeOfDay = period.Name
				timeFactor = period.CatchFactor
				return
			}
		}
	}

	// Default fallback
	timeOfDay = "Afternoon"
	timeFactor = 1.0
}

// saveGameProgress saves the game state to the save file
func saveGameProgress() {
	mu.Lock()
	defer mu.Unlock()

	gameSave := GameSave{
		Player:         player,
		WeatherFactor:  weatherFactor,
		LastActiveTime: lastActiveTime,
		AutoFishing:    autoFishing,
		SaveTime:       time.Now(),
	}

	data, err := json.Marshal(gameSave)
	if err != nil {
		fmt.Printf("Error marshalling save data: %v\n", err)
		return
	}

	err = ioutil.WriteFile(saveFile, data, 0644)
	if err != nil {
		fmt.Printf("Error writing save file: %v\n", err)
	}
}

// loadGameProgress loads the game state from the save file
func loadGameProgress() bool {
	data, err := ioutil.ReadFile(saveFile)
	if err != nil {
		// Save file doesn't exist or can't be read
		return false
	}

	var gameSave GameSave
	err = json.Unmarshal(data, &gameSave)
	if err != nil {
		fmt.Printf("Error unmarshalling save data: %v\n", err)
		return false
	}

	// Restore game state
	player = gameSave.Player
	weatherFactor = gameSave.WeatherFactor
	lastActiveTime = gameSave.LastActiveTime
	autoFishing = gameSave.AutoFishing

	// Print load message with timestamp
	saveTimeStr := gameSave.SaveTime.Format("Jan 2 15:04:05")
	fmt.Printf("Loaded game save from %s (%d fish caught)\n", saveTimeStr, len(player.FishCaught))

	return true
}

// startAutoSaveRoutine periodically saves the game progress
func startAutoSaveRoutine() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			saveGameProgress()
		case <-stopIdle:
			return
		}
	}
}

// Generate a random fishing duration
func getRandomFishingDuration() time.Duration {
	if testMode {
		// 5-10 seconds in test mode
		seconds := rand.Intn(6) + 5 // 5-10 range
		return time.Second * time.Duration(seconds)
	} else {
		// 5 seconds to 5 minutes in normal mode
		seconds := rand.Intn(295) + 5 // 5-300 seconds (5 min max)
		return time.Second * time.Duration(seconds)
	}
}

// Start all background routines with panic recovery
func startBackgroundRoutines() {
	// Start idle routine
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered from panic in idle routine:", r)
			}
		}()
		idleRoutine()
	}()

	// Start time of day routine
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered from panic in time routine:", r)
			}
		}()

		// Check time every minute
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		// Initial update (already done in initializeGame)

		for {
			select {
			case <-ticker.C:
				updateTimeOfDay()
			case <-stopIdle:
				return
			}
		}
	}()

	// Start weather routine
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered from panic in weather routine:", r)
			}
		}()
		// Update weather every 15 minutes
		ticker := time.NewTicker(15 * time.Minute)
		defer ticker.Stop()

		// Initial update
		updateWeatherFactor()

		for {
			select {
			case <-ticker.C:
				updateWeatherFactor()
			case <-stopIdle:
				return
			}
		}
	}()

	// Start auto-fishing routine
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered from panic in auto-fishing routine:", r)
			}
		}()
		autoFishingRoutine()
	}()
}

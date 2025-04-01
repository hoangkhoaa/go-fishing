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
	saveDir        string     // Directory for save files
	todaySaveFile  string     // Path to today's save file

	// History tracking
	dailyCatches     map[string][]game.Fish // Map of date strings to fish catches
	dateList         []string               // List of dates with catches, sorted
	viewingDate      string                 // Currently viewed date in history
	isViewingHistory bool                   // Whether user is viewing history

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

// DailySave represents the saveable game state for a single day
type DailySave struct {
	FishCaught []game.Fish // Fish caught on this day
	Date       string      // Date in YYYY-MM-DD format
	SaveTime   time.Time   // When the game was last saved
}

// GameSave represents the main saveable game state (excluding daily catches)
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

	// Set up save directory structure
	setupSaveDirectory()

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

// setupSaveDirectory creates the save directory structure
func setupSaveDirectory() {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		// Fallback to temp directory if needed
		cwd = os.TempDir()
	}

	// Create main saves directory
	saveDir = filepath.Join(cwd, "saves")
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		fmt.Printf("Error creating save directory: %v\n", err)
	}

	// Set today's save file path
	today := time.Now().Format("2006-01-02")
	todaySaveFile = filepath.Join(saveDir, today+".json")
	viewingDate = today
}

// Initialize game state and variables
func initializeGame() {
	// Initialize maps and slices
	dailyCatches = make(map[string][]game.Fish)
	dateList = []string{}
	isViewingHistory = false

	// Try to load saved game and history
	if !loadGameProgress() {
		// If loading fails, start a new game
		player = game.NewPlayer()
		lastActiveTime = time.Now()
		autoFishing = false
		weatherFactor = 1.0
	}

	// Load today's catches if they exist
	loadTodayCatches()

	// Load all available dates and their catches
	loadAllDailyCatches()

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

// saveGameProgress saves the main game state (without fish data)
func saveGameProgress() {
	mu.Lock()
	defer mu.Unlock()

	// Save main game state (without fish catches)
	gameSave := GameSave{
		Player:         player,
		WeatherFactor:  weatherFactor,
		LastActiveTime: lastActiveTime,
		AutoFishing:    autoFishing,
		SaveTime:       time.Now(),
	}

	// Backup player's full fish list
	allFish := player.FishCaught

	// To avoid duplicating data, clear the fish list as it's saved by day
	player.FishCaught = []game.Fish{}

	data, err := json.Marshal(gameSave)
	if err != nil {
		fmt.Printf("Error marshalling save data: %v\n", err)
		return
	}

	// Save main game state
	mainSaveFile := filepath.Join(saveDir, "game_state.json")
	err = ioutil.WriteFile(mainSaveFile, data, 0644)
	if err != nil {
		fmt.Printf("Error writing main save file: %v\n", err)
	}

	// Restore player's fish list for runtime use
	player.FishCaught = allFish

	// Save today's catches separately
	saveTodayCatches()
}

// saveTodayCatches saves only the fish caught today
func saveTodayCatches() {
	today := time.Now().Format("2006-01-02")

	// Create daily save object
	dailySave := DailySave{
		FishCaught: player.FishCaught,
		Date:       today,
		SaveTime:   time.Now(),
	}

	data, err := json.Marshal(dailySave)
	if err != nil {
		fmt.Printf("Error marshalling daily save data: %v\n", err)
		return
	}

	// Save today's catches
	err = ioutil.WriteFile(todaySaveFile, data, 0644)
	if err != nil {
		fmt.Printf("Error writing daily save file: %v\n", err)
	}

	// Update in-memory cache
	dailyCatches[today] = player.FishCaught

	// Update date list if needed
	if !contains(dateList, today) {
		dateList = append(dateList, today)
	}
}

// loadGameProgress loads the main game state
func loadGameProgress() bool {
	mainSaveFile := filepath.Join(saveDir, "game_state.json")

	data, err := ioutil.ReadFile(mainSaveFile)
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

	// Restore game state (without fish catches)
	player = gameSave.Player
	player.FishCaught = []game.Fish{} // Clear any fish data that might be in the main save
	weatherFactor = gameSave.WeatherFactor
	lastActiveTime = gameSave.LastActiveTime
	autoFishing = gameSave.AutoFishing

	// Print load message with timestamp
	saveTimeStr := gameSave.SaveTime.Format("Jan 2 15:04:05")
	fmt.Printf("Loaded game state from %s\n", saveTimeStr)

	return true
}

// loadTodayCatches loads the fish caught today
func loadTodayCatches() {
	data, err := ioutil.ReadFile(todaySaveFile)
	if err != nil {
		// No catches today yet or file can't be read
		return
	}

	var dailySave DailySave
	err = json.Unmarshal(data, &dailySave)
	if err != nil {
		fmt.Printf("Error unmarshalling daily save data: %v\n", err)
		return
	}

	// Set the fish caught for today
	player.FishCaught = dailySave.FishCaught

	// Calculate total weight and value
	player.TotalWeight = 0
	player.TotalValue = 0
	for _, fish := range player.FishCaught {
		player.TotalWeight += fish.Weight
		player.TotalValue += fish.Value
	}

	// Update in-memory cache
	today := time.Now().Format("2006-01-02")
	dailyCatches[today] = player.FishCaught

	// Print load message
	fmt.Printf("Loaded %d fish caught today\n", len(player.FishCaught))
}

// loadAllDailyCatches scans the save directory and loads all daily catches
func loadAllDailyCatches() {
	// Scan save directory for daily save files
	files, err := ioutil.ReadDir(saveDir)
	if err != nil {
		fmt.Printf("Error scanning save directory: %v\n", err)
		return
	}

	dateList = []string{}

	// Load each daily save file
	for _, file := range files {
		if file.IsDir() || file.Name() == "game_state.json" {
			continue
		}

		// Extract date from filename (remove .json extension)
		date := file.Name()
		if len(date) > 5 && date[len(date)-5:] == ".json" {
			date = date[:len(date)-5]
		}

		// Check if date is in valid format (YYYY-MM-DD)
		if _, err := time.Parse("2006-01-02", date); err != nil {
			continue
		}

		// Add to date list
		dateList = append(dateList, date)

		// Load fish catches for this date
		filePath := filepath.Join(saveDir, file.Name())
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			continue
		}

		var dailySave DailySave
		err = json.Unmarshal(data, &dailySave)
		if err != nil {
			continue
		}

		// Store in the map
		dailyCatches[date] = dailySave.FishCaught
	}
}

// contains checks if a string is in a slice
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// setViewingDate changes the current date being viewed
func setViewingDate(date string) {
	viewingDate = date
	isViewingHistory = date != time.Now().Format("2006-01-02")
}

// getFishCaughtOnDate returns the fish caught on a specific date
func getFishCaughtOnDate(date string) []game.Fish {
	// If viewing today, use the player's current catches
	today := time.Now().Format("2006-01-02")
	if date == today {
		return player.FishCaught
	}

	// Otherwise, get fish from the daily catches map
	if catches, ok := dailyCatches[date]; ok {
		return catches
	}

	// If no catches found for this date
	return []game.Fish{}
}

// getAvailableDates returns the dates that have fish catches
func getAvailableDates() []string {
	return dateList
}

// getFishCaughtDetails returns statistics for fish caught on a date
func getFishCaughtDetails(date string) (int, int, int) {
	catches := getFishCaughtOnDate(date)

	totalWeight := 0
	totalValue := 0

	for _, fish := range catches {
		totalWeight += fish.Weight
		totalValue += fish.Value
	}

	return len(catches), totalWeight, totalValue
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

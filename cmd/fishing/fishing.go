package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/user/fishing-game/game"
)

// Custom message type for fishing animation
type tickMsg time.Time

func tick() tea.Cmd {
	// Use 250ms ticks for smoother animation progress
	return tea.Tick(time.Millisecond*250, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Custom message type for auto-fishing
type autoTickMsg time.Time

func autoTick() tea.Cmd {
	// Change auto-fishing to happen every 10 seconds
	return tea.Tick(time.Second*10, func(t time.Time) tea.Msg {
		return autoTickMsg(t)
	})
}

// Custom message type for catch result
type catchResultMsg struct {
	success bool
	fish    game.Fish
}

// The main fishing animation function is now in model.go as part of the Update method
// and the catch completion logic is in the completeFishing method

func (m model) updateAutoFishingAnimation() (tea.Model, tea.Cmd) {
	// Progress through auto-fishing animation
	m.autoFishTick++

	// Every 10 ticks, try to catch a fish
	if m.autoFishTick >= 10 {
		// Calculate catch chance with weather factor and time of day factor
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

			// Reset tick counter and return catch result
			m.autoFishTick = 0
			return m, func() tea.Msg {
				return catchResultMsg{
					success: success,
					fish:    fish,
				}
			}
		}

		// If not successful, reset counter but keep auto-fishing
		m.autoFishTick = 0
		m.autoFishMsg = "No bites yet, still fishing..."
	} else {
		// Update auto-fishing message with animation dots
		dots := strings.Repeat(".", (m.autoFishTick%3)+1)
		m.autoFishMsg = fmt.Sprintf("Auto-fishing in progress%s", dots)
	}

	// Continue auto-fishing
	return m, autoTick()
}

func chooseFish() game.Fish {
	mu.Lock()
	defer mu.Unlock()

	// Decide whether to catch trash (10-15% chance)
	trashChance := rand.Float64()
	if trashChance < 0.12 {
		trashItems := game.GetTrashItems()
		if len(trashItems) > 0 {
			return trashItems[rand.Intn(len(trashItems))]
		}
	}

	// Small chance for legendary fish (0.5-2% depending on time/weather)
	legendaryChance := rand.Float64()
	legendaryThreshold := 0.005 // Base 0.5% chance

	// Better weather increases legendary chance
	if weatherFactor > 1.2 {
		legendaryThreshold += 0.005
	}

	// Night time is best for most legendary catches
	if timeOfDay == "Night" {
		legendaryThreshold += 0.01
	}

	if legendaryChance < legendaryThreshold {
		legendaryFish := game.GetLegendaryFish()
		// Filter for ones that prefer current time
		timeSpecificLegendary := []game.Fish{}
		for _, fish := range legendaryFish {
			if fish.PreferredTime == timeOfDay || fish.PreferredTime == "" {
				timeSpecificLegendary = append(timeSpecificLegendary, fish)
			}
		}

		if len(timeSpecificLegendary) > 0 {
			return timeSpecificLegendary[rand.Intn(len(timeSpecificLegendary))]
		} else if len(legendaryFish) > 0 {
			return legendaryFish[rand.Intn(len(legendaryFish))]
		}
	}

	// Get fish that prefer current time of day or have no specific time preference
	timeFish := game.GetFishByTimeOfDay(timeOfDay)
	if len(timeFish) == 0 {
		// Fallback to all fish if no time-appropriate fish
		timeFish = availableFish
	}

	// Calculate total rarity, adjusted by weather and time factors
	totalRarity := 0
	adjustedRarities := make([]int, len(timeFish))

	for i, fish := range timeFish {
		// Adjust rarity based on weather
		adjustedRarity := fish.Rarity

		// Make rare fish more common in good weather
		if weatherFactor > 1.2 && fish.Rarity <= 2 {
			adjustedRarity += 1
		}

		// Make common fish less common in bad weather
		if weatherFactor < 0.8 && fish.Rarity >= 8 {
			adjustedRarity -= 1
		}

		// Time of day specific adjustments
		if fish.PreferredTime == timeOfDay {
			// Strongly boost fish during their preferred time (+3 to rarity)
			adjustedRarity += 3
		} else {
			// Additional habitat-based adjustments as fallback
			switch timeOfDay {
			case "Morning":
				// Surface feeders are more common in morning
				if strings.Contains(strings.ToLower(fish.Habitat), "surface") {
					adjustedRarity += 1
				}
			case "Afternoon":
				// Deep water fish more common in afternoon
				if strings.Contains(strings.ToLower(fish.Habitat), "deep") {
					adjustedRarity += 1
				}
			case "Evening":
				// Predatory fish more active in evening
				if fish.Weight > 20 {
					adjustedRarity += 1
				}
			case "Night":
				// Nocturnal fish more common at night
				if strings.Contains(strings.ToLower(fish.Color), "black") ||
					strings.Contains(strings.ToLower(fish.Color), "dark") {
					adjustedRarity += 1
				}
			}
		}

		if adjustedRarity < 1 {
			adjustedRarity = 1
		}

		adjustedRarities[i] = adjustedRarity
		totalRarity += adjustedRarity
	}

	// Choose a random number between 0 and totalRarity
	randomNum := rand.Intn(totalRarity)

	// Pick a fish based on the random number and adjusted fish rarity
	currentSum := 0
	for i, fish := range timeFish {
		currentSum += adjustedRarities[i]
		if randomNum < currentSum {
			return fish
		}
	}

	// Default return (should never happen)
	return timeFish[0]
}

package main

import (
	"math/rand"
	"time"

	"github.com/user/fishing-game/game"
)

// Background processes for idle catching and weather updates

func idleRoutine() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			processCatchesWhileAway()
		case <-stopIdle:
			return
		}
	}
}

func autoFishingRoutine() {
	// Auto-fishing timer (initial value)
	duration := getRandomFishingDuration()
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if autoFishing {
				// Only catch fish in the background if we're not currently showing fishing in the UI
				// Check the current state without locking since this is just a rough check
				if currentUIState != "fishing" && currentUIState != "fishResult" {
					mu.Lock()
					// Calculate catch chance with weather factor
					catchChance := (rand.Float64() * 10) + float64(player.RodStrength) + float64(player.BaitStrength)
					catchChance *= weatherFactor

					if catchChance >= 5 {
						// Choose a fish based on rarity
						totalRarity := 0
						for _, fish := range availableFish {
							totalRarity += fish.Rarity
						}

						randomNum := rand.Intn(totalRarity)

						currentSum := 0
						var chosenFish game.Fish
						for _, fish := range availableFish {
							currentSum += fish.Rarity
							if randomNum < currentSum {
								chosenFish = fish
								break
							}
						}

						if (chosenFish != game.Fish{}) {
							player.AddFish(chosenFish)
							// Auto-save when a fish is caught in background
							saveGameProgress()
						}
					}
					mu.Unlock()
				}

				// Set a new random fishing duration
				newDuration := getRandomFishingDuration()
				ticker.Reset(newDuration)
			}
		case <-stopIdle:
			return
		}
	}
}

func processCatchesWhileAway() {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()
	minutesAway := now.Sub(lastActiveTime).Minutes()

	if minutesAway < 1 {
		return
	}

	// Calculate how many fish were caught while away
	catchChance := idleCatchRate * minutesAway * weatherFactor
	wholeCatches := int(catchChance)

	// Chance for an additional catch
	fractionalCatch := catchChance - float64(wholeCatches)
	if rand.Float64() < fractionalCatch {
		wholeCatches++
	}

	// Process the catches - simplified to avoid potential issues
	for i := 0; i < wholeCatches; i++ {
		// Choose a random fish directly to avoid complexity
		if len(availableFish) > 0 {
			randomIndex := rand.Intn(len(availableFish))
			player.AddFish(availableFish[randomIndex])
		}
	}

	// Auto-save after processing idle catches
	if wholeCatches > 0 {
		saveGameProgress()
	}

	lastActiveTime = now
}

func updateWeatherFactor() {
	// Simplified weather system to reduce complexity
	mu.Lock()
	defer mu.Unlock()

	// Simply vary the weather factor between 0.7 and 1.3
	weatherFactor = 0.7 + rand.Float64()*0.6
}

// Track current UI state for coordination with background processes
var currentUIState string

// Update the current UI state
func updateCurrentUIState(state string) {
	currentUIState = state
}

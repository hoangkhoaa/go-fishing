package game

// Player represents the player's stats and inventory
type Player struct {
	Money        int
	FishCaught   []Fish
	TotalWeight  int
	TotalValue   int
	FishingRod   string
	RodStrength  int
	Bait         string
	BaitStrength int
}

// NewPlayer creates a new player with default values
func NewPlayer() Player {
	return Player{
		Money:        50,
		FishCaught:   []Fish{},
		TotalWeight:  0,
		TotalValue:   0,
		FishingRod:   "Basic Rod",
		RodStrength:  1,
		Bait:         "Worm",
		BaitStrength: 1,
	}
}

// AddFish adds a fish to the player's inventory
func (p *Player) AddFish(fish Fish) {
	p.FishCaught = append(p.FishCaught, fish)
	p.TotalWeight += fish.Weight
	p.TotalValue += fish.Value
}

// SellAllFish sells all fish in the player's inventory
func (p *Player) SellAllFish() int {
	if len(p.FishCaught) == 0 {
		return 0
	}

	totalValue := 0
	for _, fish := range p.FishCaught {
		totalValue += fish.Value
	}

	p.Money += totalValue

	// Reset fish inventory
	p.FishCaught = []Fish{}
	p.TotalWeight = 0
	p.TotalValue = 0

	return totalValue
}

// BuyRod allows the player to buy a fishing rod
func (p *Player) BuyRod(rodName string, cost int, strength int) bool {
	if p.Money < cost {
		return false
	}

	p.Money -= cost
	p.FishingRod = rodName
	p.RodStrength = strength
	return true
}

// BuyBait allows the player to buy bait
func (p *Player) BuyBait(baitName string, cost int, strength int) bool {
	if p.Money < cost {
		return false
	}

	p.Money -= cost
	p.Bait = baitName
	p.BaitStrength = strength
	return true
}

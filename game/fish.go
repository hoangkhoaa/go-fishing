package game

// Fish represents a fish that can be caught
type Fish struct {
	Name          string
	Weight        int
	Rarity        int // 1 = extremely rare, 10 = very common
	Value         int
	CatchMsg      string
	Color         string // Primary color of the fish
	Pattern       string // Pattern type (spotted, striped, plain, etc)
	Habitat       string // Where the fish is typically found
	PreferredTime string // Time of day when this fish is most active: Morning, Afternoon, Evening, Night, or "" for no preference
	IsTrash       bool   // Whether this is a trash item rather than a fish
	IsLegendary   bool   // Whether this is a legendary/mythical creature
}

// GetAllFish returns a slice of all available fish in the game
func GetAllFish() []Fish {
	regularFish := []Fish{
		// Common Fish - Rarity 8-10
		{"Minnow", 1, 10, 2, "You caught a tiny Minnow!", "Silver", "Plain", "Freshwater", "Morning", false, false},
		{"Goldfish", 1, 10, 3, "You caught a Goldfish!", "Gold", "Plain", "Pond", "Afternoon", false, false},
		{"Carp", 4, 9, 5, "You caught a Carp!", "Brown", "Mottled", "Freshwater", "Afternoon", false, false},
		{"Perch", 3, 9, 6, "You caught a Perch!", "Yellow", "Striped", "Lake", "Evening", false, false},
		{"Bluegill", 2, 9, 4, "You caught a Bluegill!", "Blue", "Spotted", "Freshwater", "Morning", false, false},
		{"Trout", 3, 8, 7, "You caught a Trout!", "Rainbow", "Spotted", "Stream", "Morning", false, false},
		{"Sunfish", 2, 8, 5, "You caught a Sunfish!", "Orange", "Spotted", "Pond", "Afternoon", false, false},
		{"Crappie", 2, 8, 5, "You caught a Crappie!", "Silver", "Mottled", "Lake", "Evening", false, false},
		{"Bullhead", 4, 8, 6, "You caught a Bullhead!", "Black", "Plain", "Lake", "Night", false, false},
		{"Bream", 3, 8, 5, "You caught a Bream!", "Bronze", "Plain", "Freshwater", "", false, false},

		// Moderately Common Fish - Rarity 6-7
		{"Bass", 5, 7, 10, "You caught a Bass!", "Green", "Spotted", "Lake", "Evening", false, false},
		{"Catfish", 8, 7, 12, "You caught a Catfish!", "Gray", "Mottled", "River", "Night", false, false},
		{"Pike", 7, 7, 11, "You caught a Pike!", "Green", "Striped", "Lake", "Evening", false, false},
		{"Walleye", 6, 7, 10, "You caught a Walleye!", "Yellow", "Mottled", "Lake", "Night", false, false},
		{"Rainbow Trout", 4, 7, 9, "You caught a Rainbow Trout!", "Rainbow", "Spotted", "Stream", "Morning", false, false},
		{"Salmon", 8, 6, 15, "You caught a Salmon!", "Pink", "Plain", "River", "Morning", false, false},
		{"Tilapia", 5, 6, 8, "You caught a Tilapia!", "Silver", "Plain", "Lake", "", false, false},
		{"Yellowtail", 7, 6, 12, "You caught a Yellowtail!", "Yellow", "Striped", "Ocean", "Afternoon", false, false},
		{"Rock Bass", 4, 6, 8, "You caught a Rock Bass!", "Brown", "Spotted", "Lake", "", false, false},
		{"Channel Catfish", 9, 6, 14, "You caught a Channel Catfish!", "Gray", "Plain", "River", "Night", false, false},

		// Uncommon Fish - Rarity 4-5
		{"Halibut", 15, 5, 25, "You caught a Halibut!", "Brown", "Mottled", "Ocean Floor", "Afternoon", false, false},
		{"Sea Bass", 12, 5, 20, "You caught a Sea Bass!", "Black", "Plain", "Ocean", "Evening", false, false},
		{"Snapper", 10, 5, 18, "You caught a Snapper!", "Red", "Plain", "Reef", "Afternoon", false, false},
		{"Flounder", 8, 5, 16, "You caught a Flounder!", "Sand", "Spotted", "Ocean Floor", "Night", false, false},
		{"Grouper", 14, 5, 22, "You caught a Grouper!", "Brown", "Mottled", "Reef", "Evening", false, false},
		{"Cod", 11, 5, 19, "You caught a Cod!", "Gray", "Spotted", "Deep Sea", "Morning", false, false},
		{"Mahi-Mahi", 15, 4, 28, "You caught a beautiful Mahi-Mahi!", "Blue-Green", "Spotted", "Open Ocean", "Afternoon", false, false},
		{"Snook", 13, 4, 24, "You caught a Snook!", "Silver", "Black Stripe", "Coastal", "Night", false, false},
		{"Amberjack", 16, 4, 26, "You caught an Amberjack!", "Silver", "Yellow", "Reef", "Morning", false, false},
		{"Lake Trout", 12, 4, 22, "You caught a Lake Trout!", "Silver", "Spotted", "Deep Lake", "Morning", false, false},

		// Rare Fish - Rarity 2-3
		{"Tuna", 30, 3, 45, "You caught a massive Tuna!", "Blue", "Silver Belly", "Open Ocean", "Afternoon", false, false},
		{"Tarpon", 40, 3, 50, "You caught a mighty Tarpon!", "Silver", "Iridescent", "Coastal", "Evening", false, false},
		{"Barracuda", 25, 3, 40, "You caught a toothy Barracuda!", "Silver", "Striped", "Reef", "Evening", false, false},
		{"Cobia", 35, 3, 48, "You caught a powerful Cobia!", "Brown", "White Stripe", "Coastal", "Afternoon", false, false},
		{"Sturgeon", 45, 3, 55, "You caught an ancient Sturgeon!", "Gray", "Armored", "River", "Night", false, false},
		{"Striped Bass", 22, 3, 38, "You caught a huge Striped Bass!", "Silver", "Black Stripes", "Coastal", "Morning", false, false},
		{"Redfish", 20, 2, 35, "You caught a prized Redfish!", "Red", "Spotted Tail", "Coastal", "Evening", false, false},
		{"King Mackerel", 28, 2, 42, "You caught a King Mackerel!", "Silver", "Spotted", "Open Ocean", "Morning", false, false},
		{"Bonefish", 18, 2, 32, "You caught a Bonefish!", "Silver", "Dark Back", "Flats", "Morning", false, false},
		{"Permit", 25, 2, 40, "You caught a Permit!", "Silver", "Yellow Fins", "Flats", "Afternoon", false, false},

		// Very Rare Fish - Rarity 1
		{"Marlin", 180, 1, 200, "You caught a massive Marlin!", "Blue", "Striped", "Deep Ocean", "Afternoon", false, false},
		{"Swordfish", 150, 1, 180, "You caught a magnificent Swordfish!", "Blue-Black", "Plain", "Deep Ocean", "Night", false, false},
		{"Sailfish", 130, 1, 175, "You caught a beautiful Sailfish!", "Blue", "Spotted Sail", "Tropical Ocean", "Morning", false, false},
		{"Giant Trevally", 100, 1, 150, "You caught a Giant Trevally!", "Silver", "Dark Back", "Reef", "Evening", false, false},
		{"Goliath Grouper", 300, 1, 250, "You caught a massive Goliath Grouper!", "Brown", "Mottled", "Reef", "Afternoon", false, false},
		{"Arapaima", 180, 1, 190, "You caught a prehistoric Arapaima!", "Red", "Scaled", "Amazon", "Evening", false, false},
		{"Giant Squid", 400, 1, 300, "You caught a rare Giant Squid!", "Red", "Tentacled", "Deep Ocean", "Night", false, false},
		{"Mekong Giant Catfish", 280, 1, 280, "You caught a Mekong Giant Catfish!", "Gray", "Plain", "Mekong River", "Night", false, false},
		{"Bluefin Tuna", 500, 1, 400, "You caught a prized Bluefin Tuna!", "Blue", "Silver Belly", "Open Ocean", "Morning", false, false},
		{"Golden Dorado", 80, 1, 120, "You caught a spectacular Golden Dorado!", "Gold", "Patterned", "South American Rivers", "Afternoon", false, false},
	}

	// Legendary/Mythical Creatures - Even rarer than rarity 1
	legendaryFish := []Fish{
		// Legendary creatures (extremely rare, valuable, and time-specific)
		{"Kraken", 800, 1, 1000, "You caught the mythical KRAKEN! Its tentacles nearly capsize your boat!", "Dark Purple", "Tentacled", "Abyss", "Night", false, true},
		{"Loch Ness Monster", 1200, 1, 1500, "You've captured proof of Nessie! The scientific community is in shock!", "Green", "Prehistoric", "Deep Lake", "Night", false, true},
		{"Megalodon", 2000, 1, 2000, "MEGALODON! You've caught a living prehistoric shark thought extinct for millions of years!", "Gray", "Ancient", "Deep Ocean", "Night", false, true},
		{"Mermaid", 120, 1, 5000, "A MERMAID has been caught in your net! She grants you a wish before returning to the sea.", "Iridescent", "Scaled", "Tropical Ocean", "Evening", false, true},
		{"Golden Carp", 50, 1, 800, "The legendary GOLDEN CARP! Legend says it brings wealth and prosperity!", "Gold", "Glowing", "Sacred Lake", "Morning", false, true},
		{"Phoenix Fish", 30, 1, 1200, "A PHOENIX FISH! Its scales glow like embers and it's warm to the touch!", "Fiery Red", "Glowing", "Volcanic Vent", "Afternoon", false, true},
		{"Ghost Whale", 1500, 1, 1800, "A GHOST WHALE has appeared! Its translucent body glows with an otherworldly light.", "Pale Blue", "Translucent", "Phantom Depths", "Night", false, true},
		{"Dragon Eel", 200, 1, 1600, "A DRAGON EEL! It breathes small flames and has scales harder than steel!", "Crimson", "Armored", "Undersea Cave", "Evening", false, true},
		{"Abyssal Anglerfish", 80, 1, 1300, "An ABYSSAL ANGLERFISH! Its light mesmerizes you with hypnotic patterns!", "Black", "Bioluminescent", "Hadal Zone", "Night", false, true},
		{"Moonlight Jellyfish", 40, 1, 900, "A MOONLIGHT JELLYFISH! It seems to channel the very essence of moonlight!", "Silver", "Glowing", "Midnight Surface", "Night", false, true},
	}

	// Trash items (common, worthless, and a nuisance)
	trashItems := []Fish{
		{"Old Boot", 2, 9, 0, "You caught an old boot. What a disappointment!", "Brown", "Worn", "Bottom", "", true, false},
		{"Tin Can", 1, 9, 0, "You caught a rusty tin can. Not exactly treasure...", "Rusty", "Dented", "Bottom", "", true, false},
		{"Plastic Bottle", 1, 10, 0, "You caught a plastic bottle. Please recycle it!", "Clear", "Crumpled", "Surface", "", true, false},
		{"Seaweed Clump", 1, 8, 0, "Just a tangled clump of seaweed. Nothing to see here.", "Green", "Tangled", "Everywhere", "", true, false},
		{"Driftwood", 3, 8, 1, "A piece of driftwood. Could be useful for crafting?", "Tan", "Weathered", "Surface", "", true, false},
		{"Broken Fishing Rod", 4, 7, 2, "Someone else's broken fishing rod. Unlucky for them!", "Wood", "Broken", "Bottom", "", true, false},
		{"Shopping Bag", 1, 10, 0, "A waterlogged shopping bag. Save the turtles!", "Plastic", "Soggy", "Surface", "", true, false},
		{"Car Tire", 15, 6, 5, "An entire car tire! How did that get here?", "Black", "Rubber", "Bottom", "", true, false},
		{"Waterlogged Phone", 1, 7, 3, "Someone's waterlogged phone. Maybe recoverable?", "Black", "Electronic", "Bottom", "", true, false},
		{"Treasure Chest", 20, 2, 50, "A small treasure chest! It's mostly decorative but worth something!", "Wooden", "Metal-bound", "Deep Bottom", "", true, false},
	}

	// Combine all categories
	allFish := append(regularFish, legendaryFish...)
	return append(allFish, trashItems...)
}

// GetRareFish returns only the rare fish in the game
func GetRareFish() []Fish {
	allFish := GetAllFish()
	rareFish := []Fish{}

	for _, fish := range allFish {
		if fish.Rarity <= 2 && !fish.IsTrash {
			rareFish = append(rareFish, fish)
		}
	}

	return rareFish
}

// GetLegendaryFish returns only the legendary fish in the game
func GetLegendaryFish() []Fish {
	allFish := GetAllFish()
	legendaryFish := []Fish{}

	for _, fish := range allFish {
		if fish.IsLegendary {
			legendaryFish = append(legendaryFish, fish)
		}
	}

	return legendaryFish
}

// GetTrashItems returns only the trash items in the game
func GetTrashItems() []Fish {
	allFish := GetAllFish()
	trashItems := []Fish{}

	for _, fish := range allFish {
		if fish.IsTrash {
			trashItems = append(trashItems, fish)
		}
	}

	return trashItems
}

// GetFishByTimeOfDay returns fish that prefer a specific time of day
func GetFishByTimeOfDay(timeOfDay string) []Fish {
	allFish := GetAllFish()
	timeFish := []Fish{}

	for _, fish := range allFish {
		// Include fish that prefer this time or have no preference
		if fish.PreferredTime == timeOfDay || fish.PreferredTime == "" {
			timeFish = append(timeFish, fish)
		}
	}

	return timeFish
}

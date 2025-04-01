package main

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/user/fishing-game/game"
)

// Fish shapes for visualization - more variety
var fishShapes = []string{
	"><(((°>",
	"><((°>",
	"><(°>",
	"><°>",
	"><>",
	"<><",
	"<·)))><",
	"}}<><",
	"-<><",
	"o><",
	"-=><>",
	"»:>",
	"·.¸¸><(((º>",
	"c=Ø>",
	"[]==O>",
}

// Curved rod frames pointing to the right - more compact version
var fishermanFrames = []string{
	`
    \
     \
      \
       \
        |
~~~~~~~~~
`,
	`
     \
      \
       \
        )
        |
~~~~~~~~~
`,
	`
      \
       \
        )
         )
         |
~~~~~~~~~
`,
}

// Curved rod with caught fish pointing to the right - more compact
var fishermanWithFish = `
     \
      \
       )
        )
        %s
~~~~~~~~~
`

// Generate a fish pattern that scales with terminal width
func generateFishPattern(fish game.Fish) string {
	// Start with base pattern
	var pattern string

	// For very small screens, use minimal fish shapes
	if smallTerminal {
		switch fish.Pattern {
		case "Plain", "Striped", "Spotted", "Mottled":
			pattern = "><>"
		default:
			pattern = "><>"
		}
	} else {
		// Choose basic shape based on pattern, making sizes more consistent
		// and suitable for smaller terminal widths
		if fish.Weight > 100 {
			// Extra large fish - but still compact
			switch fish.Pattern {
			case "Plain":
				pattern = "<=====(Ö>"
			case "Striped":
				pattern = "<####(Ö>"
			case "Spotted":
				pattern = "<•••••(Ö>"
			case "Mottled":
				pattern = "<≈≈≈≈≈(Ö>"
			default:
				pattern = "<=====(Ö>"
			}
		} else if fish.Weight > 30 {
			// Large fish
			switch fish.Pattern {
			case "Plain":
				pattern = "<====(Ö>"
			case "Striped":
				pattern = "<###(Ö>"
			case "Spotted":
				pattern = "<••••(Ö>"
			case "Mottled":
				pattern = "<≈≈≈≈(Ö>"
			default:
				pattern = "<====(Ö>"
			}
		} else if fish.Weight > 10 {
			// Medium fish
			switch fish.Pattern {
			case "Plain":
				pattern = "<===(Ö>"
			case "Striped":
				pattern = "<##(Ö>"
			case "Spotted":
				pattern = "<••(Ö>"
			case "Mottled":
				pattern = "<≈≈(Ö>"
			default:
				pattern = "<===(Ö>"
			}
		} else if fish.Weight > 3 {
			// Small-medium fish
			switch fish.Pattern {
			case "Plain":
				pattern = "<==(Ö>"
			case "Striped":
				pattern = "<#(Ö>"
			case "Spotted":
				pattern = "<•(Ö>"
			case "Mottled":
				pattern = "<≈(Ö>"
			default:
				pattern = "<==(Ö>"
			}
		} else {
			// Tiny fish
			switch fish.Pattern {
			case "Plain":
				pattern = "<=(Ö>"
			case "Striped":
				pattern = "<#(Ö>"
			case "Spotted":
				pattern = "<•(Ö>"
			case "Mottled":
				pattern = "<≈(Ö>"
			default:
				pattern = "<=(Ö>"
			}
		}
	}

	// Set color for the fish based on its color property
	var colorCode string
	switch fish.Color {
	case "Red":
		colorCode = "#FF5555"
	case "Blue":
		colorCode = "#5555FF"
	case "Green":
		colorCode = "#55FF55"
	case "Yellow":
		colorCode = "#FFFF55"
	case "Orange":
		colorCode = "#FFAA55"
	case "Purple":
		colorCode = "#AA55FF"
	case "Pink":
		colorCode = "#FF55AA"
	case "Brown":
		colorCode = "#AA5500"
	case "Gold":
		colorCode = "#FFAA00"
	case "Silver":
		colorCode = "#AAAAAA"
	case "Gray":
		colorCode = "#555555"
	case "Black":
		colorCode = "#000000"
	case "Rainbow":
		colorCode = "#FF9900" // Just use orange for rainbow in terminal
	default:
		colorCode = "#FFFFFF" // Default white
	}

	// Apply color to the pattern
	coloredPattern := lipgloss.NewStyle().Foreground(lipgloss.Color(colorCode)).Render(pattern)
	return coloredPattern
}

// Variable to track if we're in a small terminal
var smallTerminal bool

// Function to update the terminal size check
func updateTerminalSizeCheck(width int) {
	smallTerminal = width < 40
}

// Function to get weather indicator based on weather factor
func getWeatherIndicator() string {
	if weatherFactor > 1.2 {
		// Sunny weather
		return "  \\   /\n   .─.\n  /   \\"
	} else if weatherFactor > 0.9 {
		// Partly cloudy
		return "   \\_\n  _(   )\n (___(__)  "
	} else {
		// Rainy/poor weather
		return "     \n  __//__\n  \\\\//  \n   ||   \n   ||   "
	}
}

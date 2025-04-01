package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/user/fishing-game/game"
)

func (m model) View() string {
	// Update last active time
	lastActiveTime = time.Now()

	var s string

	// Adapt title based on width
	if m.width < 40 {
		s = "🎣 Fishing Game" + "\n"
	} else {
		s = titleStyle.Render(compactGameTitle()) + "\n"
	}

	// Stats are shown in all states - adjust width
	s += renderStats(m.width) + "\n"

	// Content depends on the current state
	switch m.state {
	case "menu":
		s += m.renderMenu()
	case "fishing":
		s += m.renderFishing()
	case "autoFishing":
		s += m.renderAutoFishing()
	case "fishResult":
		s += m.renderFishResult()
	case "inventory":
		s += m.renderInventory()
	case "history":
		s += m.renderHistory()
	case "viewHistoryCatches":
		s += m.renderHistoryCatches()
	}

	// Show message if present
	if m.message != "" {
		s += "\n" + boxStyle.Render(m.message)
	}

	// Simplified help text at bottom
	var helpText string
	if m.state == "menu" {
		helpText = infoStyle.Render("↑↓ | Enter | a:Auto | s:Save | q:Quit")
	} else if m.state == "fishResult" && !autoFishing {
		helpText = infoStyle.Render("Any key: Continue")
	} else if m.state == "inventory" {
		helpText = infoStyle.Render("↑↓:Navigate | h:History | a:Auto | s:Save | q:Back")
	} else if m.state == "history" {
		helpText = infoStyle.Render("↑↓:Navigate | Enter:View | q:Back")
	} else if m.state == "viewHistoryCatches" {
		helpText = infoStyle.Render("↑↓:Navigate | 1-4:Sort | q:Back")
	} else if m.state != "fishResult" {
		helpText = infoStyle.Render("a:Auto | s:Save | q:Back")
	}

	if helpText != "" {
		s += "\n" + helpText
	}

	return s
}

func (m model) renderMenu() string {
	menuBox := strings.Builder{}
	menuBox.WriteString("MENU:\n\n")

	for i, item := range m.menuItems {
		menuLine := fmt.Sprintf("%d. %s", i+1, item)
		if i == m.selectedItem {
			menuBox.WriteString(highlightedMenuItemStyle.Render(menuLine))
		} else {
			menuBox.WriteString(menuItemStyle.Render(menuLine))
		}
		menuBox.WriteString("\n")
	}

	return boxStyle.Render(menuBox.String())
}

func (m model) renderFishing() string {
	// Get the current time period for additional information
	var currentPeriod TimeOfDay
	for _, period := range timePeriods {
		if period.Name == timeOfDay {
			currentPeriod = period
			break
		}
	}

	// Find the animation frame (3 different rod positions)
	fishermanFrame := fishermanFrames[m.fishingState%3]

	// Calculate progress and time remaining
	progress := int(m.fishingProgress * 100)
	remainingTimeMs := int64(float64(m.fishingDuration) * (1.0 - m.fishingProgress))
	remainingSeconds := remainingTimeMs / 1000

	// Content with progress bar
	content := strings.Builder{}

	// Header with time of day information
	timeHeader := fmt.Sprintf("%s %s Fishing", currentPeriod.Icon, timeOfDay)
	timeInfo := fmt.Sprintf("(Catch Rate: %.1fx)", timeFactor)

	// Adjust styles based on terminal width
	if m.width < 40 {
		content.WriteString(accentStyle.Render(timeHeader) + "\n\n")
	} else {
		content.WriteString(accentStyle.Render(timeHeader) + " " + infoStyle.Render(timeInfo) + "\n\n")
		content.WriteString(infoStyle.Render(currentPeriod.Description) + "\n\n")
	}

	// Progress bar
	if progress < 100 {
		timeLabel := ""
		if remainingSeconds > 60 {
			timeLabel = fmt.Sprintf("~%d min remaining", remainingSeconds/60)
		} else {
			timeLabel = fmt.Sprintf("~%d sec remaining", remainingSeconds)
		}

		// Responsive progress bar based on terminal width
		var progressBar string
		if m.width < 30 {
			// Ultra compact for tiny terminals
			progressBar = fmt.Sprintf("[%d%%]", progress)
		} else {
			// Calculate how many blocks to fill based on progress
			barWidth := 20
			if m.width < 60 {
				barWidth = 10
			}

			blocks := int((float64(barWidth) * m.fishingProgress) + 0.5)
			progressBar = "["
			progressBar += strings.Repeat("█", blocks)
			progressBar += strings.Repeat("░", barWidth-blocks)
			progressBar += "]"

			// Add percentage next to the bar
			progressBar = fmt.Sprintf("%s %d%%", progressBar, progress)
		}

		if m.width >= 40 {
			content.WriteString(infoStyle.Render(progressBar) + "\n")
			content.WriteString(infoStyle.Render(timeLabel) + "\n\n")
		} else {
			content.WriteString(infoStyle.Render(progressBar) + "\n\n")
		}
	}

	// Add fishing animation
	if m.message != "" {
		content.WriteString(m.message + "\n\n")
	}
	content.WriteString(fishermanFrame)

	return boxStyle.Render(content.String())
}

func (m model) renderAutoFishing() string {
	content := strings.Builder{}

	// Responsive header based on terminal width
	if m.width < 40 {
		content.WriteString(successStyle.Render("🎣 AUTO-FISHING 🎣") + "\n\n")
	} else {
		// Make auto-fishing status more obvious and compact
		autoFishHeader := successStyle.Render("🎣 AUTO-FISHING ACTIVE 🎣")
		content.WriteString(autoFishHeader + "\n\n")
	}

	// Adapt auto-fishing message based on width
	if m.width < 40 {
		dots := strings.Repeat(".", (m.autoFishTick%3)+1)
		content.WriteString(fmt.Sprintf("Fishing%s\n\n", dots))
	} else {
		content.WriteString(fmt.Sprintf("%s\n\n", m.autoFishMsg))
	}

	// Show animation frame only if there's enough space
	if m.width >= 30 {
		// Show animation frame - more compact
		frameIndex := m.autoFishTick % len(fishermanFrames)
		content.WriteString(fishermanFrames[frameIndex])
	} else {
		// For very small terminals, just show a simple animation
		content.WriteString(fmt.Sprintf("\n%s\n", strings.Repeat("~", m.autoFishTick%3+1)))
	}

	// Show fishing time info
	fishCaughtSoFar := len(player.FishCaught)
	if m.width < 40 {
		content.WriteString(fmt.Sprintf("\n\nFish: %d", fishCaughtSoFar))
	} else {
		// Show fishing time ranges
		if testMode {
			content.WriteString(fmt.Sprintf("\n\nTotal fish: %d | Time per catch: 5-10 seconds (test mode)", fishCaughtSoFar))
		} else {
			content.WriteString(fmt.Sprintf("\n\nTotal fish: %d | Time per catch: 5 sec - 5 min", fishCaughtSoFar))
		}
	}

	return boxStyle.Render(content.String())
}

func (m model) renderFishResult() string {
	content := strings.Builder{}

	// Adjust styles for responsiveness
	resultWidth := 40
	if m.width < 50 {
		resultWidth = m.width - 10 // Adjust for smaller terminals
	}

	if m.catchSuccess {
		// Find the fish details from availableFish
		var fishDetails game.Fish
		for _, fish := range availableFish {
			if fish.Name == m.caughtFish.Name {
				fishDetails = fish
				break
			}
		}

		// Create appropriate header based on fish type
		var catchHeader string
		var headerColor string
		var fishNameColor string

		if fishDetails.IsLegendary {
			// Legendary fish header
			headerColor = "#FF00FF"   // Bright magenta for legendary
			fishNameColor = "#FFFF00" // Bright yellow for legendary names
			catchHeader = "🏆 LEGENDARY CATCH! 🏆"
		} else if fishDetails.IsTrash {
			// Trash item header
			headerColor = "#777777"   // Gray for trash
			fishNameColor = "#AAAAAA" // Light gray for trash names
			catchHeader = "📦 TRASH FOUND 📦"
		} else {
			// Regular fish header
			headerColor = "#008800"   // Green for regular catch
			fishNameColor = "#00FFFF" // Cyan for regular fish names
			catchHeader = "🎣 CATCH! 🎣"
		}

		// Create a more prominent and colorful catch header - responsive width
		formattedHeader := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFF00")).
			Background(lipgloss.Color(headerColor)).
			Padding(0, 2).
			Align(lipgloss.Center).
			Width(resultWidth).
			Render(catchHeader)

		content.WriteString(formattedHeader + "\n\n")

		// Display fish name in larger, centered format - responsive width
		fishName := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(fishNameColor)).
			Width(resultWidth).
			Align(lipgloss.Center).
			Render(m.caughtFish.Name)

		content.WriteString(fishName + "\n\n")

		// Fish info in compact format
		if m.width < 50 {
			content.WriteString(fmt.Sprintf("%dlbs|$%d\n",
				m.caughtFish.Weight, m.caughtFish.Value))
		} else {
			content.WriteString(fmt.Sprintf("Weight: %d lbs | Value: $%d\n",
				m.caughtFish.Weight, m.caughtFish.Value))

			// Add time of day preference if it exists and enough screen space
			if fishDetails.PreferredTime != "" {
				timeInfo := fmt.Sprintf("Most active during: %s", fishDetails.PreferredTime)
				content.WriteString(infoStyle.Render(timeInfo) + "\n")
			}

			// Add special description for legendary or trash items
			if fishDetails.IsLegendary {
				content.WriteString(lipgloss.NewStyle().
					Foreground(lipgloss.Color("#FF00FF")).
					Render("A legendary creature of myth and wonder!") + "\n")
			} else if fishDetails.IsTrash {
				content.WriteString(lipgloss.NewStyle().
					Foreground(lipgloss.Color("#777777")).
					Render("Just some trash from the water...") + "\n")
			}
		}

		// Generate pattern based on fish properties
		fishPattern := generateFishPattern(m.caughtFish)
		fishingLine := fishPattern

		// Show fish graphic only if there's enough space
		if m.width >= 30 {
			// Add spacer for better positioning
			content.WriteString("\n")

			// Show simplified fisherman with fish
			fishermanWithCatch := fmt.Sprintf(fishermanWithFish, fishingLine)
			content.WriteString(fishermanWithCatch + "\n")
		}
	} else {
		// Improved "no catch" message - responsive width
		noCatchMsg := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF5555")).
			Width(resultWidth).
			Align(lipgloss.Center).
			Render("No catch!")

		content.WriteString(noCatchMsg + "\n\n")

		// Show fisherman only if there's enough space
		if m.width >= 30 {
			content.WriteString(fishermanFrames[0])
		}
	}

	// Only show auto-continuing message if auto-fishing is enabled
	if autoFishing {
		content.WriteString("\n" +
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("#AAAAAA")).
				Align(lipgloss.Center).
				Width(resultWidth).
				Italic(true).
				Render("(3s)"))
	}

	return boxStyle.Render(content.String())
}

func (m model) renderInventory() string {
	content := strings.Builder{}

	// Sort indicator and header
	sortLabel := ""
	switch m.inventorySort {
	case "weight":
		sortLabel = "Sorted by Weight (heaviest first)"
	case "value":
		sortLabel = "Sorted by Value (most valuable first)"
	case "quantity":
		sortLabel = "Sorted by Quantity (most caught first)"
	default:
		sortLabel = "Sorted by Rarity (rarest first)"
	}

	// Adapt header based on terminal width with sort info
	if m.width < 40 {
		content.WriteString(successStyle.Render("INVENTORY") + "\n")
		content.WriteString(infoStyle.Render(sortLabel) + "\n\n")
	} else {
		content.WriteString(successStyle.Render("YOUR INVENTORY") + "\n")
		content.WriteString(infoStyle.Render(sortLabel) + "\n\n")
	}

	if len(player.FishCaught) == 0 {
		if m.width < 40 {
			content.WriteString("No fish yet!")
		} else {
			content.WriteString("You haven't caught any fish yet!")
		}
		return boxStyle.Render(content.String())
	}

	// Count fish by type and track totals
	fishCount := make(map[string]int)
	fishTotalWeight := make(map[string]int)
	fishTotalValue := make(map[string]int)

	// Create a slice of unique fish
	type FishSummary struct {
		Name   string
		Count  int
		Weight int
		Value  int
	}

	var fishList []FishSummary

	for _, fish := range player.FishCaught {
		fishCount[fish.Name]++
		fishTotalWeight[fish.Name] += fish.Weight
		fishTotalValue[fish.Name] += fish.Value
	}

	// Convert map to slice for sorting
	for name, count := range fishCount {
		fishList = append(fishList, FishSummary{
			Name:   name,
			Count:  count,
			Weight: fishTotalWeight[name],
			Value:  fishTotalValue[name],
		})
	}

	// Sort the fish list based on the current sort mode
	switch m.inventorySort {
	case "weight":
		// Sort by weight (descending)
		sort.Slice(fishList, func(i, j int) bool {
			return fishList[i].Weight > fishList[j].Weight
		})
	case "value":
		// Sort by value (descending)
		sort.Slice(fishList, func(i, j int) bool {
			return fishList[i].Value > fishList[j].Value
		})
	case "quantity":
		// Sort by quantity (descending)
		sort.Slice(fishList, func(i, j int) bool {
			return fishList[i].Count > fishList[j].Count
		})
	default:
		// Create a map to efficiently look up fish rarity by name
		rarityMap := make(map[string]int)
		for _, f := range availableFish {
			rarityMap[f.Name] = f.Rarity
		}

		// Sort by rarity (most rare first, then alphabetically for same rarity)
		sort.Slice(fishList, func(i, j int) bool {
			// Get rarities from the map
			rarityI := rarityMap[fishList[i].Name]
			rarityJ := rarityMap[fishList[j].Name]

			// Sort by rarity (lower rarity number = more rare)
			if rarityI != rarityJ {
				return rarityI < rarityJ
			}

			// If same rarity, fall back to alphabetical
			return fishList[i].Name < fishList[j].Name
		})
	}

	// Format header based on terminal width
	if m.width >= 70 {
		content.WriteString(fmt.Sprintf("%-20s %-6s %-8s %-8s\n",
			accentStyle.Render("FISH"),
			accentStyle.Render("QTY"),
			accentStyle.Render("WEIGHT"),
			accentStyle.Render("VALUE")))
		content.WriteString(strings.Repeat("─", 46) + "\n")
	} else if m.width >= 40 {
		content.WriteString(fmt.Sprintf("%-12s %-3s %-5s %-5s\n",
			accentStyle.Render("FISH"),
			accentStyle.Render("#"),
			accentStyle.Render("WT"),
			accentStyle.Render("VAL")))
		content.WriteString(strings.Repeat("─", 32) + "\n")
	} else {
		// Super compact for very small terminals
		content.WriteString(fmt.Sprintf("%-5s #  WT  $\n",
			accentStyle.Render("FISH")))
		content.WriteString(strings.Repeat("─", 18) + "\n")
	}

	// Calculate pagination
	totalItems := len(fishList)
	totalPages := (totalItems + m.itemsPerPage - 1) / m.itemsPerPage // Ceiling division

	// Ensure page is in valid range
	if totalItems > 0 && m.inventoryPage >= totalPages {
		m.inventoryPage = totalPages - 1
	}
	if m.inventoryPage < 0 {
		m.inventoryPage = 0
	}

	// Calculate start and end indices for current page
	startIndex := m.inventoryPage * m.itemsPerPage
	endIndex := startIndex + m.itemsPerPage
	if endIndex > totalItems {
		endIndex = totalItems
	}

	// Get the slice of fish for the current page
	currentPageFish := fishList
	if totalItems > 0 {
		currentPageFish = fishList[startIndex:endIndex]
	}

	// Display fish inventory for current page
	for _, fish := range currentPageFish {
		// Find fish details to determine if it's legendary or trash
		var fishDetails game.Fish
		var fishIndicator string

		for _, f := range availableFish {
			if f.Name == fish.Name {
				fishDetails = f
				break
			}
		}

		// Add indicators for special fish types
		if fishDetails.IsLegendary {
			fishIndicator = "🏆 "
		} else if fishDetails.IsTrash {
			fishIndicator = "📦 "
		} else {
			fishIndicator = ""
		}

		if m.width >= 70 {
			// Highlight row based on sort
			fishNameStyle := lipgloss.NewStyle()
			if m.inventorySort == "name" {
				// Use a color based on rarity
				rarityColor := "#FFFFFF" // Default white

				// Special colors for legendary and trash
				if fishDetails.IsLegendary {
					rarityColor = "#FF00FF" // Magenta for legendary
				} else if fishDetails.IsTrash {
					rarityColor = "#777777" // Gray for trash
				} else {
					// Regular color gradient from yellow (rare) to olive (common)
					switch fishDetails.Rarity {
					case 1: // Very rare
						rarityColor = "#FFFF00" // Bright yellow
					case 2:
						rarityColor = "#E6E600"
					case 3, 4:
						rarityColor = "#CCCC00"
					case 5, 6:
						rarityColor = "#B3B300"
					case 7, 8:
						rarityColor = "#999900"
					case 9, 10: // Very common
						rarityColor = "#808000" // Olive
					}
				}
				fishNameStyle = fishNameStyle.Foreground(lipgloss.Color(rarityColor))
			}

			weightStyle := lipgloss.NewStyle()
			if m.inventorySort == "weight" {
				weightStyle = weightStyle.Foreground(lipgloss.Color("#FFFF00"))
			}

			valueStyle := lipgloss.NewStyle()
			if m.inventorySort == "value" {
				valueStyle = valueStyle.Foreground(lipgloss.Color("#FFFF00"))
			}

			quantityStyle := lipgloss.NewStyle()
			if m.inventorySort == "quantity" {
				quantityStyle = quantityStyle.Foreground(lipgloss.Color("#FFFF00"))
			}

			// Fixed-width formatted fish name to ensure alignment (include indicator)
			fishNameFormatted := fmt.Sprintf("%-20.20s", fishIndicator+fish.Name)

			content.WriteString(fmt.Sprintf("%s %s %s %s\n",
				fishNameStyle.Render(fishNameFormatted),
				quantityStyle.Render(fmt.Sprintf("%-6d", fish.Count)),
				weightStyle.Render(fmt.Sprintf("%-8d", fish.Weight)),
				valueStyle.Render(fmt.Sprintf("$%-7d", fish.Value))))
		} else if m.width >= 40 {
			// For medium screens, abbreviate fish names longer than 12 chars
			displayName := fish.Name
			if len(displayName) > 10 { // Account for indicator space
				displayName = displayName[:8] + ".."
			}

			// Add indicator before the name
			displayName = fishIndicator + displayName

			// Fixed-width formatted display name to ensure alignment
			displayNameFormatted := fmt.Sprintf("%-12.12s", displayName)

			// Highlight sorted column
			nameStyle := lipgloss.NewStyle()
			if m.inventorySort == "name" {
				// Use a color based on rarity or special type
				rarityColor := "#FFFFFF" // Default white

				// Special colors for legendary and trash
				if fishDetails.IsLegendary {
					rarityColor = "#FF00FF" // Magenta for legendary
				} else if fishDetails.IsTrash {
					rarityColor = "#777777" // Gray for trash
				} else {
					// Regular color gradient
					switch fishDetails.Rarity {
					case 1: // Very rare
						rarityColor = "#FFFF00" // Bright yellow
					case 2:
						rarityColor = "#E6E600"
					case 3, 4:
						rarityColor = "#CCCC00"
					case 5, 6:
						rarityColor = "#B3B300"
					case 7, 8:
						rarityColor = "#999900"
					case 9, 10: // Very common
						rarityColor = "#808000" // Olive
					}
				}
				nameStyle = nameStyle.Foreground(lipgloss.Color(rarityColor))
			}

			countStyle := lipgloss.NewStyle()
			if m.inventorySort == "quantity" {
				countStyle = countStyle.Foreground(lipgloss.Color("#FFFF00"))
			}

			weightStyle := lipgloss.NewStyle()
			if m.inventorySort == "weight" {
				weightStyle = weightStyle.Foreground(lipgloss.Color("#FFFF00"))
			}

			valueStyle := lipgloss.NewStyle()
			if m.inventorySort == "value" {
				valueStyle = valueStyle.Foreground(lipgloss.Color("#FFFF00"))
			}

			content.WriteString(fmt.Sprintf("%s %s %s %s\n",
				nameStyle.Render(displayNameFormatted),
				countStyle.Render(fmt.Sprintf("%-3d", fish.Count)),
				weightStyle.Render(fmt.Sprintf("%-5d", fish.Weight)),
				valueStyle.Render(fmt.Sprintf("$%-4d", fish.Value))))
		} else {
			// For tiny screens, very compact display
			displayName := fish.Name
			if len(displayName) > 4 { // Even shorter for tiny screens
				displayName = displayName[:2] + ".."
			}

			// For tiny screens, use shorter indicators
			tinyIndicator := ""
			if fishDetails.IsLegendary {
				tinyIndicator = "★"
			} else if fishDetails.IsTrash {
				tinyIndicator = "□"
			}

			displayName = tinyIndicator + displayName

			// Fixed-width formatted for tiny displays
			displayNameFormatted := fmt.Sprintf("%-5.5s", displayName)

			// No highlighting for very small screens - too cluttered
			content.WriteString(fmt.Sprintf("%s %-2d %-3d $%-3d\n",
				displayNameFormatted, fish.Count, fish.Weight, fish.Value))
		}
	}

	// Add pagination info
	if totalPages > 1 {
		content.WriteString("\n")
		pageInfo := fmt.Sprintf("Page %d/%d", m.inventoryPage+1, totalPages)
		if m.width >= 60 {
			content.WriteString(infoStyle.Render(pageInfo + " (↑/↓ to navigate)"))
		} else {
			content.WriteString(infoStyle.Render(pageInfo))
		}
	}

	// Add summary section with total fish caught and value
	content.WriteString("\n")
	if m.width >= 70 {
		content.WriteString(strings.Repeat("─", 46) + "\n")
		content.WriteString(successStyle.Render(fmt.Sprintf("TOTAL: %d fish | Weight: %dlbs | Value: $%d",
			len(player.FishCaught), player.TotalWeight, player.TotalValue)))
	} else if m.width >= 40 {
		content.WriteString(strings.Repeat("─", 32) + "\n")
		content.WriteString(successStyle.Render(fmt.Sprintf("TOTAL: %d fish | $%d",
			len(player.FishCaught), player.TotalValue)))
	} else {
		content.WriteString(strings.Repeat("─", 18) + "\n")
		content.WriteString(successStyle.Render(fmt.Sprintf("TOT:%d|$%d",
			len(player.FishCaught), player.TotalValue)))
	}

	// Add sorting help
	content.WriteString("\n\n")
	if m.width >= 60 {
		content.WriteString(infoStyle.Render("Sort: [1]Rarity [2]Weight [3]Value [4]Qty | Navigate: [↑]Up [↓]Down"))
	} else if m.width >= 30 {
		content.WriteString(infoStyle.Render("[1]Rarity [2]Wt [3]Val [4]Qty | [↑/↓]Nav"))
	} else {
		content.WriteString(infoStyle.Render("1:R 2:W 3:V 4:Q | ↑/↓"))
	}

	return boxStyle.Render(content.String())
}

// A more compact game title
func compactGameTitle() string {
	return `🎣 Fishing Game 🎣`
}

func renderStats(width int) string {
	// Auto-fishing status
	autoStatus := "OFF"
	autoStatusStyle := infoStyle
	if autoFishing {
		autoStatus = "ON"
		autoStatusStyle = successStyle
	}

	// Find the current time period for icon and description
	var currentPeriod TimeOfDay
	for _, period := range timePeriods {
		if period.Name == timeOfDay {
			currentPeriod = period
			break
		}
	}

	// Ultra-compact stats
	statsBuilder := strings.Builder{}

	// Show fish count and auto status
	if width < 40 {
		// Super compact view - just the essentials
		statsBuilder.WriteString(fmt.Sprintf("Fish: %d | %s",
			len(player.FishCaught), currentPeriod.Icon))
	} else if width < 60 {
		// Compact view
		statsBuilder.WriteString(fmt.Sprintf("Fish: %d | Auto: %s | %s %s",
			len(player.FishCaught),
			autoStatusStyle.Render(autoStatus),
			currentPeriod.Icon,
			timeOfDay))
	} else {
		// Full view with all details
		timeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#88CCFF"))

		// Show more details in widescreen
		statsBuilder.WriteString(fmt.Sprintf("Fish: %d | Auto: %s | %s %s",
			len(player.FishCaught),
			autoStatusStyle.Render(autoStatus),
			currentPeriod.Icon,
			timeStyle.Render(timeOfDay+" - "+currentPeriod.Description)))

		// Show test mode status on wider displays
		if testMode {
			statsBuilder.WriteString(" | " + infoStyle.Render("TEST MODE"))
		}

		// Show catch factor if there's space
		catchInfo := fmt.Sprintf(" | Catch Rate: %.1fx", timeFactor)
		statsBuilder.WriteString(infoStyle.Render(catchInfo))
	}

	return boxStyle.Render(statsBuilder.String())
}

// renderHistory displays a list of available dates with fishing records
func (m model) renderHistory() string {
	content := strings.Builder{}

	// Header
	content.WriteString(historyHeaderStyle.Render("FISHING HISTORY") + "\n\n")

	if len(m.historyDates) == 0 {
		content.WriteString("No fishing records found yet.")
		return boxStyle.Render(content.String())
	}

	// Display available dates in a calendar-like format
	for i, date := range m.historyDates {
		// Parse date for formatting
		dateObj, err := time.Parse("2006-01-02", date)
		if err != nil {
			// If we can't parse it, just show the raw date
			dateStr := date
			if i == m.historyDateIndex {
				content.WriteString(highlightedDateStyle.Render(dateStr))
			} else {
				content.WriteString(dateStyle.Render(dateStr))
			}
		} else {
			// Format date in a more readable format
			dateStr := dateObj.Format("Mon, Jan 2 2006")

			// Get fish count for this date
			count, weight, value := getFishCaughtDetails(date)
			dateInfo := fmt.Sprintf("%s (%d fish, %dlbs, $%d)", dateStr, count, weight, value)

			// Highlight selected date
			if i == m.historyDateIndex {
				content.WriteString(highlightedDateStyle.Render(dateInfo))
			} else {
				content.WriteString(dateStyle.Render(dateStr) + fmt.Sprintf(" - %d fish, %dlbs, $%d", count, weight, value))
			}
		}
		content.WriteString("\n")
	}

	// Add pagination info if needed
	if len(m.historyDates) > 10 {
		content.WriteString("\n")
		pageInfo := fmt.Sprintf("Page %d/%d", (m.historyDateIndex/10)+1, (len(m.historyDates)+9)/10)
		content.WriteString(infoStyle.Render(pageInfo))
	}

	return boxStyle.Render(content.String())
}

// renderHistoryCatches displays fish caught on a specific date
func (m model) renderHistoryCatches() string {
	content := strings.Builder{}

	// Parse date for pretty formatting
	dateObj, err := time.Parse("2006-01-02", viewingDate)
	var dateHeader string
	if err != nil {
		dateHeader = viewingDate
	} else {
		dateHeader = dateObj.Format("Monday, January 2, 2006")
	}

	// Determine if this is today or a past date
	isToday := viewingDate == time.Now().Format("2006-01-02")

	// Header shows the date
	if isToday {
		content.WriteString(historyHeaderStyle.Render("TODAY'S CATCHES - "+dateHeader) + "\n\n")
	} else {
		content.WriteString(historyHeaderStyle.Render("PAST CATCHES - "+dateHeader) + "\n\n")
	}

	// Get fish caught on this date
	fishCaught := getFishCaughtOnDate(viewingDate)

	if len(fishCaught) == 0 {
		content.WriteString("No fish caught on this date.")
		return boxStyle.Render(content.String())
	}

	// Count fish by type and track totals
	fishCount := make(map[string]int)
	fishTotalWeight := make(map[string]int)
	fishTotalValue := make(map[string]int)

	// Create a slice of unique fish
	type FishSummary struct {
		Name   string
		Count  int
		Weight int
		Value  int
	}

	var fishList []FishSummary

	for _, fish := range fishCaught {
		fishCount[fish.Name]++
		fishTotalWeight[fish.Name] += fish.Weight
		fishTotalValue[fish.Name] += fish.Value
	}

	// Convert map to slice for sorting
	for name, count := range fishCount {
		fishList = append(fishList, FishSummary{
			Name:   name,
			Count:  count,
			Weight: fishTotalWeight[name],
			Value:  fishTotalValue[name],
		})
	}

	// Sort the fish list based on the current sort mode
	switch m.inventorySort {
	case "weight":
		// Sort by weight (descending)
		sort.Slice(fishList, func(i, j int) bool {
			return fishList[i].Weight > fishList[j].Weight
		})
	case "value":
		// Sort by value (descending)
		sort.Slice(fishList, func(i, j int) bool {
			return fishList[i].Value > fishList[j].Value
		})
	case "quantity":
		// Sort by quantity (descending)
		sort.Slice(fishList, func(i, j int) bool {
			return fishList[i].Count > fishList[j].Count
		})
	default:
		// Create a map to efficiently look up fish rarity by name
		rarityMap := make(map[string]int)
		for _, f := range availableFish {
			rarityMap[f.Name] = f.Rarity
		}

		// Sort by rarity (most rare first, then alphabetically for same rarity)
		sort.Slice(fishList, func(i, j int) bool {
			// Get rarities from the map
			rarityI := rarityMap[fishList[i].Name]
			rarityJ := rarityMap[fishList[j].Name]

			// Sort by rarity (lower rarity number = more rare)
			if rarityI != rarityJ {
				return rarityI < rarityJ
			}

			// If same rarity, fall back to alphabetical
			return fishList[i].Name < fishList[j].Name
		})
	}

	// Format header based on terminal width
	if m.width >= 70 {
		content.WriteString(fmt.Sprintf("%-20s %-6s %-8s %-8s\n",
			accentStyle.Render("FISH"),
			accentStyle.Render("QTY"),
			accentStyle.Render("WEIGHT"),
			accentStyle.Render("VALUE")))
		content.WriteString(strings.Repeat("─", 46) + "\n")
	} else if m.width >= 40 {
		content.WriteString(fmt.Sprintf("%-12s %-3s %-5s %-5s\n",
			accentStyle.Render("FISH"),
			accentStyle.Render("#"),
			accentStyle.Render("WT"),
			accentStyle.Render("VAL")))
		content.WriteString(strings.Repeat("─", 32) + "\n")
	} else {
		// Super compact for very small terminals
		content.WriteString(fmt.Sprintf("%-5s #  WT  $\n",
			accentStyle.Render("FISH")))
		content.WriteString(strings.Repeat("─", 18) + "\n")
	}

	// Calculate pagination
	totalItems := len(fishList)
	totalPages := (totalItems + m.itemsPerPage - 1) / m.itemsPerPage // Ceiling division

	// Ensure page is in valid range
	if totalItems > 0 && m.inventoryPage >= totalPages {
		m.inventoryPage = totalPages - 1
	}
	if m.inventoryPage < 0 {
		m.inventoryPage = 0
	}

	// Calculate start and end indices for current page
	startIndex := m.inventoryPage * m.itemsPerPage
	endIndex := startIndex + m.itemsPerPage
	if endIndex > totalItems {
		endIndex = totalItems
	}

	// Get the slice of fish for the current page
	currentPageFish := fishList
	if totalItems > 0 {
		currentPageFish = fishList[startIndex:endIndex]
	}

	// Display fish inventory for current page
	for _, fish := range currentPageFish {
		// Find fish details to determine if it's legendary or trash
		var fishDetails game.Fish
		var fishIndicator string

		for _, f := range availableFish {
			if f.Name == fish.Name {
				fishDetails = f
				break
			}
		}

		// Add indicators for special fish types
		if fishDetails.IsLegendary {
			fishIndicator = "🏆 "
		} else if fishDetails.IsTrash {
			fishIndicator = "📦 "
		} else {
			fishIndicator = ""
		}

		if m.width >= 70 {
			// Highlight row based on sort
			fishNameStyle := lipgloss.NewStyle()
			if m.inventorySort == "name" {
				// Use a color based on rarity
				rarityColor := "#FFFFFF" // Default white

				// Special colors for legendary and trash
				if fishDetails.IsLegendary {
					rarityColor = "#FF00FF" // Magenta for legendary
				} else if fishDetails.IsTrash {
					rarityColor = "#777777" // Gray for trash
				} else {
					// Regular color gradient from yellow (rare) to olive (common)
					switch fishDetails.Rarity {
					case 1: // Very rare
						rarityColor = "#FFFF00" // Bright yellow
					case 2:
						rarityColor = "#E6E600"
					case 3, 4:
						rarityColor = "#CCCC00"
					case 5, 6:
						rarityColor = "#B3B300"
					case 7, 8:
						rarityColor = "#999900"
					case 9, 10: // Very common
						rarityColor = "#808000" // Olive
					}
				}
				fishNameStyle = fishNameStyle.Foreground(lipgloss.Color(rarityColor))
			}

			weightStyle := lipgloss.NewStyle()
			if m.inventorySort == "weight" {
				weightStyle = weightStyle.Foreground(lipgloss.Color("#FFFF00"))
			}

			valueStyle := lipgloss.NewStyle()
			if m.inventorySort == "value" {
				valueStyle = valueStyle.Foreground(lipgloss.Color("#FFFF00"))
			}

			quantityStyle := lipgloss.NewStyle()
			if m.inventorySort == "quantity" {
				quantityStyle = quantityStyle.Foreground(lipgloss.Color("#FFFF00"))
			}

			// Fixed-width formatted fish name to ensure alignment (include indicator)
			fishNameFormatted := fmt.Sprintf("%-20.20s", fishIndicator+fish.Name)

			content.WriteString(fmt.Sprintf("%s %s %s %s\n",
				fishNameStyle.Render(fishNameFormatted),
				quantityStyle.Render(fmt.Sprintf("%-6d", fish.Count)),
				weightStyle.Render(fmt.Sprintf("%-8d", fish.Weight)),
				valueStyle.Render(fmt.Sprintf("$%-7d", fish.Value))))
		} else if m.width >= 40 {
			// For medium screens, abbreviate fish names longer than 12 chars
			displayName := fish.Name
			if len(displayName) > 10 { // Account for indicator space
				displayName = displayName[:8] + ".."
			}

			// Add indicator before the name
			displayName = fishIndicator + displayName

			// Fixed-width formatted display name to ensure alignment
			displayNameFormatted := fmt.Sprintf("%-12.12s", displayName)

			// Highlight sorted column
			nameStyle := lipgloss.NewStyle()
			if m.inventorySort == "name" {
				// Use a color based on rarity or special type
				rarityColor := "#FFFFFF" // Default white

				// Special colors for legendary and trash
				if fishDetails.IsLegendary {
					rarityColor = "#FF00FF" // Magenta for legendary
				} else if fishDetails.IsTrash {
					rarityColor = "#777777" // Gray for trash
				} else {
					// Regular color gradient
					switch fishDetails.Rarity {
					case 1: // Very rare
						rarityColor = "#FFFF00" // Bright yellow
					case 2:
						rarityColor = "#E6E600"
					case 3, 4:
						rarityColor = "#CCCC00"
					case 5, 6:
						rarityColor = "#B3B300"
					case 7, 8:
						rarityColor = "#999900"
					case 9, 10: // Very common
						rarityColor = "#808000" // Olive
					}
				}
				nameStyle = nameStyle.Foreground(lipgloss.Color(rarityColor))
			}

			countStyle := lipgloss.NewStyle()
			if m.inventorySort == "quantity" {
				countStyle = countStyle.Foreground(lipgloss.Color("#FFFF00"))
			}

			weightStyle := lipgloss.NewStyle()
			if m.inventorySort == "weight" {
				weightStyle = weightStyle.Foreground(lipgloss.Color("#FFFF00"))
			}

			valueStyle := lipgloss.NewStyle()
			if m.inventorySort == "value" {
				valueStyle = valueStyle.Foreground(lipgloss.Color("#FFFF00"))
			}

			content.WriteString(fmt.Sprintf("%s %s %s %s\n",
				nameStyle.Render(displayNameFormatted),
				countStyle.Render(fmt.Sprintf("%-3d", fish.Count)),
				weightStyle.Render(fmt.Sprintf("%-5d", fish.Weight)),
				valueStyle.Render(fmt.Sprintf("$%-4d", fish.Value))))
		} else {
			// For tiny screens, very compact display
			displayName := fish.Name
			if len(displayName) > 4 { // Even shorter for tiny screens
				displayName = displayName[:2] + ".."
			}

			// For tiny screens, use shorter indicators
			tinyIndicator := ""
			if fishDetails.IsLegendary {
				tinyIndicator = "★"
			} else if fishDetails.IsTrash {
				tinyIndicator = "□"
			}

			displayName = tinyIndicator + displayName

			// Fixed-width formatted for tiny displays
			displayNameFormatted := fmt.Sprintf("%-5.5s", displayName)

			// No highlighting for very small screens - too cluttered
			content.WriteString(fmt.Sprintf("%s %-2d %-3d $%-3d\n",
				displayNameFormatted, fish.Count, fish.Weight, fish.Value))
		}
	}

	// Add pagination info
	if totalPages > 1 {
		content.WriteString("\n")
		pageInfo := fmt.Sprintf("Page %d/%d", m.inventoryPage+1, totalPages)
		if m.width >= 60 {
			content.WriteString(infoStyle.Render(pageInfo + " (↑/↓ to navigate)"))
		} else {
			content.WriteString(infoStyle.Render(pageInfo))
		}
	}

	// Add summary section with total fish caught and value
	content.WriteString("\n")

	// Calculate total weight and value for this date
	totalFish := len(fishCaught)
	var totalWeight, totalValue int
	for _, fish := range fishCaught {
		totalWeight += fish.Weight
		totalValue += fish.Value
	}

	if m.width >= 70 {
		content.WriteString(strings.Repeat("─", 46) + "\n")
		content.WriteString(successStyle.Render(fmt.Sprintf("%s TOTAL: %d fish | Weight: %dlbs | Value: $%d",
			dateObj.Format("Jan 2"), totalFish, totalWeight, totalValue)))
	} else if m.width >= 40 {
		content.WriteString(strings.Repeat("─", 32) + "\n")
		content.WriteString(successStyle.Render(fmt.Sprintf("%s: %d fish | $%d",
			dateObj.Format("Jan 2"), totalFish, totalValue)))
	} else {
		content.WriteString(strings.Repeat("─", 18) + "\n")
		content.WriteString(successStyle.Render(fmt.Sprintf("%s: %d | $%d",
			dateObj.Format("01/02"), totalFish, totalValue)))
	}

	// Add sorting help
	content.WriteString("\n\n")
	if m.width >= 60 {
		content.WriteString(infoStyle.Render("Sort: [1]Rarity [2]Weight [3]Value [4]Qty | Navigate: [↑]Up [↓]Down"))
	} else if m.width >= 30 {
		content.WriteString(infoStyle.Render("[1]Rarity [2]Wt [3]Val [4]Qty | [↑/↓]Nav"))
	} else {
		content.WriteString(infoStyle.Render("1:R 2:W 3:V 4:Q | ↑/↓"))
	}

	return boxStyle.Render(content.String())
}

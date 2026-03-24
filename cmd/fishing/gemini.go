package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	geminiAPIURL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent"
	geminiTimeout = 8 * time.Second
)

// geminiRequest is the request body sent to the Gemini REST API
type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

// geminiResponse holds the parsed Gemini API response
type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// hasGeminiToken returns true when a GEMINI_API_KEY env variable is set
func hasGeminiToken() bool {
	return os.Getenv("GEMINI_API_KEY") != ""
}

// generateFishContent calls Gemini to produce a unique ASCII-art fish and a one-line catch message.
// Both return values are empty strings when the API key is absent or the call fails.
func generateFishContent(fishName, fishColor string, weight int, isLegendary bool) (art string, catchMsg string) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", ""
	}

	legendaryNote := ""
	if isLegendary {
		legendaryNote = " (a legendary mythical creature)"
	}

	prompt := fmt.Sprintf(
		`For a terminal fishing game I caught a %s%s (%s coloured, %d lbs).
Please provide TWO things:
1. Unique ASCII art of this fish using only standard printable ASCII characters. Maximum 3 lines tall, maximum 25 characters wide. The fish must face right. Be creative – use characters like < > ( ) = ~ * ° # - | Ø o.
2. A single exciting catch announcement for this specific fish. Maximum 70 characters. Be vivid and unique.

Format your response EXACTLY like this (no extra text, no markdown):
ART:<your ascii art, can span multiple lines>
MSG:<your catch announcement>`,
		fishName, legendaryNote, fishColor, weight,
	)

	reqBody := geminiRequest{
		Contents: []geminiContent{
			{Parts: []geminiPart{{Text: prompt}}},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", ""
	}

	ctx, cancel := context.WithTimeout(context.Background(), geminiTimeout)
	defer cancel()

	url := fmt.Sprintf("%s?key=%s", geminiAPIURL, apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return "", ""
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", ""
	}

	var geminiResp geminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return "", ""
	}

	if len(geminiResp.Candidates) == 0 ||
		len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", ""
	}

	return parseGeminiResponse(geminiResp.Candidates[0].Content.Parts[0].Text)
}

// parseGeminiResponse extracts the art block and catch message from the structured Gemini output.
func parseGeminiResponse(text string) (art string, catchMsg string) {
	var artLines []string
	inArt := false

	for _, line := range strings.Split(text, "\n") {
		switch {
		case strings.HasPrefix(line, "ART:"):
			inArt = true
			rest := strings.TrimSpace(strings.TrimPrefix(line, "ART:"))
			if rest != "" {
				artLines = append(artLines, rest)
			}
		case strings.HasPrefix(line, "MSG:"):
			inArt = false
			catchMsg = strings.TrimSpace(strings.TrimPrefix(line, "MSG:"))
		default:
			if inArt && line != "" {
				artLines = append(artLines, line)
			}
		}
	}

	if len(artLines) > 0 {
		art = strings.TrimSpace(strings.Join(artLines, "\n"))
	}
	return art, catchMsg
}

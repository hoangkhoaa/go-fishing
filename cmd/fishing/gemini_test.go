package main

import "testing"

// TestParseGeminiResponse verifies that the structured Gemini output is correctly split
// into ASCII art and catch message portions.
func TestParseGeminiResponse(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantArt     string
		wantCatchMsg string
	}{
		{
			name: "single-line art",
			input: "ART:><(((°>\nMSG:You caught a magnificent fish!",
			wantArt:      "><(((°>",
			wantCatchMsg: "You caught a magnificent fish!",
		},
		{
			name: "multi-line art",
			input: "ART:><(((°>\n   ><°>\nMSG:A rare one!",
			wantArt:      "><(((°>\n   ><°>",
			wantCatchMsg: "A rare one!",
		},
		{
			name: "empty response",
			input: "",
			wantArt:      "",
			wantCatchMsg: "",
		},
		{
			name: "art only, no message",
			input: "ART:><(((°>",
			wantArt:      "><(((°>",
			wantCatchMsg: "",
		},
		{
			name: "message only",
			input: "MSG:What a catch!",
			wantArt:      "",
			wantCatchMsg: "What a catch!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotArt, gotMsg := parseGeminiResponse(tt.input)
			if gotArt != tt.wantArt {
				t.Errorf("art: got %q, want %q", gotArt, tt.wantArt)
			}
			if gotMsg != tt.wantCatchMsg {
				t.Errorf("catchMsg: got %q, want %q", gotMsg, tt.wantCatchMsg)
			}
		})
	}
}

// TestCheckAndUpdatePersonalBest verifies that personal bests are tracked correctly.
func TestCheckAndUpdatePersonalBest(t *testing.T) {
	// Reset personal bests for isolated test
	personalBests = make(map[string]int)

	// First catch – always a new record
	if !checkAndUpdatePersonalBest("Trout", 5) {
		t.Error("expected first catch to be a new record")
	}
	if personalBests["Trout"] != 5 {
		t.Errorf("expected personal best 5, got %d", personalBests["Trout"])
	}

	// Lighter catch – not a new record
	if checkAndUpdatePersonalBest("Trout", 3) {
		t.Error("expected lighter catch not to be a new record")
	}
	if personalBests["Trout"] != 5 {
		t.Errorf("expected personal best to remain 5, got %d", personalBests["Trout"])
	}

	// Heavier catch – new record
	if !checkAndUpdatePersonalBest("Trout", 8) {
		t.Error("expected heavier catch to be a new record")
	}
	if personalBests["Trout"] != 8 {
		t.Errorf("expected personal best 8, got %d", personalBests["Trout"])
	}

	// Different species is independent
	if !checkAndUpdatePersonalBest("Bass", 12) {
		t.Error("expected first Bass catch to be a new record")
	}
}

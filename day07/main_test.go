package main

import "testing"

func Test_hand_SetType(t *testing.T) {
	tests := []struct {
		cards      string
		withJokers bool
		want       string
	}{
		{"A99AA", true, "full house"},
		{"TTTTJ", true, "five of a kind"},
		{"TTTJJ", true, "five of a kind"},
		{"TTJJJ", true, "five of a kind"},
		{"TJJJJ", true, "five of a kind"},
		{"JJJJJ", true, "five of a kind"},
		{"QTTTJ", true, "four of a kind"},
		{"QTTJJ", true, "four of a kind"},
		{"QTJJJ", true, "four of a kind"},
		{"Q3TTJ", true, "three of a kind"},
		{"Q3TJJ", true, "three of a kind"},
		{"QQTTJ", true, "full house"},
		{"QQT3J", true, "three of a kind"},
		{"22445", true, "two pair"},
		{"2345J", true, "one pair"},
	}
	for _, tt := range tests {
		t.Run(tt.cards, func(t *testing.T) {
			h := NewHand(tt.cards, "0", tt.withJokers)
			if h.typ.String() != tt.want {
				t.Errorf("hand.SetType() = %v, want %v", h.typ, tt.want)
			}
		})
	}
}

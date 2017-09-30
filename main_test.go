package main

import (
	"bytes"
	"testing"
)

func TestParseFile(t *testing.T) {
	buf := bytes.NewBufferString(`Animal Rights
128
Kick 10001000
Snare 00001000
HiHat 00100010
`)
	song := parseFile(buf)
	if song.Name != "Animal Rights" {
		t.Errorf("expected song name to be Animal Rights but found %s", song.Name)
	}
	if song.bpm != 128 {
		t.Errorf("expected song bpm to be 128 but found %d", song.bpm)
	}
	if len(song.Patterns) != 3 {
		t.Errorf("expected song to have 3 patterns but found %d", len(song.Patterns))
	}
	expectedPatterns := map[string]string{
		"Kick":  "10001000",
		"Snare": "00001000",
		"HiHat": "00100010",
	}
	for _, pattern := range song.Patterns {
		_, ok := expectedPatterns[pattern.Name]
		if !ok {
			t.Error("expected pattern to exist but it doesn't")
		}
	}
}
func TestPlay(t *testing.T) {
	patterns := []*Pattern{
		{"Kick", "10001000"},
		{"Snare", "00001000"},
		{"HiHat", "00100010"},
	}
	song := &Song{
		Name:     "Animal Rights",
		bpm:      128,
		Patterns: patterns,
	}

	ch := make(chan string)
	done := make(chan struct{}, 0)
	go song.Play(ch, done)
	var buffer bytes.Buffer
	for i := 0; i < 16; i++ {
		buffer.WriteString(<-ch)
	}
	done <- struct{}{}
	close(ch)

	expectedString := "Kick_HiHat_Kick+Snare_HiHat_Kick_HiHat_Kick+Snare_HiHat_"
	result := buffer.String()
	if result != expectedString {
		t.Errorf("expected string to equal Kick_HiHat_Kick+Snare_HiHat_Kick_HiHat_Kick+Snare_HiHat_ but found %s", result)
	}
}

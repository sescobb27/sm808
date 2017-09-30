package main

import (
	"bytes"
	"testing"
	"time"
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
	if song.maxSteps != 8 {
		t.Errorf("expected song to have 8 steps but found %d", song.maxSteps)
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
	t.Run("8 steps", func(t *testing.T) {
		song1Patterns := []*Pattern{
			{"Kick", "10001000"},
			{"Snare", "00001000"},
			{"HiHat", "00100010"},
		}
		song1 := &Song{
			Name:     "Animal Rights",
			bpm:      128,
			Patterns: song1Patterns,
			maxSteps: 8,
			ticker:   time.NewTicker(time.Microsecond * 1),
		}
		ch := make(chan string)
		done := make(chan struct{}, 0)
		go song1.Play(ch, done)
		var buffer bytes.Buffer
		for i := 0; i < 16; i++ {
			buffer.WriteString(<-ch)
		}
		done <- struct{}{}
		close(ch)

		result := buffer.String()
		if result != "Kick_HiHat_Kick+Snare_HiHat_Kick_HiHat_Kick+Snare_HiHat_" {
			t.Errorf("expected string to equal Kick_HiHat_Kick+Snare_HiHat_Kick_HiHat_Kick+Snare_HiHat_ but found %s", result)
		}
	})

	t.Run("mix 8 steps with 16 steps", func(t *testing.T) {
		song2Patterns := []*Pattern{
			{"Kick", "10001000"},
			{"Snare", "00001000"},
			{"HiHat", "00100010"},
			{"Bass", "1000100010001000"},
		}
		song2 := &Song{
			Name:     "Beep Boop",
			bpm:      128,
			Patterns: song2Patterns,
			maxSteps: 16,
			ticker:   time.NewTicker(time.Microsecond * 1),
		}
		ch := make(chan string)
		done := make(chan struct{}, 0)
		go song2.Play(ch, done)
		var buffer bytes.Buffer
		for i := 0; i < 16; i++ {
			buffer.WriteString(<-ch)
		}
		done <- struct{}{}
		close(ch)

		result := buffer.String()
		if result != "Kick+Bass_HiHat_Kick+Snare+Bass_HiHat_Kick+Bass_HiHat_Kick+Snare+Bass_HiHat_" {
			t.Errorf("expected string to equal Kick+Bass_HiHat_Kick+Snare+Bass_HiHat_Kick+Bass_HiHat_Kick+Snare+Bass_HiHat_ but found %s", result)
		}
	})
}

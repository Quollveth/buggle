package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type song struct {
	name   string
	artist string
	path   string
}

type station struct {
	name string
	url  string
}

type model struct {
	connected bool // are we connected to the daemon
	quitting  bool // showing "are ya sure to quit" screen

	saved struct {
		songs    []song
		stations []station
	}

	playing struct {
		song    song
		station station
	}
}

// let everything initialize to zero value
func initModel() model { return model{} }

func createFakeData() model {
	return model{
		connected: true,
		quitting:  false,
		saved: struct {
			songs    []song
			stations []station
		}{
			songs: []song{
				{
					name:   "Song 1",
					artist: "John Music",
					path:   "music/song1.ogg",
				},
				{
					name:   "Midnight Ride",
					artist: "The Wanderers",
					path:   "music/midnight_ride.ogg",
				},
				{
					name:   "Echoes of Time",
					artist: "Luna Nova",
					path:   "music/echoes_of_time.ogg",
				},
				{
					name:   "Coffeehouse Blues",
					artist: "Ella and the Jazz Cats",
					path:   "music/coffeehouse_blues.ogg",
				},
				{
					name:   "Digital Dreams",
					artist: "SynthRider",
					path:   "music/digital_dreams.ogg",
				},
			},

			stations: []station{
				{
					name: "Radio 1",
					url:  "http://radio1.com/stream",
				},
				{
					name: "Indie Vibes",
					url:  "http://indievibes.fm/stream",
				},
				{
					name: "Jazz Lounge",
					url:  "http://jazzlounge.fm/stream",
				},
				{
					name: "Synthwave Central",
					url:  "http://synthcentral.net/stream",
				},
				{
					name: "Classic Rock Live",
					url:  "http://classicrocklive.fm/stream",
				},
			},
		},

		playing: struct {
			song    song
			station station
		}{
			song: song{
				name:   "Song 1",
				artist: "John Music",
				path:   "music/song1.ogg",
			},
			station: station{
				name: "Radio 1",
				url:  "http://radio1.com/stream",
			},
		},
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string { return "" }

func main() {
	p := tea.NewProgram(initModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("peepee went poopoo: %v", err)
		os.Exit(1)
	}
}

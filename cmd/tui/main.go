package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/bubbles/paginator"
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

const NUM_TABS = 3
const (
	TAB_STATIONS int = iota
	TAB_SONGS
	TAB_PLACEHOLDER
)

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

	//-- Config
	styles styles

	//-- UI data
	activeTab int
	tabs      [NUM_TABS]string
	tabView   [NUM_TABS]func(model) string

	paginator paginator.Model
}

func initModel(config config) model {
	songs := []song{
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
	}

	stations := []station{
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
	}

	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 2
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
	p.SetTotalPages(len(songs))

	tabViews := [NUM_TABS]func(model) string{}

	tabViews[TAB_STATIONS] = stationsTabView
	tabViews[TAB_SONGS] = songTabView
	tabViews[TAB_PLACEHOLDER] = placeholderView

	return model{
		connected: true,
		quitting:  false,
		saved: struct {
			songs    []song
			stations []station
		}{
			songs:    songs,
			stations: stations,
		},
		playing: struct {
			song    song
			station station
		}{
			song:    songs[0],
			station: stations[0],
		},

		styles: createStyles(config.colors),

		paginator: p,
		//tabs:      [NUM_TABS]string{"Stations", "Songs", "Placeholder"},
		tabs:      [NUM_TABS]string{"tab 1", "tab 2", "tab 3"},
		tabView:   tabViews,
		activeTab: 0,
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit
		case "l":
			m.paginator.NextPage()
		case "h":
			m.paginator.PrevPage()
		case "tab":
			if m.activeTab == len(m.tabs)-1 {
				m.activeTab = 0
			} else {
				m.activeTab++
			}
		case "shift+tab":
			if m.activeTab == 0 {
				m.activeTab = len(m.tabs) - 1
			} else {
				m.activeTab--
			}
		}
	}

	return m, nil
}


//-- Renders


func placeholderView(m model) string { return "nothing here" }

func stationsTabView(m model) string {
	var b strings.Builder
	start, end := m.paginator.GetSliceBounds(len(m.saved.stations))
	for _, item := range m.saved.stations[start:end] {
		b.WriteString("  • " + item.name + "\n\n")
	}
	b.WriteString("  " + m.paginator.View())

	return b.String()
}

func songTabView(m model) string {
	var b strings.Builder
	start, end := m.paginator.GetSliceBounds(len(m.saved.songs))
	for _, item := range m.saved.songs[start:end] {
		b.WriteString("  • " + item.name + "\n\n")
	}
	b.WriteString("  " + m.paginator.View())

	return b.String()
}

func (m model) renderTabs() string {
	var renderedTabs []string

	for i, t := range m.tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.tabs)-1, i == m.activeTab
		if isActive {
			style = m.styles.activeTab
		} else {
			style = m.styles.inactiveTab
		}

		border, _, _, _, _ := style.GetBorder()

		if isActive && isFirst {
			border.TopLeft = "┐"
		} else if isActive && isLast {
			border.TopRight = "┬"
		} else if !isActive && isFirst {
			border.TopLeft = "┬"
		} else if !isActive && isLast {
			border.TopRight = "┬"
		}

		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}

func (m model) View() string {
	allTabs := m.renderTabs()

	borderFinish := lipgloss.NewStyle().Foreground(m.styles.outerWindow.GetForeground()).Render(strings.Repeat("─", m.styles.outerWindow.GetWidth()-lipgloss.Width(allTabs)))

	row := lipgloss.JoinHorizontal(lipgloss.Top, allTabs, borderFinish)

	return m.styles.outerWindow.Render(row)
}

func main() {
	p := tea.NewProgram(initModel(createDefaultConfig()))

	if _, err := p.Run(); err != nil {
		fmt.Printf("peepee went poopoo: %v", err)
		os.Exit(1)
	}
}

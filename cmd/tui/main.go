package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"

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
	winHeight, winWidth int

	activeTab int
	tabs      [NUM_TABS]string
	tabView   [NUM_TABS]func(model) string

	paginator paginator.Model
}

func initModel(config config) model {
	winWidth, winHeight, _ := term.GetSize(os.Stdin.Fd())

	winWidth -= 2
	winHeight -= 2

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

	styles := createStyles(config.colors, winWidth, winHeight)

	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 2
	p.ActiveDot = lipgloss.NewStyle().Foreground(config.colors.listActiveDot).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(config.colors.listInactiveDot).Render("•")
	p.SetTotalPages(len(songs))

	tabViews := [NUM_TABS]func(model) string{}

	tabViews[TAB_STATIONS] = stationsTabView
	tabViews[TAB_SONGS] = songTabView
	tabViews[TAB_PLACEHOLDER] = placeholderView

	return model{
		connected: true,
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
		winWidth:  winWidth,
		winHeight: winHeight,

		styles:    styles,
		paginator: p,
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

//
//-- Renders
//

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
		active := i == m.activeTab

		if active {
			style = m.styles.activeTab
		} else {
			style = m.styles.inactiveTab
		}

		border, _, _, _, _ := style.GetBorder()

		if i == 0 {
			if active {
				border.BottomLeft = "│"
			} else {
				border.BottomLeft = "├"
			}
		}

		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}

func (m model) renderPlaying() string {
	return "now playing goes here"
}

//
//-- Main Render
//

func (m model) View() string {
	allTabs := m.renderTabs()

	borderLength := m.winWidth - lipgloss.Width(allTabs) + 1
	borderLine := strings.Repeat("─", borderLength)
	borderLine = "\n\n" + borderLine + "╮"

	styledBorder := lipgloss.NewStyle().Foreground(m.styles.activeTab.GetBorderTopForeground()).Render(borderLine)

	tabsRow := lipgloss.JoinHorizontal(lipgloss.Top, allTabs, styledBorder)

	//-- build

	var b strings.Builder

	playing := m.renderPlaying()

	b.WriteString(lipgloss.JoinVertical(
		lipgloss.Left,
		tabsRow,

		m.styles.middleWindow.Height(
			m.winHeight-(lipgloss.Height(tabsRow)+lipgloss.Height(playing)),
		).Render(m.tabView[m.activeTab](m)),

		m.styles.bottomWindow.Render(playing),
	))

	return b.String()
}

func main() {
	p := tea.NewProgram(initModel(createDefaultConfig()))

	if _, err := p.Run(); err != nil {
		fmt.Printf("peepee went poopoo: %v", err)
		os.Exit(1)
	}
}

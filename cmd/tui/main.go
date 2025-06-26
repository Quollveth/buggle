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
	name     string
	artist   string
	album    string
	path     string
	duration int
}

type station struct {
	name string
	desc string
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

	activeTab    int
	tabs         [NUM_TABS]string
	tabView      [NUM_TABS]func(model) string
	selectedItem int

	paginator paginator.Model
}

func initModel(config config) model {
	winWidth, winHeight, _ := term.GetSize(os.Stdin.Fd())

	winHeight -= 2
	winHeight -= 2

	songs := []song{}

	for i := range 50 {
		songs = append(songs, song{
			name:   fmt.Sprint("Song ", i),
			artist: "John Music",
			album:  "Music - The Sequel",
		})
	}

	stations := []station{}

	for i := range 75 {
		stations = append(stations, station{
			name: fmt.Sprint("Station ", i),
			desc: fmt.Sprint("Placeholder test station n ", i),
		})
	}

	styles := createStyles(config.colors, winWidth, winHeight)

	p := paginator.New()
	p.Type = paginator.Dots
	p.ActiveDot = lipgloss.NewStyle().Foreground(config.colors.listActiveDot).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(config.colors.listInactiveDot).Render("•")
	// defaults to 1 to avoid division by zero panic, will be set to proper value before first render
	p.PerPage = 1
	p.SetTotalPages(1)

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
		tabs:      [NUM_TABS]string{"Stations", "Songs", "Placeholder"},
		tabView:   tabViews,
		activeTab: 2,
	}
}

func (m model) Init() tea.Cmd { return nil }

// these are technically constants but i can't actually make them const
var SONG_HEIGHT int = lipgloss.Height(renderSongItem(song{}, false, styles{})) - 1
var STATION_HEIGHT int = lipgloss.Height(renderStationItem(station{}, false, styles{})) - 1

func (m *model) updatePaginator() {
	m.paginator.Page = 0
	m.selectedItem = 0

	switch m.activeTab {
	case TAB_SONGS:
		// add one extra line for the paginator dots             ↓
		m.paginator.PerPage = (m.styles.middleWindow.GetHeight() - 1) / SONG_HEIGHT
		m.paginator.SetTotalPages(len(m.saved.songs))
	case TAB_STATIONS:
		m.paginator.PerPage = (m.styles.middleWindow.GetHeight() - 1) / STATION_HEIGHT
		m.paginator.SetTotalPages(len(m.saved.stations))
	}
}

func (m *model) updateSizes() {
	tabsRow := m.renderTabLine()
	playing := m.renderBottomWindow()
	middleWindowHeight := m.winHeight - (lipgloss.Height(tabsRow) + lipgloss.Height(playing))

	m.styles.middleWindow = m.styles.middleWindow.Height(middleWindowHeight)
	m.styles.leftPane = m.styles.leftPane.Height(middleWindowHeight)
	m.styles.rightPane = m.styles.rightPane.Height(middleWindowHeight)

	m.updatePaginator()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit

		case "l":
			m.selectedItem = 0
			m.paginator.NextPage()

		case "h":
			m.selectedItem = 0
			m.paginator.PrevPage()

		case "k":
			if m.selectedItem > 0 {
				m.selectedItem -= 1
			}

		case "j":
			if m.selectedItem < m.paginator.PerPage-1 {
				m.selectedItem += 1
			}

		case "tab":
			if m.activeTab == len(m.tabs)-1 {
				m.activeTab = 0
			} else {
				m.activeTab++
			}
			m.updatePaginator()

		case "shift+tab":
			if m.activeTab == 0 {
				m.activeTab = len(m.tabs) - 1
			} else {
				m.activeTab--
			}
			m.updatePaginator()
		}
	}

	return m, nil
}

//
//-- Renders
//

func renderSongItem(song song, selected bool, styles styles) string {
	var b strings.Builder

	var bullet string
	if selected {
		bullet = ">"
	} else {
		bullet = "•"
	}

	// lipgloss will stagger all the text if the newline is the last character but not if it's the first
	b.WriteString(styles.textPrimary.Render("   "+bullet, song.name))
	b.WriteString(styles.textSecondary.Render("\n    ", song.artist))
	b.WriteString(styles.textSecondary.Render("\n    ", song.album))
	b.WriteString("\n\n")

	return b.String()
}

func renderStationItem(station station, selected bool, styles styles) string {
	var b strings.Builder

	var bullet string
	if selected {
		bullet = ">"
	} else {
		bullet = "•"
	}

	// lipgloss will stagger all the text if the newline is the last character but not if it's the first
	b.WriteString(styles.textPrimary.Render("   "+bullet, station.name))
	b.WriteString(styles.textSecondary.Render("\n    ", station.desc))
	b.WriteString("\n\n")

	return b.String()
}

//

func placeholderView(m model) string { return "nothing here" }

func stationsTabView(m model) string {
	var b strings.Builder
	start, end := m.paginator.GetSliceBounds(len(m.saved.stations))
	for idx, item := range m.saved.stations[start:end] {
		b.WriteString(renderStationItem(item, m.selectedItem == idx, m.styles))
	}
	b.WriteString("  " + m.paginator.View())

	return b.String()
}

func songTabView(m model) string {
	var b strings.Builder
	start, end := m.paginator.GetSliceBounds(len(m.saved.songs))
	for idx, item := range m.saved.songs[start:end] {
		b.WriteString(renderSongItem(item, m.selectedItem == idx, m.styles))
	}
	b.WriteString(m.paginator.View())

	return b.String()
}

//

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

func (m model) renderBottomWindow() string {
	/*
	str := fmt.Sprintf(
		"terminal width: %v | middle window width: %v | expected bottom window width: %v | real bottom window width: %v",
		m.winWidth,
		lipgloss.Width(m.renderMiddleWindow()),
		m.styles.bottomWindow.GetWidth(),
		lipgloss.Width(m.styles.bottomWindow.Render("")),
	)

	rendered := m.styles.bottomWindow.Render(str)
	*/

	return m.styles.bottomWindow.Render("")
}

//
//-- Main Render
//

func (m model) renderTabLine() string {
	allTabs := m.renderTabs()

	borderLength := m.styles.leftPane.GetWidth() - lipgloss.Width(allTabs) + 1
	if borderLength < 0 {
		return ""
	}

	borderLine := strings.Repeat("─", borderLength)
	borderLine = "\n\n" + borderLine + "┬─"
	borderLine += strings.Repeat("─", m.styles.rightPane.GetWidth()) + "╮"

	styledBorder := lipgloss.NewStyle().Foreground(m.styles.activeTab.GetBorderTopForeground()).Render(borderLine)

	return lipgloss.JoinHorizontal(lipgloss.Left, allTabs, styledBorder)
}

func (m model) renderMiddleWindow() string {
	tab := m.styles.leftPane.Render(m.tabView[m.activeTab](m))
	info := m.styles.rightPane.Render("info text here")

	return m.styles.middleWindow.Render(lipgloss.JoinHorizontal(lipgloss.Left, tab, info))
}

func (m model) View() string {
	tabsRow := m.renderTabLine()

	var b strings.Builder

	b.WriteString(lipgloss.JoinVertical(
		lipgloss.Top,
		tabsRow,
		m.renderMiddleWindow(),
		m.renderBottomWindow(),
	))

	return b.String()
}

func main() {
	// these are constants and are only here so i don't have to count, they should never be zero
	if SONG_HEIGHT == 0 {
		panic("\x1b[1;91mSong renderer size is zero\x1b[1;0m\n")
	}
	if STATION_HEIGHT == 0 {
		panic("\x1b[1;91mStation renderer size is zero\x1b[1;0m\n")
	}

	model := initModel(createDefaultConfig())
	model.updateSizes()

	p := tea.NewProgram(model)

	if _, err := p.Run(); err != nil {
		fmt.Printf("\x1b[1;91mpeepee went poopoo:\x1b[1;0m %v", err)
		os.Exit(1)
	}
}

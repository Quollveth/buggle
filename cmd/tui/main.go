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
	firstRender         bool
	winHeight, winWidth int
	middleWinHeight     int

	activeTab    int
	tabs         [NUM_TABS]string
	tabView      [NUM_TABS]func(model) string
	selectedItem int

	paginator paginator.Model
}

func initModel(config config) model {
	winWidth, winHeight, _ := term.GetSize(os.Stdin.Fd())

	winWidth -= 2
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
		})
	}

	styles := createStyles(config.colors, winWidth, winHeight)

	p := paginator.New()
	p.Type = paginator.Dots
	p.ActiveDot = lipgloss.NewStyle().Foreground(config.colors.listActiveDot).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(config.colors.listInactiveDot).Render("•")

	p.PerPage = 3
	p.SetTotalPages(0)

	tabViews := [NUM_TABS]func(model) string{}

	tabViews[TAB_STATIONS] = stationsTabView
	tabViews[TAB_SONGS] = songTabView
	tabViews[TAB_PLACEHOLDER] = placeholderView

	return model{
		firstRender: true,
		connected:   true,
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
		winWidth:        winWidth,
		winHeight:       winHeight,
		middleWinHeight: 0,

		styles:    styles,
		paginator: p,
		tabs:      [NUM_TABS]string{"Stations", "Songs", "Placeholder"},
		tabView:   tabViews,
		activeTab: 1,
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m *model) updatePaginator() {
	m.paginator.Page = 0
	m.selectedItem = 0

	switch m.activeTab {
	case TAB_STATIONS:
		m.paginator.PerPage = 9
		m.paginator.SetTotalPages(len(m.saved.stations))
	case TAB_SONGS:
		m.paginator.PerPage = 4
		m.paginator.SetTotalPages(len(m.saved.songs))
	}
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

func renderSong(song song, selected bool, styles styles) string {
	var b strings.Builder

	// lipgloss will stagger all the text if the newline is the last character but not if it's the first
	var bullet string
	if selected {
		bullet = ">"
	} else {
		bullet = "•"
	}

	b.WriteString(styles.textPrimary.Render("   "+bullet, song.name))
	b.WriteString(styles.textSecondary.Render("\n    ", song.artist))
	b.WriteString(styles.textSecondary.Render("\n    ", song.album))
	b.WriteString("\n\n")

	return b.String()
}

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
	for idx, item := range m.saved.songs[start:end] {
		b.WriteString(renderSong(item, m.selectedItem == idx, m.styles))
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
	return fmt.Sprintf("tab: %v | per page: %v | total pages: %v", m.activeTab, m.paginator.PerPage, m.paginator.TotalPages)
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

	playing := m.renderPlaying()

	if m.firstRender {
		m.middleWinHeight = m.winHeight - (lipgloss.Height(tabsRow) + lipgloss.Height(playing))
		m.firstRender = false
	}

	//-- build

	var b strings.Builder

	b.WriteString(lipgloss.JoinVertical(
		lipgloss.Left,
		tabsRow,

		m.styles.middleWindow.Height(m.middleWinHeight).Render(m.tabView[m.activeTab](m)),
		m.styles.bottomWindow.Render(playing),
	))

	return b.String()
}

func main() {
	model := initModel(createDefaultConfig())

	p := tea.NewProgram(model)

	if _, err := p.Run(); err != nil {
		fmt.Printf("peepee went poopoo: %v", err)
		os.Exit(1)
	}
}

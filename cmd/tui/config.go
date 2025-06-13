package main

import (
	"github.com/charmbracelet/lipgloss"
)

type colors struct {
	outBorder,
	hightlight,
	listActiveDot,
	listInactiveDot lipgloss.Color
}

type config struct {
	colors colors
}

func createDefaultConfig() config {
	return config{

		colors: colors{
			outBorder:       lipgloss.Color("#00FFFF"),
			hightlight:      lipgloss.Color("#FFFF00"),
			listActiveDot:   lipgloss.Color("252"),
			listInactiveDot: lipgloss.Color("238"),
		},
	}
}

type styles struct {
	middleWindow, bottomWindow,
	inactiveTab, activeTab lipgloss.Style
}

func borderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

func createStyles(colors colors, width, height int) styles {
	inactiveTabBorder := borderWithBottom("┴", "─", "┴")
	activeTabBorder := borderWithBottom("┘", " ", "└")
	inactiveTab := lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(colors.outBorder).Padding(0, 1)
	activeTab := inactiveTab.Border(activeTabBorder, true)

	bottomWindow := lipgloss.NewStyle().Width(width).Border(lipgloss.RoundedBorder()).BorderForeground(colors.outBorder).UnsetBorderTop()
	middleWindow := bottomWindow.UnsetBorderBottom()

	return styles{
		inactiveTab:  inactiveTab,
		activeTab:    activeTab,
		middleWindow: middleWindow,
		bottomWindow: bottomWindow,
	}
}

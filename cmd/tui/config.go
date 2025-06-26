package main

import (
	"github.com/charmbracelet/lipgloss"
)

type colors struct {
	outBorder,
	hightlight,
	listActiveDot,
	textPrimary,
	textSecondary,
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
			textPrimary:     lipgloss.Color("15"),
			textSecondary:   lipgloss.Color("245"),
		},
	}
}

type styles struct {
	middleWindow, bottomWindow,
	leftPane, rightPane,
	textPrimary, textSecondary,
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

	middleWindow := lipgloss.NewStyle().Width(width)
	bottomWindow := middleWindow.Border(lipgloss.RoundedBorder()).BorderForeground(colors.outBorder).UnsetBorderTop()

	textPrimary := lipgloss.NewStyle().Foreground(colors.textPrimary)
	textSecondary := lipgloss.NewStyle().Foreground(colors.textSecondary)

	rightPane := lipgloss.NewStyle().Width((width/2)-2).Border(borderWithBottom("┴", "─", "┤")).BorderForeground(colors.outBorder).UnsetBorderTop().UnsetBorderLeft()
	leftPane  := rightPane.Border(borderWithBottom("├", "─", "┴")).UnsetBorderTop()

	return styles{
		inactiveTab:   inactiveTab,
		activeTab:     activeTab,
		middleWindow:  middleWindow,
		bottomWindow:  bottomWindow,
		textPrimary:   textPrimary,
		textSecondary: textSecondary,
		leftPane: leftPane,
		rightPane: rightPane,
	}
}

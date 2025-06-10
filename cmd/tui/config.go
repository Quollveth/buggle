package main

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
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
			outBorder:       lipgloss.Color("#FFFF00"),
			hightlight:      lipgloss.Color("#00FFFF"),
			listActiveDot:   lipgloss.Color("252"),
			listInactiveDot: lipgloss.Color("238"),
		},
	}
}

type styles struct {
	outerWindow,
	inactiveTab,
	activeTab lipgloss.Style
}

func borderWithTop(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.TopLeft = left
	border.Top = middle
	border.TopRight = right
	return border
}

func createStyles(colors colors) styles {
	winWidth, winHeight, _ := term.GetSize(os.Stdin.Fd())

	outBorder := borderWithTop("┌", " ", "┐")
	outStyle := lipgloss.NewStyle().Width(winWidth-2).Height(winHeight-2).Border(outBorder, true).BorderForeground(colors.outBorder).UnsetBorderTop()

	inactiveTabBorder := borderWithTop("┬", "─", "┬")
	activeTabBorder := borderWithTop("┐", " ", "┌")
	inactiveTabStyle := lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(colors.hightlight).Padding(0, 1)
	activeTabStyle := inactiveTabStyle.Border(activeTabBorder, true)

	return styles{
		outerWindow: outStyle,
		inactiveTab: inactiveTabStyle,
		activeTab:   activeTabStyle,
	}
}

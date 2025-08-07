package main

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	primaryColor   = lipgloss.Color("#FFD700")
	secondaryColor = lipgloss.Color("#FFA500")
	accentColor    = lipgloss.Color("#FFFF00")
	successColor   = lipgloss.Color("#32CD32")
	errorColor     = lipgloss.Color("#FF6347")
	textColor      = lipgloss.Color("#FFFFFF")
	dimTextColor   = lipgloss.Color("#CCCCCC")
	bgColor        = lipgloss.Color("#1A1A1A")
	borderColor    = lipgloss.Color("#444444")

	// titleStyle = lipgloss.NewStyle().
	// 		Foreground(primaryColor).
	// 		Bold(true).
	// 		Padding(0, 1).
	// 		Border(lipgloss.RoundedBorder()).
	// 		BorderForeground(primaryColor)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true).
			MarginTop(1)

	normalTextStyle = lipgloss.NewStyle().
			Foreground(textColor)

	dimTextStyle = lipgloss.NewStyle().
			Foreground(dimTextColor)

	successStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(0, 1).
			MarginBottom(1)

	buttonStyle = lipgloss.NewStyle().
			Foreground(bgColor).
			Background(primaryColor).
			Padding(0, 2).
			MarginRight(1).
			Bold(true)

	selectedButtonStyle = lipgloss.NewStyle().
				Foreground(bgColor).
				Background(accentColor).
				Padding(0, 2).
				MarginRight(1).
				Bold(true)

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1, 2).
			MarginBottom(1)

	headerStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Align(lipgloss.Center).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(primaryColor).
			Padding(0, 2)

	statusActiveStyle = lipgloss.NewStyle().
				Foreground(successColor).
				Bold(true)

	statusInactiveStyle = lipgloss.NewStyle().
				Foreground(dimTextColor)
)

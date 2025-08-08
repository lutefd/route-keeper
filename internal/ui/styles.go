package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	primaryColor   = lipgloss.AdaptiveColor{Light: "#FFC107", Dark: "#FFD700"}
	secondaryColor = lipgloss.AdaptiveColor{Light: "#FFB74D", Dark: "#FFA000"}
	accentColor    = lipgloss.AdaptiveColor{Light: "#FFD54F", Dark: "#FFC107"}
	successColor   = lipgloss.AdaptiveColor{Light: "#4CAF50", Dark: "#81C784"}
	errorColor     = lipgloss.AdaptiveColor{Light: "#F44336", Dark: "#E57373"}
	textColor      = lipgloss.AdaptiveColor{Light: "#333333", Dark: "#E0E0E0"}
	dimTextColor   = lipgloss.AdaptiveColor{Light: "#757575", Dark: "#9E9E9E"}
	bgColor        = lipgloss.AdaptiveColor{Light: "#FAFAFA", Dark: "#1E1E1E"}
	borderColor    = lipgloss.AdaptiveColor{Light: "#E0E0E0", Dark: "#424242"}
	highlightColor = lipgloss.AdaptiveColor{Light: "#FFF9C4", Dark: "#2A2A2A"}

	houstonNormal = `
   ╭─────────╮
   │  ◕   ◕  │
   │   ───   │
   ╰─────────╯
`

	houstonHappy = `
   ╭─────────╮
   │  ◕   ◕  │
   │   ╰─╯   │
   ╰─────────╯
`

	houstonThinking = `
   ╭─────────╮
   │  -   -  │
   │   ───   │
   ╰─────────╯
`

	houstonStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Align(lipgloss.Center)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Italic(true).
			MarginTop(1).
			MarginBottom(1)

	normalTextStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Background(lipgloss.NoColor{})

	dimTextStyle = lipgloss.NewStyle().
			Foreground(dimTextColor).
			Faint(true)

	successStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), false, false, false, false).
			BorderBottom(true).
			BorderForeground(primaryColor).
			Padding(0, 1).
			MarginBottom(1).
			Width(50).
			Foreground(textColor).
			Background(bgColor)

	focusedInputStyle = inputStyle.Copy().
				BorderForeground(accentColor).
				Background(highlightColor).
				Foreground(textColor).
				Bold(true)

	buttonStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(0, 2).
			MarginRight(2).
			Bold(true).
			Background(lipgloss.NoColor{})

	selectedButtonStyle = buttonStyle.Copy().
				Foreground(lipgloss.Color("#000000")).
				Background(accentColor).
				BorderForeground(accentColor)

	headerStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Align(lipgloss.Center).
			MarginBottom(1).
			Underline(true).
			UnderlineSpaces(true)

	statusActiveStyle = lipgloss.NewStyle().
				Foreground(successColor).
				Bold(true)

	statusInactiveStyle = lipgloss.NewStyle().
				Foreground(dimTextColor).
				Faint(true)
)

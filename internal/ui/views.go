package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func (m *MainModel) mainMenuView() string {
	var houston string
	switch m.MenuIndex {
	case 0:
		houston = houstonStyle.Render(houstonNormal)
	case 1:
		houston = houstonStyle.Foreground(accentColor).Render(houstonHappy)
	default:
		houston = houstonStyle.Render(houstonThinking)
	}

	header := headerStyle.Render("🚀 ROUTE KEEPER")
	subtitle := subtitleStyle.Render("Houston, we have connectivity!")

	menuItems := []string{
		"Select Profile",
		"Create New Profile",
		"Quit",
	}

	maxWidth := 0
	for _, item := range menuItems {
		if len(item) > maxWidth {
			maxWidth = len(item)
		}
	}
	menuWidth := maxWidth + 8

	var menuStrings []string
	for i, item := range menuItems {
		prefix := "  "
		if i == m.MenuIndex {
			prefix = "→ "
		}
		paddedItem := fmt.Sprintf("%-*s", maxWidth, item)
		if i == m.MenuIndex {
			menuStrings = append(menuStrings, selectedButtonStyle.Width(menuWidth).Render(prefix+paddedItem))
		} else {
			menuStrings = append(menuStrings, buttonStyle.Width(menuWidth).Render(prefix+paddedItem))
		}
	}

	menu := lipgloss.JoinVertical(lipgloss.Left, menuStrings...)

	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().Padding(0, 4, 0, 0).Render(houston),
		lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			"",
			subtitle,
			"",
			menu,
			"",
			dimTextStyle.Render("Use ↑/↓ to navigate • Enter to select • q to quit"),
		),
	)

	return lipgloss.NewStyle().
		Padding(2, 4).
		Render(content)
}

func (m *MainModel) profileListView() string {
	header := headerStyle.Render("📋 SELECT PROFILE")
	profiles := m.ProfilesManager.GetProfiles()

	if len(profiles) == 0 {
		empty := lipgloss.NewStyle().
			Foreground(primaryColor).
			Italic(true).
			Render("No profiles found")

		instructions := lipgloss.JoinVertical(
			lipgloss.Left,
			"",
			dimTextStyle.Render("Press 'c' to create a new profile"),
			dimTextStyle.Render("or press Esc to go back"),
		)

		content := lipgloss.JoinVertical(
			lipgloss.Center,
			header,
			"",
			empty,
			"",
			instructions,
		)

		return lipgloss.NewStyle().
			Padding(2, 4).
			Render(content)
	}

	var profileItems []string
	for i, profile := range profiles {
		status := statusInactiveStyle.Render("○")
		if i == m.ProfileIndex {
			status = statusActiveStyle.Render("●")
		}

		url := profile.GetFullURL()
		interval := fmt.Sprintf("⏱  every %d min", profile.Interval)

		profileCard := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), false, false, false, false).
			BorderLeft(true).
			BorderForeground(primaryColor).
			Padding(0, 2).
			Margin(0, 0, 1, 0).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Left,
					lipgloss.JoinHorizontal(
						lipgloss.Top,
						status,
						" ",
						lipgloss.NewStyle().
							Bold(i == m.ProfileIndex).
							Foreground(primaryColor).
							Render(profile.Name),
					),
					lipgloss.NewStyle().
						MarginLeft(2).
						Render(dimTextStyle.Render(url)),
					lipgloss.NewStyle().
						MarginLeft(2).
						Render(dimTextStyle.Render(interval)),
				),
			)

		if i == m.ProfileIndex {
			profileCard = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(primaryColor).
				BorderTop(true).
				BorderRight(true).
				BorderBottom(true).
				BorderLeft(true).
				Padding(1, 2).
				Margin(0, 0, 1, 0).
				Render(profileCard)
		}

		profileItems = append(profileItems, profileCard)
	}

	instructions := lipgloss.JoinHorizontal(
		lipgloss.Left,
		dimTextStyle.Render("Enter: Run"),
		lipgloss.NewStyle().Margin(0, 2).Render("•"),
		dimTextStyle.Render("e: Edit"),
		lipgloss.NewStyle().Margin(0, 2).Render("•"),
		dimTextStyle.Render("d: Delete"),
		lipgloss.NewStyle().Margin(0, 2).Render("•"),
		dimTextStyle.Render("c: Create New"),
		lipgloss.NewStyle().Margin(0, 2).Render("•"),
		dimTextStyle.Render("Esc: Back"),
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		lipgloss.JoinVertical(lipgloss.Left, profileItems...),
		"",
		instructions,
	)

	return lipgloss.NewStyle().
		Padding(2, 4).
		MaxWidth(80).
		Render(content)
}

func (m *MainModel) profileFormView(title string) string {
	header := headerStyle.Render("⚙️  " + title)

	fields := []struct {
		label       string
		description string
	}{
		{"Profile Name", "A name to identify this profile"},
		{"Base URL", "The base URL to monitor (e.g., https://api.example.com)"},
		{"Route", "The API endpoint route (e.g., /health)"},
		{"URL Params", "Optional query parameters (e.g., key1=value1&key2=value2)"},
		{"Headers", "Request headers (e.g., Authorization=Bearer token)"},
		{"Interval (minutes)", "How often to check the endpoint (minimum 1 minute)"},
	}

	var formFields []string
	for i, field := range fields {
		isFocused := i == m.InputIndex
		var labelStyle lipgloss.Style
		if isFocused {
			labelStyle = normalTextStyle.Copy().Bold(true).Foreground(primaryColor)
		} else {
			labelStyle = normalTextStyle
		}
		inputField := m.Inputs[i].View()
		inputStyleToUse := inputStyle
		if isFocused {
			inputStyleToUse = focusedInputStyle
		}
		formField := lipgloss.JoinVertical(
			lipgloss.Left,
			labelStyle.Render(field.label),
			dimTextStyle.Italic(true).Render(field.description),
			inputStyleToUse.Render(inputField),
		)
		formFields = append(formFields, formField)
	}

	saveButtonLabel := "Save Profile"
	if m.InputIndex == len(fields) {
		saveButtonLabel = "💾 " + saveButtonLabel
	}
	saveButton := buttonStyle.Render(saveButtonLabel)
	if m.InputIndex == len(fields) {
		saveButton = selectedButtonStyle.Render(saveButtonLabel)
	}

	formContent := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		lipgloss.JoinVertical(
			lipgloss.Left,
			append(formFields, "", lipgloss.NewStyle().Align(lipgloss.Center).Render(saveButton))...,
		),
	)

	instructions := lipgloss.JoinHorizontal(
		lipgloss.Left,
		dimTextStyle.Render("Tab/↑↓: Navigate"),
		lipgloss.NewStyle().Margin(0, 2).Render("•"),
		dimTextStyle.Render("Enter: Next/Save"),
		lipgloss.NewStyle().Margin(0, 2).Render("•"),
		dimTextStyle.Render("Esc: Cancel"),
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		formContent,
		"",
		instructions,
	)

	return lipgloss.NewStyle().
		Padding(2, 4).
		MaxWidth(80).
		Render(content)
}

func (m *MainModel) runningView() string {
	header := headerStyle.Render("🔄 MONITORING")

	var status string
	if m.IsRunning {
		status = lipgloss.JoinHorizontal(
			lipgloss.Left,
			statusActiveStyle.Render("●"),
			" ",
			statusActiveStyle.Render("ACTIVE - Monitoring endpoint..."),
		)
	} else {
		status = lipgloss.JoinHorizontal(
			lipgloss.Left,
			statusInactiveStyle.Render("●"),
			" ",
			statusInactiveStyle.Render("PAUSED - Monitoring paused"),
		)
	}

	profileCard := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 2).
		Margin(1, 0, 2, 0).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.NewStyle().Bold(true).Render(m.CurrentProfile.Name),
				"",
				dimTextStyle.Render("URL:"),
				normalTextStyle.Render(m.CurrentProfile.GetFullURL()),
				"",
				lipgloss.JoinHorizontal(
					lipgloss.Left,
					dimTextStyle.Render("Interval:"),
					" ",
					normalTextStyle.Render(fmt.Sprintf("%d minutes", m.CurrentProfile.Interval)),
				),
			),
		)

	var resultsView string
	if len(m.PingResults) > 0 {
		resultsHeader := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), false, false, true, false).
			BorderForeground(borderColor).
			Padding(0, 0, 1, 0).
			Margin(0, 0, 1, 0).
			Render("📊 Recent Pings")

		var resultLines []string
		for i, result := range m.PingResults {
			if i >= 5 {
				break
			}
			timestamp := result.Timestamp.Format("15:04:05")
			var statusIcon, statusText string
			if result.Error == nil && result.StatusCode >= 200 && result.StatusCode < 300 {
				statusIcon = successStyle.Render("✓")
				statusText = successStyle.Render(fmt.Sprintf("HTTP %d", result.StatusCode))
			} else {
				statusIcon = errorStyle.Render("✗")
				if result.Error != nil {
					statusText = errorStyle.Render("ERROR: " + result.Error.Error())
				} else {
					statusText = errorStyle.Render(fmt.Sprintf("HTTP %d", result.StatusCode))
				}
			}
			duration := dimTextStyle.Render(fmt.Sprintf("(%v)", result.Duration.Truncate(time.Millisecond)))
			resultLine := lipgloss.JoinHorizontal(
				lipgloss.Left,
				dimTextStyle.Render(timestamp),
				"  ",
				statusIcon,
				" ",
				statusText,
				" ",
				duration,
			)
			resultLines = append(resultLines, resultLine)
		}
		resultsView = lipgloss.JoinVertical(
			lipgloss.Left,
			append([]string{resultsHeader}, resultLines...)...,
		)
	} else {
		resultsView = dimTextStyle.Italic(true).Render("No ping results yet...")
	}

	var instructions string
	if m.IsRunning {
		instructions = lipgloss.JoinHorizontal(
			lipgloss.Left,
			dimTextStyle.Render("s: Stop"),
			lipgloss.NewStyle().Margin(0, 2).Render("•"),
			dimTextStyle.Render("Esc/q: Exit"),
		)
	} else {
		instructions = lipgloss.JoinHorizontal(
			lipgloss.Left,
			dimTextStyle.Render("s: Start"),
			lipgloss.NewStyle().Margin(0, 2).Render("•"),
			dimTextStyle.Render("Esc/q: Exit"),
		)
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		status,
		"",
		profileCard,
		"",
		resultsView,
		"",
		instructions,
	)

	return lipgloss.NewStyle().
		Padding(2, 4).
		MaxWidth(80).
		Render(content)
}

package main

import (
	"fmt"
	"strings"
	"time"
)

func (m MainModel) mainMenuView() string {
	header := headerStyle.Render("🚀 ROUTE KEEPER")
	subtitle := subtitleStyle.Render("Houston, we have connectivity!")

	menuItems := []string{
		"Select Profile",
		"Create New Profile",
		"Quit",
	}

	var menuStrings []string
	for i, item := range menuItems {
		if i == m.menuIndex {
			menuStrings = append(menuStrings, selectedButtonStyle.Render("▶ "+item))
		} else {
			menuStrings = append(menuStrings, buttonStyle.Render("  "+item))
		}
	}

	menu := strings.Join(menuStrings, "\n")

	instructions := dimTextStyle.Render("\nUse ↑/↓ to navigate, Enter to select, q to quit")

	content := fmt.Sprintf("%s\n\n%s\n\n%s%s", header, subtitle, menu, instructions)

	return panelStyle.Render(content)
}

func (m MainModel) profileListView() string {
	header := headerStyle.Render("📋 SELECT PROFILE")

	profiles := m.profilesManager.GetProfiles()

	if len(profiles) == 0 {
		empty := errorStyle.Render("No profiles found!")
		instructions := dimTextStyle.Render("\nPress 'c' to create a new profile, or Esc to go back")
		return panelStyle.Render(fmt.Sprintf("%s\n\n%s%s", header, empty, instructions))
	}

	var profileList []string
	for i, profile := range profiles {
		status := statusInactiveStyle.Render("●")
		prefix := "  "

		if i == m.profileIndex {
			status = statusActiveStyle.Render("●")
			prefix = "▶ "
		}

		url := profile.GetFullURL()
		interval := fmt.Sprintf("every %d min", profile.Interval)

		profileInfo := fmt.Sprintf("%s%s%s\n   %s %s",
			prefix,
			status,
			normalTextStyle.Render(" "+profile.Name),
			dimTextStyle.Render(url),
			dimTextStyle.Render("("+interval+")"),
		)

		profileList = append(profileList, profileInfo)
	}

	list := strings.Join(profileList, "\n\n")
	instructions := dimTextStyle.Render("\n\nEnter: Run • e: Edit • d: Delete • c: Create New • Esc: Back")

	content := fmt.Sprintf("%s\n\n%s%s", header, list, instructions)

	return panelStyle.Render(content)
}

func (m MainModel) profileFormView(title string) string {
	header := headerStyle.Render("⚙️  " + title)

	fields := []string{
		"Profile Name:",
		"Base URL:",
		"Route:",
		"URL Params:",
		"Headers:",
		"Interval (minutes):",
	}

	var formFields []string
	for i, field := range fields {
		fieldStyle := normalTextStyle
		if i == m.inputIndex {
			fieldStyle = successStyle
		}

		formFields = append(formFields,
			fieldStyle.Render(field)+"\n"+
				inputStyle.Render(m.inputs[i].View()),
		)
	}

	form := strings.Join(formFields, "\n")

	saveButton := buttonStyle.Render("Save Profile")
	if m.inputIndex == len(m.inputs)-1 {
		saveButton = selectedButtonStyle.Render("Save Profile")
	}

	instructions := dimTextStyle.Render("\nTab/↑↓: Navigate • Enter: Next/Save • Esc: Cancel")

	content := fmt.Sprintf("%s\n\n%s\n\n%s%s", header, form, saveButton, instructions)

	return panelStyle.Render(content)
}

func (m MainModel) runningView() string {
	header := headerStyle.Render("🔄 KEEPING ALIVE")

	profileInfo := fmt.Sprintf(
		"%s\n%s\n%s",
		normalTextStyle.Render("Profile: "+m.currentProfile.Name),
		normalTextStyle.Render("URL: "+m.currentProfile.GetFullURL()),
		normalTextStyle.Render(fmt.Sprintf("Interval: %d minutes", m.currentProfile.Interval)),
	)

	var status string
	if m.isRunning {
		status = statusActiveStyle.Render("● ACTIVE - Keeping route alive...")
	} else {
		status = statusInactiveStyle.Render("● STOPPED")
	}

	var resultsView string
	if len(m.pingResults) > 0 {
		resultsView = subtitleStyle.Render("\nRecent Pings:")

		for i, result := range m.pingResults {
			if i >= 10 {
				break
			}

			timestamp := result.Timestamp.Format("15:04:05")
			var statusIcon, statusText string

			if result.Success {
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

			resultLine := fmt.Sprintf("  %s %s %s %s",
				dimTextStyle.Render(timestamp),
				statusIcon,
				statusText,
				duration,
			)

			resultsView += "\n" + resultLine
		}
	} else {
		resultsView = dimTextStyle.Render("\nNo ping results yet...")
	}

	var instructions string
	if m.isRunning {
		instructions = dimTextStyle.Render("\n\ns: Stop • Esc/q: Exit")
	} else {
		instructions = dimTextStyle.Render("\n\ns: Start • Esc/q: Exit")
	}

	content := fmt.Sprintf("%s\n\n%s\n\n%s%s%s",
		header,
		profileInfo,
		status,
		resultsView,
		instructions,
	)

	return panelStyle.Render(content)
}

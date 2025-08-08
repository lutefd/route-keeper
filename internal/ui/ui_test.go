package ui

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lutefd/route-keeper/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestNewMainModel(t *testing.T) {
	pm := models.NewProfilesManager()
	model := NewMainModel(pm)

	assert.NotNil(t, model)
	assert.Equal(t, MainMenuView, model.State)
	assert.Len(t, model.Inputs, 6)
	assert.Equal(t, "Profile name", model.Inputs[0].Placeholder)
}

func TestMainModel_Init(t *testing.T) {
	pm := models.NewProfilesManager()
	model := NewMainModel(pm)

	cmd := model.Init()
	assert.NotNil(t, cmd)
}

func TestMainModel_Update(t *testing.T) {
	t.Run("WindowSizeMsg", func(t *testing.T) {
		pm := models.NewProfilesManager()
		model := NewMainModel(pm)

		msg := tea.WindowSizeMsg{Width: 100, Height: 50}
		updatedModel, cmd := model.Update(msg)

		assert.NotNil(t, updatedModel)
		assert.Nil(t, cmd)
		assert.Equal(t, 100, updatedModel.(*MainModel).Width)
		assert.Equal(t, 50, updatedModel.(*MainModel).Height)
	})

	t.Run("tickMsg when running", func(t *testing.T) {
		pm := models.NewProfilesManager()
		model := NewMainModel(pm)
		model.IsRunning = true

		msg := tickMsg(time.Now())
		_, cmd := model.Update(msg)

		assert.NotNil(t, cmd)
	})
}

func TestMainModel_HandleKeyPress(t *testing.T) {
	t.Run("MainMenuView - up/down navigation", func(t *testing.T) {
		pm := models.NewProfilesManager()
		model := NewMainModel(pm)
		model.State = MainMenuView
		model.MenuIndex = 1

		_, _ = model.handleKeyPress(tea.KeyMsg{Type: tea.KeyUp})
		assert.Equal(t, 0, model.MenuIndex)

		_, _ = model.handleKeyPress(tea.KeyMsg{Type: tea.KeyDown})
		assert.Equal(t, 1, model.MenuIndex)

		model.MenuIndex = 2
		_, _ = model.handleKeyPress(tea.KeyMsg{Type: tea.KeyDown})
		assert.Equal(t, 2, model.MenuIndex)
	})

	t.Run("MainMenuView - enter key", func(t *testing.T) {
		pm := models.NewProfilesManager()
		model := NewMainModel(pm)
		model.State = MainMenuView
		model.MenuIndex = 0

		_, cmd := model.handleKeyPress(tea.KeyMsg{Type: tea.KeyEnter})
		assert.Nil(t, cmd)
		assert.Equal(t, ProfileListView, model.State)
	})
}

func TestMainModel_ProfileForm(t *testing.T) {
	pm := models.NewProfilesManager()
	model := NewMainModel(pm)
	model.State = CreateProfileView
	model.InputIndex = 0

	_, _ = model.handleKeyPress(tea.KeyMsg{Type: tea.KeyTab})
	assert.Equal(t, 1, model.InputIndex)

	_, _ = model.handleKeyPress(tea.KeyMsg{Type: tea.KeyShiftTab})
	assert.Equal(t, 0, model.InputIndex)
}

func TestMainModel_StartStopRunning(t *testing.T) {
	pm := models.NewProfilesManager()
	model := NewMainModel(pm)
	model.State = RunningView
	model.CurrentProfile = models.Profile{
		Name:    "Test Profile",
		BaseURL: "https://httpbin.org/get",
	}

	_, cmd := model.handleKeyPress(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	assert.True(t, model.IsRunning)
	assert.NotNil(t, cmd)

	_, _ = model.handleKeyPress(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	assert.False(t, model.IsRunning)
}

func TestMainModel_UpdateInputs(t *testing.T) {
	pm := models.NewProfilesManager()
	model := NewMainModel(pm)
	model.State = CreateProfileView

	model.Inputs[0].Focus()
	model.Inputs[0].SetValue("Test Profile")

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}}
	_, cmd := model.Update(msg)

	assert.NotNil(t, cmd)
}

func TestMainModel_CreateProfileFromInputs(t *testing.T) {
	pm := models.NewProfilesManager()
	model := NewMainModel(pm)

	model.Inputs[0].SetValue("Test Profile")
	model.Inputs[1].SetValue("https://api.example.com")
	model.Inputs[2].SetValue("/test")
	model.Inputs[3].SetValue("key1=value1,key2=value2")
	model.Inputs[4].SetValue("Authorization=Bearer token")
	model.Inputs[5].SetValue("5")

	profile := model.createProfileFromInputs()

	assert.Equal(t, "Test Profile", profile.Name)
	assert.Equal(t, "https://api.example.com", profile.BaseURL)
	assert.Equal(t, "/test", profile.Route)
	assert.Equal(t, 5, profile.Interval)
	assert.Equal(t, "value1", profile.Params["key1"])
	assert.Equal(t, "Bearer token", profile.Headers["Authorization"])
}

func TestMainModel_ResetInputs(t *testing.T) {
	pm := models.NewProfilesManager()
	model := NewMainModel(pm)

	for i := range model.Inputs {
		model.Inputs[i].SetValue("test")
	}

	model.resetInputs()

	for i, input := range model.Inputs {
		if i == 0 {
			assert.True(t, input.Focused())
		} else {
			assert.False(t, input.Focused())
		}
		assert.Empty(t, input.Value())
	}
}

func TestMainModel_PopulateInputsFromProfile(t *testing.T) {
	pm := models.NewProfilesManager()
	model := NewMainModel(pm)

	profile := models.Profile{
		Name:    "Test Profile",
		BaseURL: "https://api.example.com",
		Route:   "/test",
		Params: map[string]string{
			"key1": "value1",
		},
		Headers: map[string]string{
			"Authorization": "Bearer token",
		},
		Interval: 5,
	}

	model.populateInputsFromProfile(profile)

	assert.Equal(t, "Test Profile", model.Inputs[0].Value())
	assert.Equal(t, "https://api.example.com", model.Inputs[1].Value())
	assert.Equal(t, "/test", model.Inputs[2].Value())
	assert.Contains(t, model.Inputs[3].Value(), "key1=value1")
	assert.Contains(t, model.Inputs[4].Value(), "Authorization=Bearer token")
	assert.Equal(t, "5", model.Inputs[5].Value())
}

func TestMainModel_View(t *testing.T) {
	pm := models.NewProfilesManager()
	model := NewMainModel(pm)

	model.State = MainMenuView
	view := model.View()
	assert.Contains(t, view, "ROUTE KEEPER")
	assert.Contains(t, view, "Select Profile")

	model.State = ProfileListView
	view = model.View()
	assert.Contains(t, view, "SELECT PROFILE")

	model.State = CreateProfileView
	view = model.View()
	assert.Contains(t, view, "Create New Profile")

	model.State = EditProfileView
	view = model.View()
	assert.Contains(t, view, "Edit Profile")

	model.State = RunningView
	view = model.View()
	assert.Contains(t, view, "MONITORING")
}

package ui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lutefd/route-keeper/internal/models"
)

type ViewState int

const (
	MainMenuView ViewState = iota
	ProfileListView
	CreateProfileView
	EditProfileView
	RunningView
)

type tickMsg time.Time
type pingResultMsg models.PingResult

type MainModel struct {
	State           ViewState
	ProfilesManager *models.ProfilesManager
	PingService     *models.PingService

	MenuIndex    int
	ProfileIndex int

	EditingProfile models.Profile
	InputIndex     int
	Inputs         []textinput.Model
	IsEditing      bool

	CurrentProfile models.Profile
	IsRunning      bool
	Ticker         *time.Ticker
	PingResults    []models.PingResult

	Width  int
	Height int
}

func NewMainModel(pm *models.ProfilesManager) *MainModel {
	inputs := make([]textinput.Model, 6)

	inputs[0] = textinput.New()
	inputs[0].Placeholder = "Profile name"
	inputs[0].Focus()

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "https://api.example.com"

	inputs[2] = textinput.New()
	inputs[2].Placeholder = "/health"

	inputs[3] = textinput.New()
	inputs[3].Placeholder = "key1=value1,key2=value2"

	inputs[4] = textinput.New()
	inputs[4].Placeholder = "Authorization=Bearer token,Content-Type=application/json"

	inputs[5] = textinput.New()
	inputs[5].Placeholder = "5"

	return &MainModel{
		State:           MainMenuView,
		ProfilesManager: pm,
		PingService:     models.NewPingService(),
		Inputs:          inputs,
		PingResults:     []models.PingResult{},
	}
}

func (m *MainModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tickMsg:
		if m.IsRunning {
			return m, tea.Batch(
				m.doPing(),
				m.tick(),
			)
		}

	case pingResultMsg:
		m.PingResults = append([]models.PingResult{models.PingResult(msg)}, m.PingResults...)
		if len(m.PingResults) > 20 {
			m.PingResults = m.PingResults[:20]
		}
	}

	if m.State == CreateProfileView || m.State == EditProfileView {
		cmd := m.updateInputs(msg)
		return m, cmd
	}

	return m, nil
}

func (m *MainModel) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		if m.State == RunningView {
			return m.stopRunning(), nil
		}
		return m, tea.Quit

	case "esc":
		switch m.State {
		case ProfileListView, CreateProfileView, EditProfileView:
			m.State = MainMenuView
			m.MenuIndex = 0
		case RunningView:
			return m.stopRunning(), nil
		}
	}

	switch m.State {
	case MainMenuView:
		return m.handleMainMenuKeys(msg)
	case ProfileListView:
		return m.handleProfileListKeys(msg)
	case CreateProfileView, EditProfileView:
		return m.handleProfileFormKeys(msg)
	case RunningView:
		return m.handleRunningKeys(msg)
	}

	return m, nil
}

func (m *MainModel) handleMainMenuKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.MenuIndex > 0 {
			m.MenuIndex--
		}
	case "down", "j":
		if m.MenuIndex < 2 {
			m.MenuIndex++
		}
	case "enter":
		switch m.MenuIndex {
		case 0:
			m.State = ProfileListView
			m.ProfileIndex = 0
		case 1:
			m.State = CreateProfileView
			m.IsEditing = false
			m.resetInputs()
		case 2:
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *MainModel) handleProfileListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	profiles := m.ProfilesManager.GetProfiles()

	switch msg.String() {
	case "up", "k":
		if m.ProfileIndex > 0 {
			m.ProfileIndex--
		}
	case "down", "j":
		if m.ProfileIndex < len(profiles)-1 {
			m.ProfileIndex++
		}
	case "enter":
		if len(profiles) > 0 {
			m.CurrentProfile = profiles[m.ProfileIndex]
			m.State = RunningView
			return m.startRunning()
		}
	case "e":
		if len(profiles) > 0 {
			m.EditingProfile = profiles[m.ProfileIndex]
			m.populateInputsFromProfile(m.EditingProfile)
			m.State = EditProfileView
			m.IsEditing = true
		}
	case "d":
		if len(profiles) > 0 {
			profile := profiles[m.ProfileIndex]
			m.ProfilesManager.DeleteProfile(profile.Name)
			if m.ProfileIndex >= len(m.ProfilesManager.GetProfiles()) {
				m.ProfileIndex = len(m.ProfilesManager.GetProfiles()) - 1
			}
			if m.ProfileIndex < 0 {
				m.ProfileIndex = 0
			}
		}
	case "c":
		m.State = CreateProfileView
		m.IsEditing = false
		m.resetInputs()
	}
	return m, nil
}

func (m *MainModel) handleProfileFormKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "shift+tab":
		if m.InputIndex > 0 {
			m.Inputs[m.InputIndex].Blur()
			m.InputIndex--
			m.Inputs[m.InputIndex].Focus()
		}
	case "down", "tab":
		if m.InputIndex < len(m.Inputs)-1 {
			m.Inputs[m.InputIndex].Blur()
			m.InputIndex++
			m.Inputs[m.InputIndex].Focus()
		}
	case "enter":
		if m.InputIndex == len(m.Inputs)-1 {
			profile := m.createProfileFromInputs()
			if profile.Name != "" && profile.BaseURL != "" {
				m.ProfilesManager.AddProfile(profile)
				m.State = MainMenuView
				m.resetInputs()
			}
		} else {
			m.Inputs[m.InputIndex].Blur()
			m.InputIndex++
			if m.InputIndex >= len(m.Inputs) {
				m.InputIndex = len(m.Inputs) - 1
			}
			m.Inputs[m.InputIndex].Focus()
		}
	}
	return m, m.updateInputs(msg)
}

func (m *MainModel) handleRunningKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "s":
		if m.IsRunning {
			return m.stopRunning(), nil
		} else {
			return m.startRunning()
		}
	}
	return m, nil
}

func (m *MainModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.Inputs))
	for i := range m.Inputs {
		m.Inputs[i], cmds[i] = m.Inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m *MainModel) resetInputs() {
	for i := range m.Inputs {
		m.Inputs[i].SetValue("")
		m.Inputs[i].Blur()
	}
	m.Inputs[0].Focus()
}

func (m *MainModel) populateInputsFromProfile(profile models.Profile) {
	m.Inputs[0].SetValue(profile.Name)
	m.Inputs[1].SetValue(profile.BaseURL)
	m.Inputs[2].SetValue(profile.Route)

	var params []string
	for k, v := range profile.Params {
		params = append(params, fmt.Sprintf("%s=%s", k, v))
	}
	m.Inputs[3].SetValue(strings.Join(params, ","))

	var headers []string
	for k, v := range profile.Headers {
		headers = append(headers, fmt.Sprintf("%s=%s", k, v))
	}
	m.Inputs[4].SetValue(strings.Join(headers, ","))

	m.Inputs[5].SetValue(strconv.Itoa(profile.Interval))

	for i := range m.Inputs {
		m.Inputs[i].Blur()
	}
	m.Inputs[0].Focus()
}

func (m *MainModel) createProfileFromInputs() models.Profile {
	profile := models.Profile{
		Name:     m.Inputs[0].Value(),
		BaseURL:  m.Inputs[1].Value(),
		Route:    m.Inputs[2].Value(),
		Params:   make(map[string]string),
		Headers:  make(map[string]string),
		Interval: 5,
	}

	if paramsStr := m.Inputs[3].Value(); paramsStr != "" {
		for _, pair := range strings.Split(paramsStr, ",") {
			if kv := strings.SplitN(strings.TrimSpace(pair), "=", 2); len(kv) == 2 {
				profile.Params[kv[0]] = kv[1]
			}
		}
	}

	if headersStr := m.Inputs[4].Value(); headersStr != "" {
		for _, pair := range strings.Split(headersStr, ",") {
			if kv := strings.SplitN(strings.TrimSpace(pair), "=", 2); len(kv) == 2 {
				profile.Headers[kv[0]] = kv[1]
			}
		}
	}

	if intervalStr := m.Inputs[5].Value(); intervalStr != "" {
		if interval, err := strconv.Atoi(intervalStr); err == nil && interval > 0 {
			profile.Interval = interval
		}
	}

	return profile
}

func (m *MainModel) startRunning() (tea.Model, tea.Cmd) {
	m.IsRunning = true
	m.PingResults = []models.PingResult{}
	return m, tea.Batch(
		m.doPing(),
		m.tick(),
	)
}

func (m *MainModel) stopRunning() tea.Model {
	m.IsRunning = false
	if m.Ticker != nil {
		m.Ticker.Stop()
		m.Ticker = nil
	}
	m.State = MainMenuView
	return m
}

func (m *MainModel) tick() tea.Cmd {
	if !m.IsRunning {
		return nil
	}
	return tea.Tick(time.Duration(m.CurrentProfile.Interval)*time.Minute, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m *MainModel) doPing() tea.Cmd {
	return func() tea.Msg {
		result := m.PingService.Ping(m.CurrentProfile)
		return pingResultMsg(result)
	}
}

func (m *MainModel) View() string {
	switch m.State {
	case MainMenuView:
		return m.mainMenuView()
	case ProfileListView:
		return m.profileListView()
	case CreateProfileView:
		return m.profileFormView("Create New Profile")
	case EditProfileView:
		return m.profileFormView("Edit Profile")
	case RunningView:
		return m.runningView()
	}
	return "Unknown view"
}

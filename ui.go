package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ViewState int

const (
	MainMenuView ViewState = iota
	ProfileListView
	CreateProfileView
	EditProfileView
	RunningView
)

type MainModel struct {
	state           ViewState
	profilesManager *ProfilesManager
	pingService     *PingService

	menuIndex    int
	profileIndex int

	editingProfile Profile
	inputIndex     int
	inputs         []textinput.Model
	isEditing      bool

	currentProfile Profile
	isRunning      bool
	ticker         *time.Ticker
	pingResults    []PingResult

	width  int
	height int
}

type tickMsg time.Time
type pingResultMsg PingResult

func NewMainModel(pm *ProfilesManager) *MainModel {
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
		state:           MainMenuView,
		profilesManager: pm,
		pingService:     NewPingService(),
		inputs:          inputs,
		pingResults:     []PingResult{},
	}
}

func (m MainModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tickMsg:
		if m.isRunning {
			return m, tea.Batch(
				m.doPing(),
				m.tick(),
			)
		}

	case pingResultMsg:
		m.pingResults = append([]PingResult{PingResult(msg)}, m.pingResults...)
		if len(m.pingResults) > 20 {
			m.pingResults = m.pingResults[:20]
		}
	}

	if m.state == CreateProfileView || m.state == EditProfileView {
		cmd := m.updateInputs(msg)
		return m, cmd
	}

	return m, nil
}

func (m MainModel) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		if m.state == RunningView {
			return m.stopRunning(), nil
		}
		return m, tea.Quit

	case "esc":
		switch m.state {
		case ProfileListView, CreateProfileView, EditProfileView:
			m.state = MainMenuView
			m.menuIndex = 0
		case RunningView:
			return m.stopRunning(), nil
		}
	}

	switch m.state {
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

func (m MainModel) handleMainMenuKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.menuIndex > 0 {
			m.menuIndex--
		}
	case "down", "j":
		if m.menuIndex < 2 {
			m.menuIndex++
		}
	case "enter":
		switch m.menuIndex {
		case 0:
			m.state = ProfileListView
			m.profileIndex = 0
		case 1:
			m.state = CreateProfileView
			m.isEditing = false
			m.resetInputs()
		case 2:
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m MainModel) handleProfileListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	profiles := m.profilesManager.GetProfiles()

	switch msg.String() {
	case "up", "k":
		if m.profileIndex > 0 {
			m.profileIndex--
		}
	case "down", "j":
		if m.profileIndex < len(profiles)-1 {
			m.profileIndex++
		}
	case "enter":
		if len(profiles) > 0 {
			m.currentProfile = profiles[m.profileIndex]
			m.state = RunningView
			return m.startRunning()
		}
	case "e":
		if len(profiles) > 0 {
			m.editingProfile = profiles[m.profileIndex]
			m.populateInputsFromProfile(m.editingProfile)
			m.state = EditProfileView
			m.isEditing = true
		}
	case "d":
		if len(profiles) > 0 {
			profile := profiles[m.profileIndex]
			m.profilesManager.DeleteProfile(profile.Name)
			if m.profileIndex >= len(m.profilesManager.GetProfiles()) {
				m.profileIndex = len(m.profilesManager.GetProfiles()) - 1
			}
			if m.profileIndex < 0 {
				m.profileIndex = 0
			}
		}
	case "c":
		m.state = CreateProfileView
		m.isEditing = false
		m.resetInputs()
	}
	return m, nil
}

func (m MainModel) handleProfileFormKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "shift+tab":
		if m.inputIndex > 0 {
			m.inputs[m.inputIndex].Blur()
			m.inputIndex--
			m.inputs[m.inputIndex].Focus()
		}
	case "down", "tab":
		if m.inputIndex < len(m.inputs)-1 {
			m.inputs[m.inputIndex].Blur()
			m.inputIndex++
			m.inputs[m.inputIndex].Focus()
		}
	case "enter":
		if m.inputIndex == len(m.inputs)-1 {
			profile := m.createProfileFromInputs()
			if profile.Name != "" && profile.BaseURL != "" {
				m.profilesManager.AddProfile(profile)
				m.state = MainMenuView
				m.resetInputs()
			}
		} else {
			m.inputs[m.inputIndex].Blur()
			m.inputIndex++
			if m.inputIndex >= len(m.inputs) {
				m.inputIndex = len(m.inputs) - 1
			}
			m.inputs[m.inputIndex].Focus()
		}
	}
	return m, m.updateInputs(msg)
}

func (m MainModel) handleRunningKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "s":
		if m.isRunning {
			return m.stopRunning(), nil
		} else {
			return m.startRunning()
		}
	}
	return m, nil
}

func (m MainModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m MainModel) resetInputs() {
	for i := range m.inputs {
		m.inputs[i].SetValue("")
		m.inputs[i].Blur()
	}
	m.inputs[0].Focus()
}

func (m MainModel) populateInputsFromProfile(profile Profile) {
	m.inputs[0].SetValue(profile.Name)
	m.inputs[1].SetValue(profile.BaseURL)
	m.inputs[2].SetValue(profile.Route)

	var params []string
	for k, v := range profile.Params {
		params = append(params, fmt.Sprintf("%s=%s", k, v))
	}
	m.inputs[3].SetValue(strings.Join(params, ","))

	var headers []string
	for k, v := range profile.Headers {
		headers = append(headers, fmt.Sprintf("%s=%s", k, v))
	}
	m.inputs[4].SetValue(strings.Join(headers, ","))

	m.inputs[5].SetValue(strconv.Itoa(profile.Interval))

	for i := range m.inputs {
		m.inputs[i].Blur()
	}
	m.inputs[0].Focus()
}

func (m MainModel) createProfileFromInputs() Profile {
	profile := Profile{
		Name:     m.inputs[0].Value(),
		BaseURL:  m.inputs[1].Value(),
		Route:    m.inputs[2].Value(),
		Params:   make(map[string]string),
		Headers:  make(map[string]string),
		Interval: 5,
	}

	if paramsStr := m.inputs[3].Value(); paramsStr != "" {
		for _, pair := range strings.Split(paramsStr, ",") {
			if kv := strings.SplitN(strings.TrimSpace(pair), "=", 2); len(kv) == 2 {
				profile.Params[kv[0]] = kv[1]
			}
		}
	}

	if headersStr := m.inputs[4].Value(); headersStr != "" {
		for _, pair := range strings.Split(headersStr, ",") {
			if kv := strings.SplitN(strings.TrimSpace(pair), "=", 2); len(kv) == 2 {
				profile.Headers[kv[0]] = kv[1]
			}
		}
	}

	if intervalStr := m.inputs[5].Value(); intervalStr != "" {
		if interval, err := strconv.Atoi(intervalStr); err == nil && interval > 0 {
			profile.Interval = interval
		}
	}

	return profile
}

func (m MainModel) startRunning() (tea.Model, tea.Cmd) {
	m.isRunning = true
	m.pingResults = []PingResult{}

	return m, tea.Batch(
		m.doPing(),
		m.tick(),
	)
}

func (m MainModel) stopRunning() tea.Model {
	m.isRunning = false
	if m.ticker != nil {
		m.ticker.Stop()
		m.ticker = nil
	}
	m.state = MainMenuView
	return m
}

func (m MainModel) tick() tea.Cmd {
	if !m.isRunning {
		return nil
	}

	return tea.Tick(time.Duration(m.currentProfile.Interval)*time.Minute, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m MainModel) doPing() tea.Cmd {
	return func() tea.Msg {
		result := m.pingService.Ping(m.currentProfile)
		return pingResultMsg(result)
	}
}

func (m MainModel) View() string {
	switch m.state {
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

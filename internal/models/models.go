package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Profile struct {
	Name     string            `json:"name"`
	BaseURL  string            `json:"base_url"`
	Route    string            `json:"route"`
	Params   map[string]string `json:"params"`
	Headers  map[string]string `json:"headers"`
	Interval int               `json:"interval"`
}

func (p *Profile) GetFullURL() string {
	baseURL := strings.TrimSuffix(p.BaseURL, "/")
	u, err := url.Parse(baseURL)
	if err != nil {
		return fmt.Sprintf("%s/%s", baseURL, strings.TrimPrefix(p.Route, "/"))
	}

	if p.Route != "" {
		u.Path = filepath.Join(u.Path, strings.TrimPrefix(p.Route, "/"))
	}

	if len(p.Params) > 0 {
		query := u.Query()
		for k, v := range p.Params {
			query.Add(k, v)
		}
		u.RawQuery = query.Encode()
	}

	return u.String()
}

type PingResult struct {
	Timestamp  time.Time
	StatusCode int
	Success    bool
	Error      error
	Duration   time.Duration
}

type ProfilesManager struct {
	profiles []Profile
	filePath string
}

func NewProfilesManager() *ProfilesManager {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".route-keeper")
	os.MkdirAll(configDir, 0755)

	return &ProfilesManager{
		profiles: []Profile{},
		filePath: filepath.Join(configDir, "profiles.json"),
	}
}

func (pm *ProfilesManager) LoadProfiles() error {
	if _, err := os.Stat(pm.filePath); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(pm.filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &pm.profiles)
}

func (pm *ProfilesManager) SaveProfiles() error {
	data, err := json.MarshalIndent(pm.profiles, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(pm.filePath, data, 0644)
}

func (pm *ProfilesManager) AddProfile(profile Profile) error {
	for i, p := range pm.profiles {
		if p.Name == profile.Name {
			pm.profiles[i] = profile
			return pm.SaveProfiles()
		}
	}

	pm.profiles = append(pm.profiles, profile)
	return pm.SaveProfiles()
}

func (pm *ProfilesManager) GetProfiles() []Profile {
	return pm.profiles
}

func (pm *ProfilesManager) DeleteProfile(name string) error {
	for i, p := range pm.profiles {
		if p.Name == name {
			pm.profiles = append(pm.profiles[:i], pm.profiles[i+1:]...)
			return pm.SaveProfiles()
		}
	}
	return fmt.Errorf("profile not found")
}

type PingService struct{}

func NewPingService() *PingService {
	return &PingService{}
}

func (ps *PingService) Ping(profile Profile) PingResult {
	start := time.Now()
	result := PingResult{
		Timestamp: start,
		Success:   false,
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", profile.GetFullURL(), nil)
	if err != nil {
		result.Error = err
		result.Duration = time.Since(start)
		return result
	}

	for k, v := range profile.Headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		result.Error = err
		result.Duration = time.Since(start)
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	result.Success = resp.StatusCode >= 200 && resp.StatusCode < 300
	result.Duration = time.Since(start)

	return result
}

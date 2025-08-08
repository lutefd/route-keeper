package models

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProfile_GetFullURL(t *testing.T) {
	tests := []struct {
		name     string
		profile  Profile
		expected string
	}{
		{
			name: "basic URL",
			profile: Profile{
				BaseURL: "https://api.example.com",
				Route:   "/health",
			},
			expected: "https://api.example.com/health",
		},
		{
			name: "URL with params",
			profile: Profile{
				BaseURL: "https://api.example.com",
				Route:   "/data",
				Params: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
			expected: "https://api.example.com/data?key1=value1&key2=value2",
		},
		{
			name: "URL with existing query params",
			profile: Profile{
				BaseURL: "https://api.example.com?existing=param",
				Route:   "/endpoint",
				Params: map[string]string{
					"key1": "value1",
				},
			},
			expected: "https://api.example.com/endpoint?existing=param&key1=value1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.profile.GetFullURL()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestProfilesManager(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "route-keeper-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	pm := NewProfilesManager()
	assert.NotNil(t, pm)

	profile := Profile{
		Name:     "Test Profile",
		BaseURL:  "https://api.example.com",
		Route:    "/health",
		Interval: 5,
	}

	err = pm.AddProfile(profile)
	require.NoError(t, err)

	profiles := pm.GetProfiles()
	require.Len(t, profiles, 1)
	assert.Equal(t, profile.Name, profiles[0].Name)

	tempFile := filepath.Join(tempDir, "profiles.json")
	pm.filePath = tempFile

	err = pm.SaveProfiles()
	require.NoError(t, err)

	pm2 := NewProfilesManager()
	pm2.filePath = tempFile

	err = pm2.LoadProfiles()
	require.NoError(t, err)

	profiles = pm2.GetProfiles()
	require.Len(t, profiles, 1)
	assert.Equal(t, profile.Name, profiles[0].Name)

	err = pm2.DeleteProfile("Test Profile")
	require.NoError(t, err)
	assert.Len(t, pm2.GetProfiles(), 0)
	err = pm2.DeleteProfile("Non-existent")
	assert.Error(t, err)
}

func TestPingService(t *testing.T) {
	ps := NewPingService()
	require.NotNil(t, ps)

	profile := Profile{
		Name:    "Test Ping",
		BaseURL: "https://httpbin.org/get",
		Route:   "",
	}

	result := ps.Ping(profile)
	assert.False(t, result.Timestamp.IsZero())
}

func TestProfilesManager_EdgeCases(t *testing.T) {
	pm := NewProfilesManager()
	pm.filePath = "/invalid/path/profiles.json"

	err := pm.SaveProfiles()
	require.Error(t, err)

	err = pm.LoadProfiles()
	require.NoError(t, err)
	assert.Empty(t, pm.GetProfiles())
}

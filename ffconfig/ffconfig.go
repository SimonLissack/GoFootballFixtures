package ffconfig

import (
	"encoding/json"
)

// FFConfiguration provides the API key for getting fixtures as well as persistence options
type FFConfiguration struct {
	APIKey           string `json:"apiKey"`
	TeamsFile        string `json:"teamsFile"`
	PersistTeams     bool   `json:"persistTeams"`
	RebuildIfNoTeams bool   `json:"rebuildIfNoTeams"`
}

// LoadConfig Load the configuration from file
func LoadConfig(file []byte) (config *FFConfiguration, err error) {
	err = json.Unmarshal(file, &config)

	if err != nil {
		return nil, err
	}

	return config, nil
}

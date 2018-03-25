package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"

	c "github.com/simonlissack/footballfixtures/ffconfig"
	m "github.com/simonlissack/footballfixtures/model"
)

// TeamsCacheClient is an interface allowing reading/writing to the teams cache file
type TeamsCacheClient interface {
	// LoadTeams loads the list of teams from the cache
	LoadTeams() ([]m.Team, error)
	// SaveTeams Persists the list of teams to the cache
	SaveTeams([]m.Team) error
}

type localTeamsCacheClient struct {
	Config c.FFConfiguration
}

type localTeamsCache struct {
	Teams []m.Team `json:"teams"`
}

// NewLocalTeamsCache creates a cache client which stores the teams as JSON locally
func NewLocalTeamsCache(config c.FFConfiguration) TeamsCacheClient {
	return localTeamsCacheClient{Config: config}
}

func (ltc localTeamsCacheClient) LoadTeams() ([]m.Team, error) {
	if _, err := os.Stat(ltc.Config.TeamsFile); os.IsNotExist(err) {
		return nil, err
	}

	teamsFile, err := ioutil.ReadFile(ltc.Config.TeamsFile)

	if err != nil {
		return nil, err
	}

	var teamsCache localTeamsCache
	err = json.Unmarshal(teamsFile, &teamsCache)

	return teamsCache.Teams, err
}

func (ltc localTeamsCacheClient) SaveTeams(teams []m.Team) error {
	teamsCache := localTeamsCache{Teams: teams}
	j, err := json.Marshal(teamsCache)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(ltc.Config.TeamsFile, j, 0644)

	if err != nil {
		return err
	}

	return nil
}

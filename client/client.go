package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/simonlissack/footballfixtures/ffconfig"
	m "github.com/simonlissack/footballfixtures/model"
)

type competition struct {
	ID int `json:"id"`
}

type fdoCompetitionResponse struct {
	Count int      `json:"count"`
	Teams []m.Team `json:"teams"`
}

type fdoFixtureResponse struct {
	Count    int         `json:"count"`
	Fixtures []m.Fixture `json:"fixtures"`
}

// FootballClient is the interface for getting teams and fixtures
type FootballClient interface {
	// GetTeams gets all the available teams
	GetTeams() []m.Team
	// GetFixtures finds the fixtures in the next 30 days for a team, based on the teamID
	GetFixtures(teamID int) []m.Fixture
}

// FootballDataOrgClient is a client for football-data.org rest API
type footballDataOrgClient struct {
	config ffconfig.FFConfiguration
	teams  []m.Team
}

const (
	xAuthToken       = "X-Auth-Token"
	xResponseControl = "X-Response-Control"
	api              = "https://api.football-data.org"
	apiVersion       = "v1"
)

// End points for football-data.org
var fdoEndPoints = map[string]string{
	"Competitions": "competitions",
	"Teams":        "competitions/{cid}/teams",
	"Fixtures":     "teams/{tid}/fixtures?timeFrame=n30",
}

// NewFootballDataOrgClient creates a new FootballClient which uses football-data.org API
func NewFootballDataOrgClient(config ffconfig.FFConfiguration) FootballClient {
	return footballDataOrgClient{config: config}
}

func (fbClient footballDataOrgClient) GetTeams() (teams []m.Team) {
	competitions := fbClient.getCompetitions()
	teams = make([]m.Team, 0)
	for _, competition := range competitions {
		compTeams := fbClient.getTeamsInCompetition(competition.ID)
		teams = append(teams, compTeams...)
	}

	return
}

func (fbClient footballDataOrgClient) GetFixtures(teamID int) []m.Fixture {
	var fixtureResponse fdoFixtureResponse
	values := map[string]string{"tid": strconv.Itoa(teamID)}
	fbClient.makeMinifiedRequest(fdoEndPoints["Fixtures"], values, &fixtureResponse)

	return fixtureResponse.Fixtures
}

func (fbClient footballDataOrgClient) getCompetitions() []competition {
	competitions := make([]competition, 0)
	fbClient.makeMinifiedRequest(fdoEndPoints["Competitions"], nil, &competitions)

	return competitions
}

func (fbClient footballDataOrgClient) getTeamsInCompetition(competitionID int) []m.Team {
	var competitionResponse fdoCompetitionResponse
	values := map[string]string{"cid": strconv.Itoa(competitionID)}
	fbClient.makeMinifiedRequest(fdoEndPoints["Teams"], values, &competitionResponse)

	return competitionResponse.Teams
}

func (fbClient footballDataOrgClient) makeMinifiedRequest(endPoint string, values map[string]string, unmarshalTo interface{}) {
	fbClient.makeRequest(endPoint, "minified", values, unmarshalTo)
}

func (fbClient footballDataOrgClient) makeRequest(endPoint string, responseControl string, values map[string]string, unmarshalTo interface{}) {
	resp, _ := fbClient.sendRequest(endPoint, responseControl, values)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, unmarshalTo)
}

func (fbClient footballDataOrgClient) sendRequest(endPoint string, responseControl string, values map[string]string) (*http.Response, error) {
	client := &http.Client{}
	address := buildRequestURL(endPoint, values)
	req, _ := http.NewRequest("GET", address, nil)
	req.Header.Add(xAuthToken, fbClient.config.APIKey)
	req.Header.Add(xResponseControl, responseControl)

	return client.Do(req)
}

func buildRequestURL(endpoint string, values map[string]string) string {
	formattedEndPoint := endpoint
	for key, value := range values {
		token := fmt.Sprintf("{%s}", key)
		formattedEndPoint = strings.Replace(formattedEndPoint, token, value, 1)
	}

	return strings.Join([]string{api, apiVersion, formattedEndPoint}, "/")
}

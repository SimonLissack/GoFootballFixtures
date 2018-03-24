package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	t "time"

	"github.com/simonlissack/footballfixtures/ffconfig"
)

// Team represents a football team
type Team struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
}

// Fixture represents a scheduled match between two football teams
type Fixture struct {
	Date         t.Time `json:"date"`
	Status       string `json:"status"`
	HomeTeamName string `json:"homeTeamName"`
	HomeTeamID   int    `json:"homeTeamId"`
	AwayTeamName string `json:"awayTeamName"`
	AwayTeamID   int    `json:"awayTeamId"`
}

type competition struct {
	ID int `json:"id"`
}

type fdoCompetitionResponse struct {
	Count int    `json:"count"`
	Teams []Team `json:"teams"`
}

type fdoFixtureResponse struct {
	Count    int       `json:"count"`
	Fixtures []Fixture `json:"fixtures"`
}

// FootballClient is the interface for getting teams and fixtures
type FootballClient interface {
	// GetTeams gets all the available teams
	GetTeams() []Team
	// GetFixtures finds the fixtures in the next 30 days for a team, based on the teamID
	GetFixtures(teamID int) []Fixture
}

// FootballDataOrgClient is a client for football-data.org rest API
type footballDataOrgClient struct {
	config ffconfig.FFConfiguration
	teams  []Team
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

func (fbClient footballDataOrgClient) GetTeams() (teams []Team) {
	competitions := fbClient.getCompetitions()
	teams = make([]Team, 0)
	for _, competition := range competitions {
		compTeams := fbClient.getTeamsInCompetition(competition.ID)
		teams = append(teams, compTeams...)
	}

	return
}

func (fbClient footballDataOrgClient) GetFixtures(teamID int) []Fixture {
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

func (fbClient footballDataOrgClient) getTeamsInCompetition(competitionID int) []Team {
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

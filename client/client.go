package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/simonlissack/footballfixtures/storage"

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
	GetTeams() ([]m.Team, error)
	// GetFixtures finds the fixtures in the next 30 days for a team, based on the teamID
	GetFixtures(teamID int) []m.Fixture
	// CanMakeRequest checks if the client can make a request, returning an error if the user if being throttled
	CanMakeRequest() error
}

// FootballDataOrgClient is a client for football-data.org rest API
type footballDataOrgClient struct {
	config              ffconfig.FFConfiguration
	cache               storage.TeamsCacheClient
	teams               []m.Team
	requestsAvailable   int
	requestCounterReset int
	lastRequest         time.Time
}

const (
	// Request constants
	xAuthToken       = "X-Auth-Token"
	xResponseControl = "X-Response-Control"
	api              = "https://api.football-data.org"
	apiVersion       = "v1"
	// Response constants
	xRequestCounterReset = "X-Requestcounter-Reset"
	xRequestsAvailable   = "X-Requests-Available"
	date                 = "Date"
)

// End points for football-data.org
var fdoEndPoints = map[string]string{
	"Competitions": "competitions",
	"Teams":        "competitions/{cid}/teams",
	"Fixtures":     "teams/{tid}/fixtures?timeFrame=n30",
}

// NewFootballDataOrgClient creates a new FootballClient which uses football-data.org API
func NewFootballDataOrgClient(config ffconfig.FFConfiguration, cache storage.TeamsCacheClient) FootballClient {
	return footballDataOrgClient{
		config:              config,
		cache:               cache,
		lastRequest:         time.Time{},
		requestCounterReset: 0,
		requestsAvailable:   0,
	}
}

func (fbClient footballDataOrgClient) GetTeams() (teams []m.Team, err error) {
	if fbClient.config.PersistTeams && fbClient.cache != nil {
		teams, err = fbClient.cache.LoadTeams()

		if len(teams) > 0 || !fbClient.config.RebuildIfNoTeams {
			return
		}
	}

	competitions := fbClient.getCompetitions()
	teams = make([]m.Team, 0)
	for _, competition := range competitions {
		compTeams := fbClient.getTeamsInCompetition(competition.ID)
		teams = append(teams, compTeams...)
	}

	if fbClient.cache != nil && fbClient.config.PersistTeams {
		err = fbClient.cache.SaveTeams(teams)
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

func (fbClient footballDataOrgClient) makeMinifiedRequest(endPoint string, values map[string]string, unmarshalTo interface{}) error {
	return fbClient.makeRequest(endPoint, "minified", values, unmarshalTo)
}

func (fbClient footballDataOrgClient) makeRequest(endPoint string, responseControl string, values map[string]string, unmarshalTo interface{}) error {
	err := fbClient.CanMakeRequest()
	if err != nil {
		return err
	}

	resp, err := fbClient.sendRequest(endPoint, responseControl, values)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, unmarshalTo)

	if err != nil {
		return err
	}

	fbClient.updateRequestRate(resp.Header)

	return nil
}

func (fbClient footballDataOrgClient) sendRequest(endPoint string, responseControl string, values map[string]string) (*http.Response, error) {
	client := &http.Client{}
	address := buildRequestURL(endPoint, values)
	req, _ := http.NewRequest("GET", address, nil)
	req.Header.Add(xAuthToken, fbClient.config.APIKey)
	req.Header.Add(xResponseControl, responseControl)

	return client.Do(req)
}

func (fbClient footballDataOrgClient) CanMakeRequest() error {
	remainingSecs := requestTimeout(fbClient.lastRequest, fbClient.requestCounterReset)
	if fbClient.requestsAvailable == 0 && remainingSecs > 0 {
		return fmt.Errorf("Cannot make request, try again in %f seconds", remainingSecs)
	}
	return nil
}

func requestTimeout(lastReqTime time.Time, reqReset int) float64 {
	resetTime := lastReqTime.Add(time.Second * time.Duration(reqReset))
	return -time.Since(resetTime).Seconds()
}

func (fbClient footballDataOrgClient) updateRequestRate(headers map[string][]string) {
	fbClient.lastRequest = time.Now()
	fbClient.requestCounterReset, _ = strconv.Atoi(headers[xRequestCounterReset][0])
	fbClient.requestsAvailable, _ = strconv.Atoi(headers[xRequestsAvailable][0])
}

func buildRequestURL(endpoint string, values map[string]string) string {
	formattedEndPoint := endpoint
	for key, value := range values {
		token := fmt.Sprintf("{%s}", key)
		formattedEndPoint = strings.Replace(formattedEndPoint, token, value, 1)
	}

	return strings.Join([]string{api, apiVersion, formattedEndPoint}, "/")
}

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/simonlissack/footballfixtures/client"
	"github.com/simonlissack/footballfixtures/ffconfig"
)

var (
	configPath, teamName string
	teamID               int
)

func init() {
	flag.StringVar(&configPath, "config", "", "Path to the configuration file")
	flag.StringVar(&teamName, "team", "", "Name of team to lookup")
	flag.IntVar(&teamID, "team-id", -1, "ID of team to lookup")
	flag.Parse()

	if configPath == "" {
		log.Fatal("No configuration file specified")
	}

	if teamName == "" {
		log.Fatal("No team specified")
	}
}

func main() {
	// Load config
	configFile, err := ioutil.ReadFile(configPath)
	logFatal(err)
	config, err := ffconfig.LoadConfig(configFile)
	logFatal(err)

	client := client.NewFootballDataOrgClient(*config)

	if teamID == -1 {
		teams := client.GetTeams()
		teamID, err = lookupTeam(teams, teamName)
		logFatal(err)
	}

	// Get fixtures for the team
	fixtures := client.GetFixtures(teamID)

	printHomeFixtures(fixtures, teamID)
}

func lookupTeam(teams []client.Team, teamQuery string) (int, error) {
	for _, t := range teams {
		if t.Name == teamQuery || t.ShortName == teamQuery {
			return t.ID, nil
		}
	}
	err := fmt.Errorf("Could not find team '%s'", teamQuery)

	return -1, err
}

func printHomeFixtures(fixtures []client.Fixture, teamID int) {
	for _, f := range fixtures {

		if f.AwayTeamID == teamID {
			continue
		}

		fmt.Println(f.Date)
		fmt.Printf("%s vs %s", f.HomeTeamName, f.AwayTeamName)
		fmt.Println()
		fmt.Println()
	}
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

# FootballFixtures

Football fixtures is a Go application designed to lookup upcoming fixtures for a football team.

## Configuration

The application has a configuration struct called FFConfiguration this contains all the configuration for reading footbal data via the FootballDataOrg client.

The configuration file is a json file with the following format:

```
{
    "apiKey" : "<API-KEY>",
    "teamsFile" : "teams.json",
    "persistTeams" : true,
    "rebuildIfNoTeams": true
}
```

| Item             | Type   | Description                                                                                                        |
|------------------|--------|--------------------------------------------------------------------------------------------------------------------|
| apiKey           | string | The API key for football-data.org                                                                                  |
| teamsFile        | string | The location of the file containing all the teams                                                                  |
| persistTeams     | bool   | Whether the teams list is persisted or loaded from the API each time. It is recommended that this is set to `true` |
| rebuildIfNoTeams | bool   | Whether the teams list should be recreated, if the cache is loaded but has no results                              |

Note that the teamsFile will only be read from and written to if persistTeams is set to `true`

## Flags

If this is being run via main.go then the following flags must be specified

| Flag    | Description                                                                                        |
|---------|----------------------------------------------------------------------------------------------------|
| config  | the path to the configuration file for FFConfiguration                                             |
| team    | The team whose fixtures you want to lookup                                                         |
| team-id | The ID of the team you want to lookup (this will skip the step of loading teams from server/cache) |

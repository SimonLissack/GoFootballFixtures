package model

import (
	t "time"
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

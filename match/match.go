package match

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/robertobouses/online-football-tycoon/team"
)

type Match struct {
	HomeMatchStrategy Strategy
	AwayMatchStrategy Strategy
}

type Strategy struct {
	StrategyTeam         team.Team
	Formation            string
	PlayingStyle         string
	GameTempo            string
	PassingStyle         string
	DefensivePositioning string
	BuildUpPlay          string
	AttackFocus          string
	KeyPlayerUsage       string
}

type Result struct {
	HomeStats TeamStats
	AwayStats TeamStats
}

type TeamStats struct {
	BallPossession int
	ScoringChances int
	Goals          int
}

type MatchEventStats struct {
	HomeEvents       []EventResult
	AwayEvents       []EventResult
	HomeScoreChances int
	AwayScoreChances int
	HomeGoals        int
	AwayGoals        int
}

type MatchEvent struct {
	MatchID   uuid.UUID
	MatchInfo MatchEventStats
}

func (a AppService) PlayMatch(matchID uuid.UUID) (Result, error) {
	m, err := a.repo.GetMatchById(matchID)
	if err != nil {
		return Result{}, fmt.Errorf("error retrieving match: %w", err)
	}

	result, matchEventStats, err := m.Play()
	if err != nil {
		return Result{}, fmt.Errorf("error playing match: %w", err)
	}

	matchDate := time.Now()
	homeTeamId := m.HomeMatchStrategy.StrategyTeam.Id
	awayTeamId := m.AwayMatchStrategy.StrategyTeam.Id

	log.Println("homeTeamId, awayTeamId", homeTeamId, awayTeamId)

	err = a.repo.PostMatch(homeTeamId, awayTeamId, matchDate, result.HomeStats.Goals, result.AwayStats.Goals)
	if err != nil {
		return Result{}, fmt.Errorf("error PostMatch: %w", err)
	}

	matchEvent := MatchEvent{
		MatchID:   matchID,
		MatchInfo: matchEventStats,
	}

	err = a.repo.PostMatchEvent(matchEvent)
	if err != nil {
		return Result{}, err
	}

	return result, nil
}

func (m Match) Play() (Result, MatchEventStats, error) {
	lineup := m.HomeMatchStrategy.StrategyTeam.Players

	rivalLineup := m.AwayMatchStrategy.StrategyTeam.Players

	homeTeam := m.HomeMatchStrategy.StrategyTeam
	awayTeam := m.AwayMatchStrategy.StrategyTeam

	log.Println("rivalLineup", rivalLineup)

	numberOfMatchEvents, err := CalculateNumberOfMatchEvents(m.HomeMatchStrategy.GameTempo, m.AwayMatchStrategy.GameTempo)
	if err != nil {
		log.Println("error on numberOfMatchEvents", err)
		return Result{}, MatchEventStats{}, err
	}
	log.Println("numberOfMatchEvents", numberOfMatchEvents)

	numberOfLineupEvents, numberOfRivalEvents, err := DistributeMatchEvents(m.HomeMatchStrategy.StrategyTeam, m.AwayMatchStrategy.StrategyTeam, numberOfMatchEvents)
	if err != nil {
		log.Println("error al distribuir numberOfMatchEvents", err)
		return Result{}, MatchEventStats{}, err
	}
	log.Println("numberOfLineupEvents, numberOfRivalEvents", numberOfLineupEvents, numberOfRivalEvents)

	matchEventStats := GenerateEvents(homeTeam, awayTeam, numberOfLineupEvents, numberOfRivalEvents)
	breakMatch := EventResult{Minute: 45, Event: "Descanso"}
	endMatch := EventResult{Minute: 90, Event: "Final del Partido"}
	allEvents := append(matchEventStats.HomeEvents, matchEventStats.AwayEvents...)
	allEvents = append(allEvents, breakMatch, endMatch)
	sort.Slice(allEvents, func(i, j int) bool {
		return allEvents[i].Minute < allEvents[j].Minute
	})

	var totalHomeTechnique, totalHomeMental, totalHomePhysique int
	for _, player := range lineup {
		totalHomeTechnique += totalHomeTechnique + player.Technique
		totalHomeMental += totalHomeMental + player.Mental
		totalHomePhysique += totalHomePhysique + player.Physique
	}

	var totalAwayTechnique, totalAwayMental, totalAwayPhysique int
	for _, player := range rivalLineup {
		totalAwayTechnique += totalAwayTechnique + player.Technique
		totalAwayMental += totalAwayMental + player.Mental
		totalAwayPhysique += totalAwayPhysique + player.Physique
	}

	strategy := m.HomeMatchStrategy

	resultOfStrategy, err := CalculateResultOfStrategy(lineup, strategy.Formation, strategy.PlayingStyle, strategy.GameTempo, strategy.PassingStyle, strategy.DefensivePositioning, strategy.BuildUpPlay, strategy.AttackFocus, strategy.KeyPlayerUsage)
	if err != nil {

		return Result{}, MatchEventStats{}, fmt.Errorf("error in calculating the result of the strategy: %w", err)
	}

	totalHomePhysique = totalHomePhysique + int(resultOfStrategy["teamPhysique"])

	lineupTotalQuality, rivalTotalQuality, allQuality, err := CalculateTotalQuality(totalHomeTechnique, totalHomeMental, totalHomePhysique, totalAwayTechnique, totalAwayMental, totalAwayPhysique)
	if err != nil {
		log.Println("Error calculating total quality:", err)
		return Result{}, MatchEventStats{}, err
	}
	log.Printf("Total Quality: player %d, rival %d, total quality %d\n", lineupTotalQuality, rivalTotalQuality, allQuality)
	lineupPercentagePossession, rivalPercentagePossession, err := CalculateBallPossession(totalHomeTechnique, totalHomeMental, lineupTotalQuality, rivalTotalQuality, allQuality, resultOfStrategy["teamPossession"])
	if err != nil {
		log.Println("Error CalculateBallPossession:", err)
		return Result{}, MatchEventStats{}, err
	}

	result := Result{
		HomeStats: TeamStats{
			BallPossession: lineupPercentagePossession,
			ScoringChances: matchEventStats.HomeScoreChances,
			Goals:          matchEventStats.HomeGoals,
		},
		AwayStats: TeamStats{
			BallPossession: rivalPercentagePossession,
			ScoringChances: matchEventStats.AwayScoreChances,
			Goals:          matchEventStats.AwayGoals,
		},
	}

	return result, matchEventStats, nil
}

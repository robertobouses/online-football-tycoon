package match

import (
	"database/sql"

	_ "embed"
)

//go:embed sql/get_matches.sql
var getMatchesQuery string

//go:embed sql/get_match_teams.sql
var getMatchTeamsQuery string

//go:embed sql/get_match_strategies.sql
var getMatchStrategiesQuery string

//go:embed sql/get_match_players.sql
var getMatchPlayersQuery string

//go:embed sql/post_match.sql
var postMatchQuery string

//go:embed sql/post_match_events.sql
var postMatchEventsQuery string

//go:embed sql/get_pending_matches.sql
var getPendingMatches string

func NewRepository(db *sql.DB) (*Repository, error) {
	getMatchesStmt, err := db.Prepare(getMatchesQuery)
	if err != nil {
		return nil, err
	}

	getMatchTeamsStmt, err := db.Prepare(getMatchTeamsQuery)
	if err != nil {
		return nil, err
	}

	getMatchStrategiesStmt, err := db.Prepare(getMatchStrategiesQuery)
	if err != nil {
		return nil, err
	}

	getMatchPlayersStmt, err := db.Prepare(getMatchPlayersQuery)
	if err != nil {
		return nil, err
	}

	postMatchStmt, err := db.Prepare(postMatchQuery)
	if err != nil {
		return nil, err
	}

	postMatchEventsStmt, err := db.Prepare(postMatchEventsQuery)
	if err != nil {
		return nil, err
	}
	getPendingMatchesStmt, err := db.Prepare(getPendingMatches)
	if err != nil {
		return nil, err
	}

	return &Repository{
		db:                 db,
		getMatches:         getMatchesStmt,
		getMatchTeams:      getMatchTeamsStmt,
		getMatchStrategies: getMatchStrategiesStmt,
		getMatchPlayers:    getMatchPlayersStmt,
		postMatch:          postMatchStmt,
		postMatchEvents:    postMatchEventsStmt,
		getPendingMatches:  getPendingMatchesStmt,
	}, nil
}

type Repository struct {
	db                 *sql.DB
	getMatches         *sql.Stmt
	getMatchTeams      *sql.Stmt
	getMatchStrategies *sql.Stmt
	getMatchPlayers    *sql.Stmt
	postMatch          *sql.Stmt
	postMatchEvents    *sql.Stmt
	getPendingMatches  *sql.Stmt
}

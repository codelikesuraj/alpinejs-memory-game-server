package leaderboard

import (
	"context"
	"database/sql"
	"encore.dev/storage/sqldb"
	"errors"
)

var highscoredb = sqldb.NewDatabase("leaderboards", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})

func fetchHighScores(ctx context.Context) ([]HighScore, error) {
	rows, err := highscoredb.Query(ctx, `
		SELECT duration, username
		FROM leaderboards
		ORDER BY duration ASC, username ASC
		LIMIT 15
	`)
	if err != nil {
		return nil, err
	}

	var highscores []HighScore
	for rows.Next() {
		var lb HighScore
		if err := rows.Scan(&lb.Duration, &lb.Username); err != nil {
			return nil, err
		}
		highscores = append(highscores, lb)
	}
	return highscores, nil
}

func addHighscore(ctx context.Context, highscore HighScore) (string, error) {
	var existingDuration int
	duration, username := highscore.Duration, highscore.Username

	err := highscoredb.QueryRow(ctx, `
		SELECT duration
		FROM leaderboards
		WHERE username = $1
	`, username).Scan(&existingDuration)

	if err != nil {
		// If no rows, insert
		if errors.Is(err, sql.ErrNoRows) {
			_, err = highscoredb.Exec(ctx, `
				INSERT INTO leaderboards (duration, username)
				VALUES ($1, $2)
			`, duration, username)
		}
		if err != nil {
			return "highscore added", nil
		}

		// any other error
		return "", err
	}

	// username exists, compare durations
	if duration < existingDuration {
		_, err = highscoredb.Exec(ctx, `
			UPDATE leaderboards
			SET duration = $1,
			    updated_at = CURRENT_TIMESTAMP
			WHERE username = $2
		`, highscore.Duration, highscore.Username)
		if err != nil {
			return "highscore updated", nil
		}
	}

	// no updated needed (new duration is worse)
	return "no updated needed", nil
}

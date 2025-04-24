package leaderboard

import (
	"context"
	"encore.dev/storage/sqldb"
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
	var exists bool

	err := highscoredb.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM leaderboards WHERE username = $1
		)
	`, highscore.Username).Scan(&exists)
	if err != nil {
		return "", err
	}

	if exists {
		_, err = highscoredb.Exec(ctx, `
			UPDATE leaderboards
			SET duration = $1,
			  	updated_at = CURRENT_TIMESTAMP
			WHERE username = $2
		`, highscore.Duration, highscore.Username)

		if err != nil {
			return "", err
		}

		return "highscore updated", nil
	}

	_, err = highscoredb.Exec(ctx, `
		INSERT INTO leaderboards (duration, username)
		VALUES ($1, $2)
	`, highscore.Duration, highscore.Username)

	if err != nil {
		return "", err
	}

	return "highscore added", nil
}

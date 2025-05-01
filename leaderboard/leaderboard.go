package leaderboard

import (
	"context"
	"encore.app/email"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
	"strings"
)

type HighScore struct {
	Duration int    `json:"duration"`
	Username string `json:"username"`
}

type Resp struct {
	HighScores []HighScore `json:"highscores"`
}

//encore:api public method=POST path=/leaderboard
func Add(ctx context.Context, params *HighScore) error {
	duration, username := params.Duration, strings.Trim(params.Username, " ")

	if duration < 1 {
		return &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "duration cannot be less than a second",
		}
	}

	if l := len(username); l < 3 || l > 16 {
		return &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "invalid username - numbers of characters should be between 3 & 16",
		}
	}

	info, err := addHighscore(ctx, HighScore{duration, username})
	if err != nil {
		rlog.Error("error adding high score",
			"err", err,
		)

		return &errs.Error{
			Code:    errs.Internal,
			Message: "Oops, something went wrong",
		}
	}

	rlog.Info(
		info,
		"username", username,
		"duration", duration,
	)

	if info == "highscore added" {
		email.HighScores.Publish(ctx, &email.NewHighScoreEvent{
			HighScore: struct {
				Duration int    `json:"duration"`
				Username string `json:"username"`
			}{Duration: duration, Username: username},
		})
	}

	return nil
}

//encore:api public method=GET path=/leaderboard
func Get(ctx context.Context) (Resp, error) {
	highscores, err := fetchHighScores(ctx)
	if err != nil || len(highscores) == 0 {
		return Resp{[]HighScore{}}, err
	}

	rlog.Info("fetched highscores")

	return Resp{highscores}, nil
}

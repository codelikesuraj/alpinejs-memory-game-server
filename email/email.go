package email

import (
	"context"
	"encore.dev/pubsub"
	"fmt"
	"github.com/mailersend/mailersend-go"
	"time"
)

type NewHighScoreEvent struct {
	HighScore struct {
		Duration int    `json:"duration"`
		Username string `json:"username"`
	}
}

var secrets struct {
	MailerSendAPIKey string
}

var HighScores = pubsub.NewTopic[*NewHighScoreEvent](
	"new-username",
	pubsub.TopicConfig{
		DeliveryGuarantee: pubsub.ExactlyOnce,
	})

var _ = pubsub.NewSubscription(
	HighScores,
	"new-username",
	pubsub.SubscriptionConfig[*NewHighScoreEvent]{
		Handler: func(ctx context.Context, event *NewHighScoreEvent) error {
			ms := mailersend.NewMailersend(secrets.MailerSendAPIKey)
			ctx, cancel := context.WithTimeout(ctx, time.Second*5)
			defer cancel()

			highscore := event.HighScore

			message := ms.Email.NewMessage()
			message.SetFrom(mailersend.From{
				Name:  "Abdulbaki Suraj",
				Email: "notifications@abdulbakisuraj.com",
			})
			message.SetRecipients([]mailersend.Recipient{
				{
					Name:  "Abdulbaki Suraj",
					Email: "surajabdulbaki19@gmail.com",
				},
			})
			message.SetSubject("New HighScore - AlpineJS Memory Game")
			message.SetHTML(fmt.Sprintf("<h1>Hey Suraj,<h1><p>A new user, '%s', just completed the game in %d seconds.<p>", highscore.Username, highscore.Duration))

			_, err := ms.Email.Send(ctx, message)
			if err != nil {
				return err
			}

			return nil
		},
	},
)

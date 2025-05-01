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
			message.SetFrom(mailersend.From{Email: "notifications@abdulbakisuraj.com"})
			message.SetSubject("New highscore added")
			message.SetText(fmt.Sprintf("New user '%s' just completed a game in %d seconds", highscore.Username, highscore.Duration))
			message.SetRecipients([]mailersend.Recipient{{Email: "surajabdulbaki19@gmail.com"}})

			_, err := ms.Email.Send(ctx, message)
			if err != nil {
				return err
			}

			return nil
		},
	},
)

package cmd

import (
	"fmt"
	"os"

	slack "github.com/monochromegane/slack-incoming-webhooks"
	"github.com/udzura/collector/collectorlib"
)

var incomingWebhookUrl string

var bodyFormat = "*Domain:* %s\n" +
	"*Change detail:*\n" +
	"```\n" +
	"%s\n" +
	"```"

func NotifyToSlack(domain string, diff *collectorlib.Diff) {
	if incomingWebhookUrl == "" {
		return
	}

	payload := &slack.Payload{
		Username: "From collector",
		Attachments: []*slack.Attachment{
			{
				Pretext:    "DNS Records are changed:",
				Text:       fmt.Sprintf(bodyFormat, domain, diff.ToString()),
				Color:      "#ff6600",
				MarkdownIn: []string{"text"},
			},
		},
	}
	if emoji := os.Getenv("SLACK_ICON_EMOJI"); emoji != "" {
		payload.IconEmoji = emoji
	}
	if url := os.Getenv("SLACK_ICON_URL"); url != "" {
		payload.IconURL = url
	}

	cli := slack.Client{
		WebhookURL: incomingWebhookUrl,
	}
	err := cli.Post(payload)
	if err != nil {
		logger.Warnf("Notification error: %s, but ignored.", err.Error())
	}
}

func init() {
	incomingWebhookUrl = os.Getenv("SLACK_URL")
}

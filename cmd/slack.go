package cmd

import (
	"fmt"
	"os"
	"sort"

	slack "github.com/monochromegane/slack-incoming-webhooks"
)

var incomingWebhookUrl string

var bodyFormat = "*Domain:* %s\n" +
	"*Change detail:*\n" +
	"```\n" +
	"%s\n" +
	"```"

func NotifyToSlack(domain string, ipsBefore, ipsAfter []string) {
	if incomingWebhookUrl == "" {
		return
	}

	payload := &slack.Payload{
		Username: "From collector",
		Attachments: []*slack.Attachment{
			{
				Pretext:    "DNS Records are changed:",
				Text:       fmt.Sprintf(bodyFormat, domain, toDiff(ipsBefore, ipsAfter)),
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

func toDiff(ipsBefore, ipsAfter []string) string {
	sort.Strings(ipsBefore)
	sort.Strings(ipsAfter)
	var added, deleted, existsBoth []string
	deleted = append(deleted, ipsBefore...)
	added = append(added, ipsAfter...)
	for _, i1 := range ipsAfter {
		for idx, i2 := range deleted {
			if i1 == i2 {
				deleted = append(deleted[:idx], deleted[idx+1:]...)
				existsBoth = append(existsBoth, i1)
			}
		}
	}

	for _, i1 := range existsBoth {
		for idx, i2 := range added {
			if i1 == i2 {
				added = append(added[:idx], added[idx+1:]...)
			}
		}
	}

	f := `+%v
 %v
-%v`
	return fmt.Sprintf(f, added, existsBoth, deleted)
}

func init() {
	incomingWebhookUrl = os.Getenv("SLACK_URL")
}

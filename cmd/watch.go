package cmd

import (
	//"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

var hostedZone, domain string

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "The IP watcher",
	Long: `The IP watcher.
This subcommand is intended to run under "consul watch".`,
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(runWatcher())
	},
}

func init() {
	RootCmd.AddCommand(watchCmd)

	watchCmd.Flags().StringVarP(&hostedZone, "hosted-zone", "H", "", "Hosted zone to update")
	watchCmd.Flags().StringVarP(&domain, "domain", "D", "", "Full domain name to keep global IPs")
}

func runWatcher() int {
	svc := route53.New(session.New())
	p1 := &route53.ListHostedZonesByNameInput{
		DNSName: aws.String(hostedZone),
	}
	resp, err := svc.ListHostedZonesByName(p1)

	if err != nil {
		logger.Errorln(err.Error())
		return -1
	}

	var target *route53.HostedZone
	for _, z := range resp.HostedZones {
		if *z.Name == hostedZone+"." {
			target = z
		}
	}
	if target == nil {
		logger.Errorf("Hosted zone not found. response: %v", resp.HostedZones)
		return -1
	}

	logger.Infof("get response: %s(%s)", *target.Name, *target.Id)

	p2 := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(domain + "."),
						Type: aws.String("A"),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String("192.168.0.100"),
							},
							{
								Value: aws.String("192.168.0.101"),
							},
							{
								Value: aws.String("192.168.0.102"),
							},
						},
						TTL: aws.Int64(60),
					},
				},
			},
			Comment: aws.String("Update via collector"),
		},
		HostedZoneId: target.Id,
	}
	r2, err := svc.ChangeResourceRecordSets(p2)
	if err != nil {
		logger.Errorln(err.Error())
		return -1
	}

	logger.Infof("Success: %v", r2.ChangeInfo)

	return 0
}

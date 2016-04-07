package cmd

import (
	//"fmt"
	"bufio"
	"os"

	"github.com/spf13/cobra"
	"github.com/udzura/collector/collectorlib"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

var hostedZone string
var domains []string

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
	watchCmd.Flags().StringSliceVarP(&domains, "domain", "D", []string{}, "Full domain name to keep global IPs")
}

func runWatcher() int {
	reader := bufio.NewReader(os.Stdin)
	_, err := reader.Peek(1)
	if err != nil {
		logger.Errorln(err.Error())
		return -1
	}
	defer os.Stdin.Close()

	req, err := collectorlib.ParseRequest(reader)
	if err != nil {
		logger.Errorln(err.Error())
		return -1
	}

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

	domainModels, err := collectorlib.NewDomains(domains)
	if err != nil {
		logger.Errorln(err.Error())
		return -1
	}

	for _, domain := range domainModels {
		ips := req.IPsByTag(domain.Tag)
		if len(ips) == 0 {
			logger.Warnln("Hey, no Ip included. Skipping for fail-safe.")
			continue
		}
		logger.Infof("IPs: %v", ips)

		var oldIPs []string
		p2 := &route53.ListResourceRecordSetsInput{
			HostedZoneId:    target.Id,
			StartRecordName: aws.String(domain.FQDN + "."),
			StartRecordType: aws.String("A"),
		}
		r2, err := svc.ListResourceRecordSets(p2)
		if err != nil {
			logger.Errorln(err.Error())
			return -1
		}
		if len(r2.ResourceRecordSets) > 0 {
			rrset := r2.ResourceRecordSets[0]
			if *rrset.Type == "A" {
				for _, rr := range rrset.ResourceRecords {
					oldIPs = append(oldIPs, *rr.Value)
				}
			}
		}
		logger.Infof("existing IPs: %v", oldIPs)
		diff := collectorlib.NewDiff(oldIPs, ips)

		if diff.IsChanged() {
			var rrs []*route53.ResourceRecord
			for _, ip := range ips {
				rrs = append(rrs, &route53.ResourceRecord{
					Value: aws.String(ip),
				})
			}
			p3 := &route53.ChangeResourceRecordSetsInput{
				ChangeBatch: &route53.ChangeBatch{
					Changes: []*route53.Change{
						{
							Action: aws.String("UPSERT"),
							ResourceRecordSet: &route53.ResourceRecordSet{
								Name:            aws.String(domain.FQDN + "."),
								Type:            aws.String("A"),
								ResourceRecords: rrs,
								TTL:             aws.Int64(60),
							},
						},
					},
					Comment: aws.String("Update via collector"),
				},
				HostedZoneId: target.Id,
			}
			r3, err := svc.ChangeResourceRecordSets(p3)
			if err != nil {
				logger.Errorln(err.Error())
				return -1
			}

			collectorlib.NotifyToSlack(domain.FQDN, diff)
			logger.Infof("Success: %v", r3.ChangeInfo)
		} else {
			logger.Infof("No change, skipping.")
		}
	}

	return 0
}

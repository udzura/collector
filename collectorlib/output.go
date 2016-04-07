package collectorlib

import "strings"

func FindIPFromOutput(output string) string {
	for _, rec := range strings.Split(output, "\t") {
		pair := strings.SplitN(rec, ":", 2)
		if len(pair) != 2 {
			continue
		}
		if pair[0] == "ipaddr" {
			return pair[1]
		}
	}

	return ""
}

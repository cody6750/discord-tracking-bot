package handlers

import "testing"

func Test_printMap(t *testing.T) {
	type args struct {
		m map[string]string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			args: args{
				m: map[string]string{
					"!help":          "Displays all avaliable commands within the console channel",
					"!status":        "Displays the current status for the tracking bot",
					"!metrics":       "Displays the metrics for the current scrape",
					"!total_metrics": "Displays the total metrics for the web crawl",
					"!shutdown":      "Shuts down the bot",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printMap(tt.args.m)
		})
	}
}

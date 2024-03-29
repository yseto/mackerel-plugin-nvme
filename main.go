package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os/exec"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin"
)

type NvmePlugin struct {
	Device string
	Prefix string
}

type NvmeSmartLogRecord struct {
	PercentageUsed float64 `json:"percent_used"`
	AvailableSpare float64 `json:"avail_spare"`
}

func (n NvmePlugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := strings.Title(n.MetricKeyPrefix())
	return map[string]mp.Graphs{
		"PercentageUsed": {
			Label: labelPrefix + " Percentage Used",
			Unit:  mp.UnitPercentage,
			Metrics: []mp.Metrics{
				{Name: "PercentageUsed", Label: "Percentage Used"},
			},
		},
		"AvailableSpare": {
			Label: labelPrefix + " Available Spare",
			Unit:  mp.UnitPercentage,
			Metrics: []mp.Metrics{
				{Name: "AvailableSpare", Label: "Available Spare"},
			},
		},
	}
}

func (n NvmePlugin) FetchMetrics() (map[string]float64, error) {
	var record NvmeSmartLogRecord

	out, err := exec.Command("nvme", "smart-log", n.Device, "--output-format=json").Output()
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch nvme metrics: %s", err)
	}

	err = json.Unmarshal(out, &record)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch nvme metrics: %s", err)
	}
	return map[string]float64{
		"PercentageUsed": record.PercentageUsed,
		"AvailableSpare": record.AvailableSpare,
	}, nil
}

func (n NvmePlugin) MetricKeyPrefix() string {
	if n.Prefix == "" {
		n.Prefix = "nvme"
	}
	return n.Prefix
}

func main() {
	optDevice := flag.String("device", "/dev/nvme0", "nvme device path")
	optPrefix := flag.String("metric-key-prefix", "nvme", "Metric key prefix")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	n := NvmePlugin{
		Device: *optDevice,
		Prefix: *optPrefix,
	}
	plugin := mp.NewMackerelPlugin(n)
	plugin.Tempfile = *optTempfile
	plugin.Run()
}

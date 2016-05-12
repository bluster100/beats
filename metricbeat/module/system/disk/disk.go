// +build freebsd linux windows

package disk

import (
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/metricbeat/mb"

	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/disk"
)

func init() {
	if err := mb.Registry.AddMetricSet("system", "disk", New); err != nil {
		panic(err)
	}
}

// MetricSet for fetching system disk IO metrics.
type MetricSet struct {
	mb.BaseMetricSet
}

// New is a mb.MetricSetFactory that returns a new MetricSet.
func New(base mb.BaseMetricSet) (mb.MetricSet, error) {
	return &MetricSet{base}, nil
}

// Fetch fetches disk IO metrics from the OS.
func (m *MetricSet) Fetch() ([]common.MapStr, error) {
	stats, err := disk.IOCounters()
	if err != nil {
		return nil, errors.Wrap(err, "disk io counters")
	}

	events := make([]common.MapStr, 0, len(stats))
	for _, counters := range stats {
		event := common.MapStr{
			"name":        counters.Name,
			"read_count":  counters.ReadCount,
			"write_count": counters.WriteCount,
			"read_bytes":  counters.ReadBytes,
			"write_bytes": counters.WriteBytes,
			"read_time":   counters.ReadTime,
			"write_time":  counters.WriteTime,
			"io_time":     counters.IoTime,
		}
		events = append(events, event)

		if counters.SerialNumber != "" {
			event["serial_number"] = counters.SerialNumber
		}
	}

	return events, nil
}

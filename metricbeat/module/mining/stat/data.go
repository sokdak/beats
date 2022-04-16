package stat

import (
	"encoding/json"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/metricbeat/mb"
	"github.com/sokdak/miner-exporter/pkg/dto"
)

func parseResponse(data []byte) []mb.Event {
	status := &dto.Status{}
	err := json.Unmarshal(data, status)
	if err != nil {
		return nil
	}

	events := make([]mb.Event, 0)
	for _, dev := range status.Devices {
		events = append(events, mb.Event{
			RootFields:        common.MapStr{},
			ModuleFields:      common.MapStr{
				"miner.name": status.Miner.Name,
				"miner.version": status.Miner.Version,
				"miner.uptime": status.Miner.Uptime,
				"miner.pool": status.Miner.Pool,
				"miner.address": status.Miner.Address,
				"miner.algorithm": status.Miner.Algorithm,
				"miner.worker": status.Miner.Worker,
 			},
			MetricSetFields:   common.MapStr{
				"gpu.id": dev.GpuId,
				"gpu.model": dev.Name,
				"gpu.temperature.core": dev.CoreTemp,
				"gpu.temperature.memory": dev.MemoryTemp,
				"gpu.fan": dev.FanSpeed,
				"gpu.hashrate": dev.Hashrate,
				"gpu.power": dev.PowerConsumption,
				"share.accepted": dev.ShareAccepted,
				"share.rejected": dev.ShareRejected,
				"share.stale": dev.ShareStale,
			},
		})
	}
	return events
}

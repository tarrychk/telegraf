// +build linux

package system

import (
	"fmt"
	"net"

	"github.com/influxdb/telegraf/plugins"
)

type DockerStats struct {
	ps PS
}

func (_ *DockerStats) Description() string {
	return "Read metrics about docker containers"
}

func (_ *DockerStats) SampleConfig() string { return "" }

func (s *DockerStats) Gather(acc plugins.Accumulator) error {
	containers, err := s.ps.DockerStat()
	if err != nil {
		return fmt.Errorf("error getting docker info: %s", err)
	}

	for _, cont := range containers {
		tags := map[string]string{
			"id":      cont.Id,
			"name":    cont.Name,
			"command": cont.Command,
		}
		for k, v := range cont.Labels {
			tags[k] = v
		}

		cts := cont.CPU

		acc.Add("user", cts.User, tags)
		acc.Add("system", cts.System, tags)
		acc.Add("idle", cts.Idle, tags)
		acc.Add("nice", cts.Nice, tags)
		acc.Add("iowait", cts.Iowait, tags)
		acc.Add("irq", cts.Irq, tags)
		acc.Add("softirq", cts.Softirq, tags)
		acc.Add("steal", cts.Steal, tags)
		acc.Add("guest", cts.Guest, tags)
		acc.Add("guest_nice", cts.GuestNice, tags)

		acc.Add("cache", cont.Mem.Cache, tags)
		acc.Add("rss", cont.Mem.RSS, tags)
		acc.Add("rss_huge", cont.Mem.RSSHuge, tags)
		acc.Add("mapped_file", cont.Mem.MappedFile, tags)
		acc.Add("swap_in", cont.Mem.Pgpgin, tags)
		acc.Add("swap_out", cont.Mem.Pgpgout, tags)
		acc.Add("page_fault", cont.Mem.Pgfault, tags)
		acc.Add("page_major_fault", cont.Mem.Pgmajfault, tags)
		acc.Add("inactive_anon", cont.Mem.InactiveAnon, tags)
		acc.Add("active_anon", cont.Mem.ActiveAnon, tags)
		acc.Add("inactive_file", cont.Mem.InactiveFile, tags)
		acc.Add("active_file", cont.Mem.ActiveFile, tags)
		acc.Add("unevictable", cont.Mem.Unevictable, tags)
		acc.Add("memory_limit", cont.Mem.HierarchicalMemoryLimit, tags)
		acc.Add("total_cache", cont.Mem.TotalCache, tags)
		acc.Add("total_rss", cont.Mem.TotalRSS, tags)
		acc.Add("total_rss_huge", cont.Mem.TotalRSSHuge, tags)
		acc.Add("total_mapped_file", cont.Mem.TotalMappedFile, tags)
		acc.Add("total_swap_in", cont.Mem.TotalPgpgIn, tags)
		acc.Add("total_swap_out", cont.Mem.TotalPgpgOut, tags)
		acc.Add("total_page_fault", cont.Mem.TotalPgFault, tags)
		acc.Add("total_page_major_fault", cont.Mem.TotalPgMajFault, tags)
		acc.Add("total_inactive_anon", cont.Mem.TotalInactiveAnon, tags)
		acc.Add("total_active_anon", cont.Mem.TotalActiveAnon, tags)
		acc.Add("total_inactive_file", cont.Mem.TotalInactiveFile, tags)
		acc.Add("total_active_file", cont.Mem.TotalActiveFile, tags)
		acc.Add("total_unevictable", cont.Mem.TotalUnevictable, tags)

		for _, io := range cont.Net {
			iface, err := net.InterfaceByName(io.Name)
			if err != nil {
				continue
			}

			if iface.Flags&net.FlagLoopback == net.FlagLoopback {
				continue
			}

			if iface.Flags&net.FlagUp == 0 {
				continue
			}

			tags["interface"] = io.Name

			acc.Add("bytes_sent", io.BytesSent, tags)
			acc.Add("bytes_recv", io.BytesRecv, tags)
			acc.Add("packets_sent", io.PacketsSent, tags)
			acc.Add("packets_recv", io.PacketsRecv, tags)
			acc.Add("err_in", io.Errin, tags)
			acc.Add("err_out", io.Errout, tags)
			acc.Add("drop_in", io.Dropin, tags)
			acc.Add("drop_out", io.Dropout, tags)
		}

	}

	return nil
}

func init() {
	plugins.Add("docker", func() plugins.Plugin {
		return &DockerStats{ps: &systemPS{}}
	})
}

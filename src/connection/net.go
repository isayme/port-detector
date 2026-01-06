package connection

import (
	"fmt"
	"sort"
	"syscall"

	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

type ListenPortInfo struct {
	Fd         string
	Family     string
	Type       uint32
	Proto      string
	Port       uint32
	PID        int32
	ProcName   string
	Cmdline    string
	LocalAddr  string
	Username   string
	CPUPercent float64
	MemInfo    *process.MemoryInfoStat
}

func ListeningPorts() ([]ListenPortInfo, error) {
	conns, err := net.Connections("inet")
	if err != nil {
		return nil, err
	}

	result := make([]ListenPortInfo, 0)

	for _, c := range conns {
		if c.Type == syscall.SOCK_STREAM && c.Status != "LISTEN" {
			continue
		}

		if c.Type == syscall.SOCK_DGRAM && c.Status != "" {
			continue
		}

		info := ListenPortInfo{
			//Proto:     langs.Localize(fmt.Sprintf("PROTO_%d", c.Type)),
			Type:      c.Type,
			Family:    formatFamily(c.Family),
			Proto:     formatProto(c.Type),
			Port:      c.Laddr.Port,
			PID:       c.Pid,
			LocalAddr: c.Laddr.IP,
		}

		if c.Pid > 0 {
			if p, err := process.NewProcess(c.Pid); err == nil {
				if name, err := p.Name(); err == nil {
					info.ProcName = name
				}
				if cmd, err := p.Cmdline(); err == nil {
					info.Cmdline = cmd
				}
				info.CPUPercent, _ = p.CPUPercent()
				info.Username, _ = p.Username()
				info.MemInfo, _ = p.MemoryInfo()
			}
		}

		result = append(result, info)
	}

	sort.Slice(result, func(i, j int) bool {
		x := result[i]
		y := result[j]

		if x.Port != y.Port {
			return x.Port < y.Port
		}

		return x.Family < x.Family
	})

	return result, nil
}

var protoMap = map[uint32]string{
	syscall.SOCK_STREAM: "TCP",
	syscall.SOCK_DGRAM:  "UDP",
}

func formatProto(proto uint32) string {
	if v, ok := protoMap[proto]; ok {
		return v
	}

	return fmt.Sprintf("%d", proto)
}

var familyMap = map[uint32]string{
	syscall.AF_INET:  "IPv4",
	syscall.AF_INET6: "IPv6",
}

func formatFamily(family uint32) string {
	if v, ok := familyMap[family]; ok {
		return v
	}

	return fmt.Sprintf("%d", family)
}

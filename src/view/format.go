package view

import (
	"fmt"
	"strconv"
	"syscall"

	"github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/v4/process"
)

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

func formatUint32(v uint32) string {
	return strconv.Itoa(int(v))
}

func formatInt32(v int32) string {
	return strconv.Itoa(int(v))
}

func formatCPU(v float64) string {
	return fmt.Sprintf("%0.4f%%", v)
}

func formatMEM(s *process.MemoryInfoStat) string {
	if s == nil {
		return "-"
	}

	return humanize.Bytes(s.RSS)
}

func formatSwitch(v bool) string {
	if v {
		return "on"
	}
	return "off"
}

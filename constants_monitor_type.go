package roselite

import (
	"errors"
	"strings"
)

var ErrMonitorTypeInvalid = errors.New("invalid monitor type")

type MonitorType uint8

const (
	MonitorTypeHTTP MonitorType = iota
	MonitorTypeICMP
	MonitorTypeUnknown MonitorType = 255
)

func (m MonitorType) String() string {
	switch m {
	case MonitorTypeHTTP:
		return "HTTP"
	case MonitorTypeICMP:
		return "ICMP"
	default:
		return "Unknown"
	}
}

func MonitorTypeFromString(s string) (MonitorType, error) {
	switch strings.ToUpper(s) {
	case "HTTP":
		return MonitorTypeHTTP, nil
	case "ICMP":
		return MonitorTypeICMP, nil
	default:
		return MonitorTypeUnknown, ErrMonitorTypeInvalid
	}
}

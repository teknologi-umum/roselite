package roselite

import "errors"

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
	switch s {
	case "HTTP":
		return MonitorTypeHTTP, nil
	case "ICMP":
		return MonitorTypeICMP, nil
	default:
		return MonitorTypeUnknown, errors.New("unknown monitor type")
	}
}

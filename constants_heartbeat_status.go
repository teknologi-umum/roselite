package roselite

import "strings"

type HeartbeatStatus uint8

const (
	HeartbeatStatusUp HeartbeatStatus = iota
	HeartbeatStatusDown
	HeartbeatStatusUnknown HeartbeatStatus = 255
)

func (h HeartbeatStatus) String() string {
	switch h {
	case HeartbeatStatusUp:
		return "up"
	case HeartbeatStatusDown:
		return "down"
	default:
		return ""
	}
}

func HeartbeatStatusFromString(s string) HeartbeatStatus {
	switch strings.ToLower(s) {
	case "up":
		return HeartbeatStatusUp
	case "down":
		return HeartbeatStatusDown
	default:
		return HeartbeatStatusUnknown
	}
}

func (h HeartbeatStatus) MarshalJSON() ([]byte, error) {
	return []byte(`"` + h.String() + `"`), nil
}

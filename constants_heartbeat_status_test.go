package roselite_test

import (
	"bytes"
	"testing"

	"github.com/teknologi-umum/roselite"
)

func TestHeartbeatStatus_String(t *testing.T) {
	testCases := []struct {
		status   roselite.HeartbeatStatus
		expected string
	}{
		{
			status:   roselite.HeartbeatStatusUp,
			expected: "up",
		},
		{
			status:   roselite.HeartbeatStatusDown,
			expected: "down",
		},
		{
			status:   roselite.HeartbeatStatusUnknown,
			expected: "",
		},
	}

	for _, testCase := range testCases {
		if testCase.status.String() != testCase.expected {
			t.Errorf("expected %s, got %s", testCase.expected, testCase.status.String())
		}
	}
}

func TestHeartbeatStatusFromString(t *testing.T) {
	testCases := []struct {
		status   string
		expected roselite.HeartbeatStatus
	}{
		{
			status:   "up",
			expected: roselite.HeartbeatStatusUp,
		},
		{
			status:   "DOWN",
			expected: roselite.HeartbeatStatusDown,
		},
		{
			status:   "down",
			expected: roselite.HeartbeatStatusDown,
		},
		{
			status:   "unknown",
			expected: roselite.HeartbeatStatusUnknown,
		},
	}

	for _, testCase := range testCases {
		if roselite.HeartbeatStatusFromString(testCase.status) != testCase.expected {
			t.Errorf("expected %s, got %s", testCase.expected, testCase.status)
		}
	}
}

func TestHeartbeatStatus_MarshalJSON(t *testing.T) {
	testCases := []struct {
		status   roselite.HeartbeatStatus
		expected string
	}{
		{
			status:   roselite.HeartbeatStatusUp,
			expected: "\"up\"",
		},
	}

	for _, testCase := range testCases {
		b, err := testCase.status.MarshalJSON()
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if bytes.Compare(b, []byte(testCase.expected)) != 0 {
			t.Errorf("expected %s, got %s", testCase.expected, string(b))
		}
	}
}

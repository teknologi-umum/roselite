package roselite_test

import (
	"errors"
	"testing"

	"github.com/teknologi-umum/roselite"
)

func TestMonitorType_String(t *testing.T) {
	testCases := []struct {
		monitorType roselite.MonitorType
		expected    string
	}{
		{
			monitorType: roselite.MonitorTypeHTTP,
			expected:    "HTTP",
		},
		{
			monitorType: roselite.MonitorTypeICMP,
			expected:    "ICMP",
		},
		{
			monitorType: roselite.MonitorTypeUnknown,
			expected:    "Unknown",
		},
	}

	for _, testCase := range testCases {
		if testCase.monitorType.String() != testCase.expected {
			t.Errorf("expected %s, got %s", testCase.expected, testCase.monitorType.String())
		}
	}
}

func TestMonitorTypeFromString(t *testing.T) {
	testCases := []struct {
		monitorType string
		expected    roselite.MonitorType
		expectedErr error
	}{
		{
			monitorType: "HTTP",
			expected:    roselite.MonitorTypeHTTP,
			expectedErr: nil,
		},
		{
			monitorType: "ICMP",
			expected:    roselite.MonitorTypeICMP,
			expectedErr: nil,
		},
		{
			monitorType: "",
			expected:    roselite.MonitorTypeUnknown,
			expectedErr: roselite.ErrMonitorTypeInvalid,
		},
	}

	for _, testCase := range testCases {
		monitorType, err := roselite.MonitorTypeFromString(testCase.monitorType)
		if !errors.Is(err, testCase.expectedErr) {
			t.Errorf("expected %s, got %s", testCase.expected, testCase.monitorType)
		}

		if monitorType != testCase.expected {
			t.Errorf("expected %s, got %s", testCase.expected, testCase.monitorType)
		}
	}
}

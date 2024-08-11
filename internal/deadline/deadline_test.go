package deadline

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseDeadlineValue(t *testing.T) {
	testCases := []struct {
		dlType   string
		value    string
		expected string
		hasError bool
	}{
		{
			dlType:   TypeTime,
			value:    "120",
			expected: "120000",
			hasError: false,
		},
		{
			dlType:   TypeTime,
			value:    "0",
			expected: "",
			hasError: true,
		},
		{
			dlType:   TypeTime,
			value:    "",
			expected: "",
			hasError: true,
		},
		{
			dlType:   "UNSUPPORTED",
			value:    "ABC",
			expected: "ABC",
			hasError: false,
		},
	}

	for i, tc := range testCases {
		res, err := ParseDeadlineValue(tc.dlType, tc.value)
		if tc.hasError {
			assert.NotNil(t, err, fmt.Sprintf("[%d] error is expected", i))
		} else {
			assert.Equal(t, tc.expected, res, fmt.Sprintf("[%d] the result does not match with expected", i))
		}
	}
}

func TestFindActivated(t *testing.T) {
	currentTime := time.Date(2006, 1, 2, 15, 0, 0, 0, time.Local)

	testCases := []struct {
		deadlines   []Deadline
		placedAt    time.Time
		expected    *Deadline
		expectError bool
	}{
		{
			deadlines: []Deadline{
				{Type: TypeTime, Action: ActionSellByMarket, Value: "120000"},
				{Type: TypeTime, Action: ActionCancel, Value: "60000"},
			},
			placedAt:    time.Date(2006, 1, 2, 14, 58, 30, 0, time.Local),
			expected:    &Deadline{Type: TypeTime, Action: ActionCancel, Value: "60000"},
			expectError: false,
		},
		{
			deadlines: []Deadline{
				{Type: TypeTime, Action: ActionSellByMarket, Value: "120000"},
				{Type: TypeTime, Action: ActionCancel, Value: "60000"},
			},
			placedAt:    time.Date(2006, 1, 2, 14, 59, 30, 0, time.Local),
			expected:    nil,
			expectError: false,
		},
		{
			deadlines: []Deadline{
				{Type: TypeTime, Action: ActionSellByMarket, Value: "120002"},
				{Type: TypeTime, Action: ActionCancel, Value: "120001"},
			},
			placedAt:    time.Date(2006, 1, 2, 14, 50, 30, 0, time.Local),
			expected:    &Deadline{Type: TypeTime, Action: ActionSellByMarket, Value: "120002"},
			expectError: false,
		},
		{
			deadlines:   []Deadline{},
			placedAt:    time.Date(2006, 1, 2, 14, 50, 30, 0, time.Local),
			expected:    nil,
			expectError: false,
		},
		{
			deadlines: []Deadline{
				{Type: "UNSUPPORTED", Action: ActionSellByMarket, Value: "120002"},
				{Type: TypeTime, Action: ActionCancel, Value: "120001"},
			},
			placedAt:    time.Date(2006, 1, 2, 14, 50, 30, 0, time.Local),
			expected:    nil,
			expectError: true,
		},
		{
			deadlines: []Deadline{
				{Type: TypeTime, Action: ActionCancel, Value: ""},
			},
			placedAt:    time.Date(2006, 1, 2, 14, 50, 30, 0, time.Local),
			expected:    nil,
			expectError: true,
		},
	}

	for i, tc := range testCases {
		res, err := FindActivated(tc.deadlines, currentTime, tc.placedAt)
		if tc.expectError {
			assert.NotNil(t, err, fmt.Sprintf("[%d] error is expected", i))
		} else {
			assert.Nil(t, err, fmt.Sprintf("[%d] error is not expected", i))
			assert.Equal(t, tc.expected, res, fmt.Sprintf("[%d] the result does not match with expected", i))
		}
		if !tc.expectError && err != nil {
			t.Logf("[%d] unexpected error is returned: %v", i, err)
		}
	}
}

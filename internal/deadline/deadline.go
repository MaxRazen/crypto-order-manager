package deadline

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

const (
	TypeTime           = "TIME"
	ActionCancel       = "CANCEL"
	ActionSellByMarket = "SELL_BY_MARKET"
)

type Deadline struct {
	Type   string `datastore:"type,noindex"`
	Action string `datastore:"action,noindex"`
	// empty string or duration in miliseconds
	Value string `datastore:"value,noindex"`
}

func ParseDeadlineValue(tp, value string) (string, error) {
	if tp != TypeTime {
		return value, nil
	}

	if value == "" || value == "0" {
		return "", errors.New("deadline value must be specified for this type")
	}

	v, err := strconv.Atoi(value)
	if err != nil {
		return "", err
	}

	ms := time.Duration(time.Duration(v) * time.Second).Milliseconds()
	return strconv.FormatInt(ms, 10), nil
}

func FindActivated(deadlines []Deadline, now, startTime time.Time) (*Deadline, error) {
	for _, d := range deadlines {
		if d.Type == TypeTime {
			v, err := strconv.Atoi(d.Value)
			if err != nil {
				return nil, err
			}
			if now.After(startTime.Add(time.Duration(v) * time.Millisecond)) {
				return &d, nil
			}
		} else {
			return nil, fmt.Errorf("deadline: unsupported type found '%s'", d.Type)
		}
	}
	return nil, nil
}

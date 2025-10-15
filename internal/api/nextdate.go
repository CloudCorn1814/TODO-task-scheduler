package api

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const maxDays = 400
const dateFormat = "20060102"

func nextByDays(start, now time.Time, step int) time.Time {
	date := start.AddDate(0, 0, step)
	for !date.After(now) {
		date = date.AddDate(0, 0, step)
	}
	return date
}

func nextByYears(start, now time.Time, step int) time.Time {
	date := start.AddDate(step, 0, 0)
	for !date.After(now) {
		date = date.AddDate(step, 0, 0)
	}
	return date
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	repeat = strings.TrimSpace(repeat)
	if repeat == "" {
		return "", errors.New("repeat rule is empty")
	}

	start, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("bad start date: %w", err)
	}

	parts := strings.Fields(repeat)
	kind := parts[0]

	switch kind {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("invalid format for 'd'")
		}
		n, err := strconv.Atoi(parts[1])
		if err != nil || n < 1 || n > maxDays {
			return "", fmt.Errorf("invalid number of days, must be from 1 to %d", maxDays)
		}
		next := nextByDays(start, now, n)
		return next.Format(dateFormat), nil

	case "y":
		if len(parts) != 1 {
			return "", errors.New("invalid format for 'y'")
		}
		next := nextByYears(start, now, 1)
		return next.Format(dateFormat), nil

	default:
		return "", fmt.Errorf("unknown repeat rule: %q", kind)
	}
}

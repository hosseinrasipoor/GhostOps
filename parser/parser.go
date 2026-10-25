package parser

import (
	"errors"
	"strings"
	"time"
)

type LogEntry struct {
	Timestamp time.Time
	Level     string
	Service   string
	Message   string
}

func ParseLine(line string) (*LogEntry, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, errors.New("empty line")
	}

	parts := strings.Split(line, " ")
	if len(parts) < 4 {
		return nil, errors.New("bad line")
	}

	ts, err := time.Parse(time.RFC3339, parts[0])
	if err != nil {
		return nil, err
	}

	return &LogEntry{
		Timestamp: ts,
		Level:     parts[1],
		Service:   parts[2],
		Message:   parts[3],
	}, nil
}

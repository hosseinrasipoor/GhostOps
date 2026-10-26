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

	parts := strings.SplitN(line, " ", 4)
	if len(parts) < 4 {
		return nil, errors.New("bad line")
	}

	lvl := parts[1]
	if lvl != "INFO" && lvl != "WARN" && lvl != "ERROR" && lvl != "FATAL" {
		return nil, errors.New("bad level")
	}

	ts, err := time.Parse(time.RFC3339, parts[0])
	if err != nil {
		return nil, err
	}

	svc := parts[2]
	msg := parts[3]
	if svc == "" || msg == "" {
		return nil, errors.New("missing field")
	}

	return &LogEntry{
		Timestamp: ts,
		Level:     lvl,
		Service:   svc,
		Message:   msg,
	}, nil
}

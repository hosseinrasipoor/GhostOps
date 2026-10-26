package analyzer

import (
	"ghostops/parser"
	"strings"
	"time"
)

type Incident struct {
	Type      string    `json:"type"`
	Severity  string    `json:"severity"`
	Service   string    `json:"service"`
	Count     int       `json:"count"`
	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
	Desc      string    `json:"desc"`
}

type Result struct {
	TotalLines     int            `json:"total_lines"`
	ValidEntries   int            `json:"valid_entries"`
	MalformedLines int            `json:"malformed_lines"`
	Severities     map[string]int `json:"severities"`
	Services       map[string]int `json:"services"`
	Incidents      []Incident     `json:"incidents"`
}

type Analyzer struct {
	res Result
	m   map[string]*Incident
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{
		res: Result{
			Severities: map[string]int{
				"INFO":  0,
				"WARN":  0,
				"ERROR": 0,
				"FATAL": 0,
			},
			Services:  make(map[string]int),
			Incidents: make([]Incident, 0),
		},
		m: make(map[string]*Incident),
	}
}

func (a *Analyzer) AddBad() {
	a.res.TotalLines++
	a.res.MalformedLines++
}

func checkIncident(lvl, msg string) (string, string, bool) {
	lower := strings.ToLower(msg)

	if lvl == "FATAL" {
		if strings.Contains(lower, "pool") || strings.Contains(lower, "database") {
			return "db crash", "database pool crash", true
		}
		return "fatal crash", "system fatal crash", true
	}

	if lvl == "ERROR" {
		if strings.Contains(lower, "connect to database") || strings.Contains(lower, "timeout after") {
			return "db timeout", "database timeout", true
		}
		if strings.Contains(lower, "gateway") || strings.Contains(lower, "502") {
			return "gateway outage", "payment gateway outage", true
		}
		if strings.Contains(lower, "sql") || strings.Contains(lower, "syntax") {
			return "sql error", "sql query error", true
		}
		if strings.Contains(lower, "nullpointer") {
			return "null pointer", "null pointer exception", true
		}
		return "app error", "application error", true
	}

	if lvl == "WARN" {
		if strings.Contains(lower, "memory") {
			return "high memory", "memory warning", true
		}
		if strings.Contains(lower, "rate limit") {
			return "rate limit", "rate limit reached", true
		}
		if strings.Contains(lower, "deprecated") {
			return "deprecated api", "deprecated api call", true
		}
		if strings.Contains(lower, "pool") || strings.Contains(lower, "threshold") {
			return "db pool warning", "database pool warning", true
		}
	}

	return "", "", false
}

func (a *Analyzer) AddGood(e *parser.LogEntry) {
	a.res.TotalLines++
	a.res.ValidEntries++

	a.res.Severities[e.Level]++
	a.res.Services[e.Service]++

	name, desc, ok := checkIncident(e.Level, e.Message)
	if ok {
		k := name + ":" + e.Service
		inc, exists := a.m[k]
		if !exists {
			inc = &Incident{
				Type:      name,
				Severity:  e.Level,
				Service:   e.Service,
				Count:     1,
				FirstSeen: e.Timestamp,
				LastSeen:  e.Timestamp,
				Desc:      desc,
			}
			a.m[k] = inc
		} else {
			inc.Count++
			if e.Timestamp.Before(inc.FirstSeen) {
				inc.FirstSeen = e.Timestamp
			}
			if e.Timestamp.After(inc.LastSeen) {
				inc.LastSeen = e.Timestamp
			}
		}
	}
}

func (a *Analyzer) GetResult() Result {
	r := a.res
	for _, inc := range a.m {
		r.Incidents = append(r.Incidents, *inc)
	}
	return r
}

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// submission is one estimate request as written to the submissions log.
type submission struct {
	Time    time.Time `json:"time"`
	Name    string    `json:"name"`
	Phone   string    `json:"phone"`
	Email   string    `json:"email"`
	Address string    `json:"address"`
	Service string    `json:"service"`
	Message string    `json:"message"`
}

// submissionLog appends each estimate request to a JSON Lines file so there
// is a durable copy independent of email. Container stdout is not one: every
// deploy recreates the container and its log goes with it.
//
// A nil *submissionLog is valid and does nothing (SUBMISSIONS_FILE unset).
type submissionLog struct {
	mu sync.Mutex
	f  *os.File
}

// openSubmissionLog opens path for appending, creating it 0600 since it holds
// names, phone numbers and home addresses. An empty path returns nil, nil.
func openSubmissionLog(path string) (*submissionLog, error) {
	if path == "" {
		return nil, nil
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o600)
	if err != nil {
		return nil, fmt.Errorf("open submissions file: %w", err)
	}
	return &submissionLog{f: f}, nil
}

// append writes s as one JSON line and fsyncs, so a request is on disk
// before the visitor sees the success message.
func (l *submissionLog) append(s submission) error {
	if l == nil {
		return nil
	}
	line, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("marshal submission: %w", err)
	}
	line = append(line, '\n')

	l.mu.Lock()
	defer l.mu.Unlock()
	if _, err := l.f.Write(line); err != nil {
		return fmt.Errorf("write submission: %w", err)
	}
	if err := l.f.Sync(); err != nil {
		return fmt.Errorf("sync submissions file: %w", err)
	}
	return nil
}

package runs

import (
    "errors"
    "sync"
    "time"

    "github.com/google/uuid"
)

type RunStore struct {
    mu   sync.RWMutex
    runs map[string]*Run
}

func NewRunStore() *RunStore {
    return &RunStore{runs: make(map[string]*Run)}
}

// RunRequest is the payload the UI/backend uses to start a new run.
type RunRequest struct {
    Framework  string `json:"framework"`
    SuiteId    string `json:"suiteId"`
    ScenarioId string `json:"scenarioId"`
    // Executor preference: "docker" or "k8s". Optional.
    Executor   string `json:"executor,omitempty"`
}

func (s *RunStore) StartRun(req RunRequest) *Run {
    id := uuid.New().String()
    now := time.Now().UTC().Format(time.RFC3339)

    run := &Run{
        ID:         id,
        Status:     "running",
        StartedAt:  now,
        Framework:  req.Framework,
        SuiteId:    req.SuiteId,
        ScenarioId: req.ScenarioId,
        Executor:   req.Executor,
    }

    s.mu.Lock()
    s.runs[id] = run
    s.mu.Unlock()

    return run
}

func (s *RunStore) CompleteRun(id, status, reportURL string) {
    s.mu.Lock()
    defer s.mu.Unlock()

    run := s.runs[id]
    if run == nil {
        return
    }

    run.Status = status
    run.FinishedAt = time.Now().UTC().Format(time.RFC3339)
    run.ReportURL = reportURL
}

func (s *RunStore) GetRun(id string) (*Run, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    run := s.runs[id]
    if run == nil {
        return nil, errors.New("run not found")
    }
    return run, nil
}

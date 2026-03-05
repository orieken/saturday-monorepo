package registry

import "errors"

// Registry holds in-memory metadata for test suites per framework.
type Registry struct {
    cucumberSuites map[string]*CucumberIndex
}

func NewRegistry() *Registry {
    return &Registry{
        cucumberSuites: make(map[string]*CucumberIndex),
    }
}

func (r *Registry) ListSuites(framework string) []string {
    if framework == "cucumber" {
        res := make([]string, 0, len(r.cucumberSuites))
        for k := range r.cucumberSuites {
            res = append(res, k)
        }
        return res
    }
    return []string{}
}

func (r *Registry) GetSuiteIndex(framework, suiteId string) (*CucumberIndex, error) {
    if framework != "cucumber" {
        return nil, errors.New("unsupported framework")
    }
    idx, ok := r.cucumberSuites[suiteId]
    if !ok {
        return nil, errors.New("suite not found")
    }
    return idx, nil
}

// RegisterCucumberIndex registers or replaces a CucumberIndex in memory.
func (r *Registry) RegisterCucumberIndex(idx *CucumberIndex) {
    if idx == nil {
        return
    }
    if r.cucumberSuites == nil {
        r.cucumberSuites = make(map[string]*CucumberIndex)
    }
    r.cucumberSuites[idx.SuiteId] = idx
}

package registry

import (
    "encoding/json"
    "os"
)

// LoadCucumberIndexFromJSON bootstraps the registry from a local JSON file.
func (r *Registry) LoadCucumberIndexFromJSON(path string) error {
    data, err := os.ReadFile(path)
    if err != nil {
        return err
    }
    var idx CucumberIndex
    if err := json.Unmarshal(data, &idx); err != nil {
        return err
    }
    r.RegisterCucumberIndex(&idx)
    return nil
}

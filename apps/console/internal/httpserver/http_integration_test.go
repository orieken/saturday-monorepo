package httpserver_test

import (
    "bytes"
    "encoding/json"
    "io"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"
    "time"

    httpserver "saturday/console/internal/httpserver"
    "saturday/console/internal/registry"
    "saturday/console/internal/runs"
)

// testIndex is a minimal Cucumber index fixture with one feature and one scenario.
var testIndex = &registry.CucumberIndex{
    Framework: "cucumber",
    SuiteId:   "demo-suite",
    Features: []registry.CucumberFeatureRef{
        {
            Id:          "product-search",
            Name:        "Product Search",
            File:        "product_search.feature",
            Description: "Search for products",
            Scenarios: []registry.CucumberScenarioRef{
                {
                    Id:   "searching-for-a-product-by-name",
                    Line: 6,
                    Name: "Searching for a product by name",
                    File: "product_search.feature",
                    Tags: []string{"@smoke"},
                    Steps: []registry.CucumberStepRef{
                        {Text: "Given I am on the demo shop home page", Line: 7},
                        {Text: "When I search for \"keyboard\"", Line: 8},
                        {Text: "Then I should see results containing \"keyboard\"", Line: 9},
                    },
                },
            },
        },
    },
}

func newTestServer() *httptest.Server {
    reg := registry.NewRegistry()
    runStore := runs.NewRunStore()
    h := httpserver.New(reg, runStore)
    return httptest.NewServer(h)
}

func TestListFrameworks(t *testing.T) {
    srv := newTestServer()
    defer srv.Close()

    resp, err := http.Get(srv.URL + "/api/frameworks")
    if err != nil {
        t.Fatalf("GET /api/frameworks error: %v", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        t.Fatalf("expected 200, got %d", resp.StatusCode)
    }
    var arr []string
    if err := json.NewDecoder(resp.Body).Decode(&arr); err != nil {
        t.Fatalf("decode: %v", err)
    }
    if len(arr) != 1 || arr[0] != "cucumber" {
        t.Fatalf("unexpected frameworks: %v", arr)
    }
}

func TestIndexIngestAndList(t *testing.T) {
    srv := newTestServer()
    defer srv.Close()

    // Initially suites list is empty
    resp, err := http.Get(srv.URL + "/api/frameworks/cucumber/suites")
    if err != nil {
        t.Fatalf("GET suites error: %v", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        t.Fatalf("expected 200, got %d", resp.StatusCode)
    }
    var suites []string
    if err := json.NewDecoder(resp.Body).Decode(&suites); err != nil {
        t.Fatalf("decode suites: %v", err)
    }
    if len(suites) != 0 {
        t.Fatalf("expected no suites, got %v", suites)
    }

    // POST the index
    var buf bytes.Buffer
    if err := json.NewEncoder(&buf).Encode(testIndex); err != nil {
        t.Fatalf("encode index: %v", err)
    }
    resp2, err := http.Post(srv.URL+"/api/cucumber/index", "application/json", &buf)
    if err != nil {
        t.Fatalf("POST index error: %v", err)
    }
    defer resp2.Body.Close()
    if resp2.StatusCode != http.StatusOK {
        b, _ := io.ReadAll(resp2.Body)
        t.Fatalf("expected 200 posting index, got %d body=%s", resp2.StatusCode, string(b))
    }

    // Suites should now include demo-suite
    resp3, err := http.Get(srv.URL + "/api/frameworks/cucumber/suites")
    if err != nil {
        t.Fatalf("GET suites (after) error: %v", err)
    }
    defer resp3.Body.Close()
    if resp3.StatusCode != http.StatusOK {
        t.Fatalf("expected 200, got %d", resp3.StatusCode)
    }
    suites = nil
    if err := json.NewDecoder(resp3.Body).Decode(&suites); err != nil {
        t.Fatalf("decode suites(after): %v", err)
    }
    if len(suites) != 1 || suites[0] != testIndex.SuiteId {
        t.Fatalf("unexpected suites: %v", suites)
    }

    // Fetch scenarios for the suite
    resp4, err := http.Get(srv.URL + "/api/frameworks/cucumber/suites/" + testIndex.SuiteId + "/scenarios")
    if err != nil {
        t.Fatalf("GET scenarios error: %v", err)
    }
    defer resp4.Body.Close()
    if resp4.StatusCode != http.StatusOK {
        t.Fatalf("expected 200, got %d", resp4.StatusCode)
    }
    var got registry.CucumberIndex
    if err := json.NewDecoder(resp4.Body).Decode(&got); err != nil {
        t.Fatalf("decode scenarios: %v", err)
    }
    if got.SuiteId != testIndex.SuiteId || len(got.Features) != 1 || len(got.Features[0].Scenarios) != 1 {
        t.Fatalf("unexpected index: %+v", got)
    }
}

// TestRunEndpoint is optional and skipped unless E2E_RUNS=1 is set.
// It will likely result in a failed run status unless the environment has Node + Playwright + cucumber-js available.
func TestRunEndpoint(t *testing.T) {
    if os.Getenv("E2E_RUNS") != "1" {
        t.Skip("skipping run endpoint test; set E2E_RUNS=1 to enable")
    }
    srv := newTestServer()
    defer srv.Close()

    // Ingest index first
    var buf bytes.Buffer
    _ = json.NewEncoder(&buf).Encode(testIndex)
    resp, err := http.Post(srv.URL+"/api/cucumber/index", "application/json", &buf)
    if err != nil || resp.StatusCode != http.StatusOK {
        t.Fatalf("failed to post index (status=%d err=%v)", resp.StatusCode, err)
    }
    _ = resp.Body.Close()

    // Start a run
    payload := map[string]string{
        "framework":  "cucumber",
        "suiteId":    testIndex.SuiteId,
        "scenarioId": testIndex.Features[0].Scenarios[0].Id,
    }
    var runBuf bytes.Buffer
    _ = json.NewEncoder(&runBuf).Encode(payload)
    resp2, err := http.Post(srv.URL+"/api/runs", "application/json", &runBuf)
    if err != nil {
        t.Fatalf("POST /api/runs error: %v", err)
    }
    defer resp2.Body.Close()
    if resp2.StatusCode != http.StatusOK {
        b, _ := io.ReadAll(resp2.Body)
        t.Fatalf("expected 200 starting run, got %d body=%s", resp2.StatusCode, string(b))
    }
    var run runs.Run
    if err := json.NewDecoder(resp2.Body).Decode(&run); err != nil {
        t.Fatalf("decode run: %v", err)
    }
    if run.ID == "" {
        t.Fatalf("expected run id")
    }

    // Poll until status is not running (with timeout)
    deadline := time.Now().Add(15 * time.Second)
    for {
        time.Sleep(200 * time.Millisecond)
        resp3, err := http.Get(srv.URL + "/api/runs/" + run.ID)
        if err != nil {
            t.Fatalf("GET run error: %v", err)
        }
        if resp3.StatusCode != http.StatusOK {
            t.Fatalf("GET run unexpected status: %d", resp3.StatusCode)
        }
        var cur runs.Run
        if err := json.NewDecoder(resp3.Body).Decode(&cur); err != nil {
            _ = resp3.Body.Close()
            t.Fatalf("decode run poll: %v", err)
        }
        _ = resp3.Body.Close()
        if cur.Status != "running" {
            if cur.ReportURL == "" {
                t.Fatalf("expected reportUrl when finished: %+v", cur)
            }
            break
        }
        if time.Now().After(deadline) {
            t.Fatalf("run did not complete in time")
        }
    }
}

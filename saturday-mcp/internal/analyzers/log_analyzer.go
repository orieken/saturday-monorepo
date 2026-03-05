package analyzers

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/orieken/saturday-mcp/internal/logging"
)

// TestLogAnalyzer parses test execution logs
type TestLogAnalyzer struct {
	logger *logging.Logger
}

// TestFailureInfo represents parsed failure details
type TestFailureInfo struct {
	Title    string `json:"title"`
	File     string `json:"file"`
	Line     int    `json:"line"`
	Error    string `json:"error"`
	Code     string `json:"code,omitempty"` // The failing line of code if captured
}

// NewTestLogAnalyzer creates a new analyzer
func NewTestLogAnalyzer(logger *logging.Logger) *TestLogAnalyzer {
	return &TestLogAnalyzer{
		logger: logger,
	}
}

// Parse extracts failure details from Playwright output
func (a *TestLogAnalyzer) Parse(logOutput string) ([]TestFailureInfo, error) {
	var failures []TestFailureInfo

	// Regex to capture standard Playwright error header
	// Example: 1) [chromium] › tests/login.spec.ts:15:5 › Login Flow
	headerRegex := regexp.MustCompile(`\d+\)\s+\[.*?\]\s+›\s+(.*?):(\d+):\d+\s+›\s+(.*)`)
	
	lines := strings.Split(logOutput, "\n")
	
	var currentFailure *TestFailureInfo

	for i, line := range lines {
		// remove colors/formatting roughly (not perfect but helpful)
		cleanLine := stripAnsi(line)
		
		matches := headerRegex.FindStringSubmatch(cleanLine)
		if len(matches) > 3 {
			// Save previous failure if exists
			if currentFailure != nil {
				failures = append(failures, *currentFailure)
			}

			lineNum, _ := strconv.Atoi(matches[2])

			currentFailure = &TestFailureInfo{
				File:  matches[1],
				Title: matches[3],
				Line:  lineNum,
			}
			continue
		}

		if currentFailure != nil {
			// Heuristic: Capture the Error message (usually indent followed by Error type)
			if currentFailure.Error == "" && strings.TrimSpace(cleanLine) != "" && !strings.Contains(cleanLine, "==") {
				currentFailure.Error = strings.TrimSpace(cleanLine)
			}
			
			// Heuristic: Capture the Code snippet
			// Example:   15 |     await page.click('button#submit');
			if strings.Contains(cleanLine, "|") {
				parts := strings.SplitN(cleanLine, "|", 2)
				if len(parts) == 2 {
					codeLineNumStr := strings.TrimSpace(parts[0])
					codeLineNum, _ := strconv.Atoi(codeLineNumStr)
					
					if codeLineNum == currentFailure.Line {
						currentFailure.Code = strings.TrimSpace(parts[1])
					}
				}
			}
		}

		// Look ahead constraint? No, basic parsing for now.
		if i == len(lines)-1 && currentFailure != nil {
			failures = append(failures, *currentFailure)
		}
	}

	return failures, nil
}

func stripAnsi(str string) string {
	const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"
	var re = regexp.MustCompile(ansi)
	return re.ReplaceAllString(str, "")
}

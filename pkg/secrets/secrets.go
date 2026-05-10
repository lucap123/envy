package secrets

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Finding struct {
	File    string
	Line    int
	Pattern string
	Snippet string // redacted preview
}

func (f Finding) String() string {
	return fmt.Sprintf("  [%s] %s:%d  →  %s", f.Pattern, f.File, f.Line, f.Snippet)
}

type rule struct {
	name string
	re   *regexp.Regexp
}

var rules = []rule{
	{"AWS Access Key", regexp.MustCompile(`AKIA[0-9A-Z]{16}`)},
	{"AWS Secret Key", regexp.MustCompile(`(?i)aws_secret_access_key\s*[=:]\s*[a-zA-Z0-9/+=]{40}`)},
	{"GitHub Token", regexp.MustCompile(`gh[pors]_[a-zA-Z0-9]{36}`)},
	{"Slack Token", regexp.MustCompile(`xox[baprs]-[0-9a-zA-Z]{10,48}`)},
	{"Stripe Secret Key", regexp.MustCompile(`sk_(test|live)_[0-9a-zA-Z]{24,}`)},
	{"Private Key", regexp.MustCompile(`-----BEGIN [A-Z ]+ PRIVATE KEY-----`)},
	{"Generic Secret", regexp.MustCompile(`(?i)(secret|password|passwd|api_key|apikey|auth_token)\s*[=:]\s*["']?[a-zA-Z0-9_\-]{16,}["']?`)},
	{"Generic Bearer Token", regexp.MustCompile(`(?i)bearer\s+[a-zA-Z0-9\-_=.]{20,}`)},
}

// extensions to skip entirely (binary / lock files)
var skipExt = map[string]bool{
	".png": true, ".jpg": true, ".jpeg": true, ".gif": true,
	".pdf": true, ".zip": true, ".tar": true, ".gz": true,
	".sum": true, ".lock": true, ".ico": true, ".woff": true, ".woff2": true,
}

func ScanFile(path string) ([]Finding, error) {
	// Skip binary / irrelevant file types
	lower := strings.ToLower(path)
	for ext := range skipExt {
		if strings.HasSuffix(lower, ext) {
			return nil, nil
		}
	}

	// Skip .env.example (it's supposed to have placeholder keys)
	if strings.HasSuffix(lower, ".env.example") {
		return nil, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var findings []Finding
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		for _, r := range rules {
			if r.re.MatchString(line) {
				findings = append(findings, Finding{
					File:    path,
					Line:    lineNum,
					Pattern: r.name,
					Snippet: redact(line),
				})
				break // one finding per line is enough
			}
		}
	}

	return findings, scanner.Err()
}

// redact shows the start of the line but masks the sensitive value
func redact(line string) string {
	line = strings.TrimSpace(line)
	if len(line) > 60 {
		return line[:40] + "..." + strings.Repeat("*", 8)
	}
	// Mask everything after = or :
	for i, c := range line {
		if c == '=' || c == ':' {
			return line[:i+1] + " [REDACTED]"
		}
	}
	return "[REDACTED]"
}

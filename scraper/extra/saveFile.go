package extra

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pahulgogna/evoAI_Web/scraper/src/customTypes"
)

var invalidFilenameChars = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

func filenameFromURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		// Fallback if it's not a valid URL
		safe := invalidFilenameChars.ReplaceAllString(raw, "_")
		return safe + ".html"
	}

	// Use host + path as base name
	path := u.Path
	if path == "" || path == "/" {
		path = "index"
	}
	base := u.Host + path

	// Replace bad characters with _
	base = invalidFilenameChars.ReplaceAllString(base, "_")

	// Trim leading/trailing underscores
	base = strings.Trim(base, "_")

	if base == "" {
		base = "page"
	}

	return base + ".html"
}

func WritePageToFile(p customTypes.Page) (string, error) {
	// Ensure output directory exists
	outDir := "output"
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return "", fmt.Errorf("creating output dir: %w", err)
	}

	// Build file path
	filename := filenameFromURL(p.Source)
	fullPath := filepath.Join(outDir, filename)

	// Write file
	if err := os.WriteFile(fullPath, []byte(p.Body), 0o644); err != nil {
		return "", fmt.Errorf("writing file: %w", err)
	}

	return fullPath, nil
}

func GetJSON(data interface{}) string {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return ""
	}

	return string(jsonBytes)
}

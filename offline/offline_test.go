package offline

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/terrable-dev/terrable/config"
)

func TestPrintConfig(t *testing.T) {
	testConfig := config.TerrableConfig{
		Handlers: []config.HandlerMapping{
			{
				Name:   "Handler1",
				Source: "source1",
				Http: map[string]string{
					"GET":  "/path1",
					"POST": "/path1",
				},
			},
			{
				Name:   "Handler2",
				Source: "source2",
				Http: map[string]string{
					"GET": "/path2",
				},
			},
		},
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the function
	printConfig(testConfig, 1234)

	// Restore stdout
	w.Close()
	os.Stdout = old

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	expectedEndpoints := 3
	expectedLine := fmt.Sprintf("%d Endpoint(s) to prepare...", expectedEndpoints)
	if !strings.Contains(output, expectedLine) {
		t.Errorf("Expected output to contain '%s', but it doesn't.\nActual output:\n%s", expectedLine, output)
	}

	expectedLines := []string{
		"Starting terrable local server...",
		"   GET   http://localhost:1234/path1",
		"   POST  http://localhost:1234/path1",
		"   GET   http://localhost:1234/path2",
	}

	for _, line := range expectedLines {
		if !strings.Contains(output, line) {
			t.Errorf("Expected output to contain '%s', but it doesn't.\nActual output:\n%s", line, output)
		}
	}
}

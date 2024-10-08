package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/terrable-dev/terrable/config"
)

// Helper function to parse the test configuration
func parseTestConfig(t *testing.T) *config.TerrableConfig {
	testConfig := `
		module "simple_api" {
		  source = "terrable-dev/terrable-api/aws"
		  version = "0.0.1"
		  api_name = "simple-api"

		  handlers = {
		    HelloWorld: {
		        source = "./src/HelloWorld.ts"
		        http = {
					GET = "/"
		        }
		    },

		    HelloPost: {
		        source = "./src/HelloPost.ts"
		        http = {
					POST = "/hello-post"
		        }
		    }
		  }
		}
	`

	file, err := ParseHCL(testConfig)
	if err != nil {
		t.Fatalf("Failed to parse HCL: %v", err)
	}

	targetModule, err := FindTargetModule(file, "simple_api")
	if err != nil {
		t.Fatalf("Failed to find target module: %v", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	testFilePath := filepath.Join(cwd, "test.tf")
	terrableConfig, err := ParseModuleConfiguration(testFilePath, targetModule)
	if err != nil {
		t.Fatalf("Failed to parse module config: %v", err)
	}

	return terrableConfig
}

func TestAllHandlersMapped(t *testing.T) {
	config := parseTestConfig(t)
	expectedHandlers := []string{"HelloWorld", "HelloPost"}

	if len(config.Handlers) != len(expectedHandlers) {
		t.Errorf("Expected %d handlers, got %d", len(expectedHandlers), len(config.Handlers))
	}

	for _, expected := range expectedHandlers {
		found := false
		for _, handler := range config.Handlers {
			if handler.Name == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected handler %s not found", expected)
		}
	}
}

func TestSourceCorrectness(t *testing.T) {
	config := parseTestConfig(t)
	cwd, _ := os.Getwd()

	expectedSources := map[string]string{
		"HelloWorld": filepath.Join(cwd, "src", "HelloWorld.ts"),
		"HelloPost":  filepath.Join(cwd, "src", "HelloPost.ts"),
	}

	for _, handler := range config.Handlers {
		expected, ok := expectedSources[handler.Name]
		if !ok {
			t.Errorf("Unexpected handler: %s", handler.Name)
			continue
		}
		if handler.Source != expected {
			t.Errorf("Handler %s: expected source %s, got %s", handler.Name, expected, handler.Source)
		}
	}
}

func TestHttpConfiguration(t *testing.T) {
	config := parseTestConfig(t)

	expectedHttp := map[string]map[string]string{
		"HelloWorld": {
			"GET": "/",
		},
		"HelloPost": {
			"POST": "/hello-post",
		},
	}

	for _, handler := range config.Handlers {
		expected, ok := expectedHttp[handler.Name]
		if !ok {
			t.Errorf("Unexpected handler: %s", handler.Name)
			continue
		}

		for method, path := range expected {
			if handler.Http[method] != path {
				t.Errorf("Handler %s: expected HTTP %s path %s, got %s", handler.Name, method, path, handler.Http[method])
			}
		}
	}
}

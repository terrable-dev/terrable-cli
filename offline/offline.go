package offline

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"github.com/terrable-dev/terrable/config"
	"github.com/terrable-dev/terrable/utils"
)

func Run(filePath string, moduleName string, port string) error {
	terrableConfig, err := utils.ParseTerraformFile(filePath, moduleName)

	if err != nil {
		log.Fatalf("error running offline: %s", err)
	}

	// TODO: Validate config

	tomlConfig, err := config.ParseTerrableToml()

	if err != nil {
		panic(fmt.Errorf("error parsing .terrable.toml file: %w", err))
	}

	printConfig(*terrableConfig)

	var wg sync.WaitGroup
	defer wg.Done()

	r := mux.NewRouter()

	for _, handler := range terrableConfig.Handlers {
		go ServeHandler(&HandlerInstance{
			handlerConfig: handler,
			envVars:       tomlConfig.Environment,
		}, r)
	}

	fmt.Printf("Starting server on :%s\n", port)

	if err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%s", port), r); err != nil {
		return fmt.Errorf("could not start server on port %s. Error: %s", port, err.Error())
	}

	return nil
}

func printConfig(config config.TerrableConfig) {
	totalEndpoints := 0
	printlines := []string{}

	for _, handler := range config.Handlers {
		for method, path := range handler.Http {
			totalEndpoints += len(handler.Http)
			printlines = append(printlines, fmt.Sprintf("   %-*s http://localhost:8080%s\n", 5, method, path))
		}
	}

	fmt.Printf("Starting terrable local server... \n")
	fmt.Printf("%d Endpoint(s) to prepare...\n", totalEndpoints)
	fmt.Print(strings.Join(printlines, ""))
}

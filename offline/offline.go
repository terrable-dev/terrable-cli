package offline

import (
	"fmt"
	"log"
	"net"
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

	listener, activePort, err := getListener(port)

	if err != nil {
		return err
	}

	defer listener.Close()

	printConfig(*terrableConfig, activePort)

	var wg sync.WaitGroup
	defer wg.Wait()

	r := mux.NewRouter()

	// Start compiling and serving each handler
	for _, handler := range terrableConfig.Handlers {
		wg.Add(1)

		go func(handler config.HandlerMapping) {
			defer wg.Done()

			ServeHandler(&HandlerInstance{
				handlerConfig: handler,
				envVars:       tomlConfig.Environment,
			}, r)
		}(handler)
	}

	fmt.Printf("Starting server on :%d\n", activePort)

	server := &http.Server{
		Handler: r,
	}

	if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("could not start server on port %d. Error: %w", activePort, err)
	}

	return nil
}

func getListener(port string) (net.Listener, int, error) {
	var specificPortDesired bool = (port != "")

	if port == "" {
		port = "8080"
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%s", port))

	if err != nil {
		if specificPortDesired {
			return nil, 0, fmt.Errorf("could not start server on specified port %s: %w", port, err)
		}

		// Try with a random free-port
		listener, err = net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return nil, 0, fmt.Errorf("could not start server on any available port: %w", err)
		}
	}

	return listener, listener.Addr().(*net.TCPAddr).Port, nil
}

func printConfig(config config.TerrableConfig, port int) {
	totalEndpoints := 0
	printlines := []string{}

	for _, handler := range config.Handlers {
		for method, path := range handler.Http {
			totalEndpoints += 1
			printlines = append(printlines, fmt.Sprintf("   %-*s http://localhost:%d%s\n", 5, method, port, path))
		}
	}

	fmt.Printf("Starting terrable local server... \n")
	fmt.Printf("%d Endpoint(s) to prepare...\n", totalEndpoints)
	fmt.Print(strings.Join(printlines, ""))
}

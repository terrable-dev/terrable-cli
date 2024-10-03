package offline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/fsnotify/fsnotify"
	"github.com/terrable-dev/terrable/config"
)

type HandlerInstance struct {
	handlerConfig         config.HandlerMapping
	handlerTranspiledPath string
	readCodeMutex         sync.RWMutex
	recompileSyncLock     *sync.Once
	executionMutex        sync.Mutex
	envVars               map[string]interface{}
}

func (handlerInstance *HandlerInstance) GetExecutionPath() string {
	return handlerInstance.handlerTranspiledPath
}

func (handlerInstance *HandlerInstance) SetExecutionPath(path string) {
	handlerInstance.readCodeMutex.Lock()
	defer handlerInstance.readCodeMutex.Unlock()

	handlerInstance.handlerTranspiledPath = path
}

func (handlerInstance *HandlerInstance) CompileHandler() (inputFilePaths []string) {
	workingDirectory, _ := os.Getwd()
	fmt.Println("WD", workingDirectory)

	result := api.Build(api.BuildOptions{
		EntryPoints: []string{handlerInstance.handlerConfig.Source},
		Bundle:      true,
		Write:       true,
		Format:      api.FormatCommonJS,
		Platform:    api.PlatformNode,
		Target:      api.ES2015,
		Sourcemap:   api.SourceMapLinked,
		Metafile:    true,
		GlobalName:  "exports",
		Outdir:      filepath.Join(workingDirectory, ".terrable"),
	})

	if len(result.Errors) > 0 {
		printBuildErrors(result.Errors)
		return
	}

	handlerInstance.SetExecutionPath(filepath.ToSlash(result.OutputFiles[1].Path))
	return extractMetafileInputs(result.Metafile)
}

func extractMetafileInputs(metafileContents string) []string {
	var data Metafile

	err := json.Unmarshal([]byte(metafileContents), &data)

	if err != nil {
		println(fmt.Errorf("error parsing metafile: %w", err))
		return []string{}
	}

	var inputFiles []string

	for key := range data.Inputs {
		inputFiles = append(inputFiles, key)
	}

	return inputFiles
}

type Metafile struct {
	Inputs map[string]interface{} `json:"inputs"`
}

func printBuildErrors(result []api.Message) {
	fmt.Println("\n🚨 Build Errors:")
	fmt.Println(strings.Repeat("=", 50))

	for i, err := range result {
		fmt.Printf("Error #%d:\n", i+1)
		fmt.Printf("  File: %s\n", err.Location.File)
		fmt.Printf("  Line: %d, Column: %d\n", err.Location.Line, err.Location.Column)
		fmt.Printf("  Message: %s\n", err.Text)

		if err.Location.LineText != "" {
			fmt.Printf("  Code:\n")
			fmt.Printf("    %s\n", err.Location.LineText)
			fmt.Printf("    %s^\n", strings.Repeat(" ", err.Location.Column-1))
		}

		if i < len(result)-1 {
			fmt.Println(strings.Repeat("-", 50))
		}
	}

	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("\nTotal Errors: %d\n", len(result))
}

func (handlerInstance *HandlerInstance) WatchForChanges(inputFiles []string) {
	handlerInstance.recompileSyncLock = new(sync.Once)
	var wg sync.WaitGroup
	defer wg.Done()

	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()

	for _, file := range inputFiles {
		watcher.Add(file)
	}

	go func() {
		for event := range watcher.Events {
			if event.Op&fsnotify.Write == fsnotify.Write {
				handlerInstance.recompileSyncLock.Do(func() {
					handlerInstance.CompileHandler()
					handlerInstance.recompileSyncLock = new(sync.Once)
				})
			}
		}
	}()

	wg.Add(1)
	wg.Wait()
}

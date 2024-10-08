package main

import (
	"log"
	"os"

	"github.com/terrable-dev/terrable/config"
	"github.com/terrable-dev/terrable/offline"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "terrable",
		Version: buildInfo()["version"],

		Commands: []*cli.Command{
			{
				Name:  "offline",
				Usage: "",
				Action: func(cCtx *cli.Context) error {
					executablePath, _ := os.Getwd()
					tomlConfig, _ := config.ParseTerrableToml(executablePath)

					filePath := cCtx.String("file")
					moduleName := cCtx.String("module")

					if filePath == "" {
						filePath = tomlConfig.Offline.File
					}

					if moduleName == "" {
						moduleName = tomlConfig.Offline.Module
					}

					port := cCtx.String("port")
					err := offline.Run(filePath, moduleName, port)

					if err != nil {
						return err
					}

					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "file",
						Aliases:  []string{"f"},
						Required: false,
						Usage:    "Path to the Terraform file",
					},
					&cli.StringFlag{
						Name:     "module",
						Aliases:  []string{"m"},
						Required: false,
						Usage:    "Name of the terraform module to try and run locally",
					},
					&cli.StringFlag{
						Name:     "port",
						Aliases:  []string{"p"},
						Required: false,
						Value:    "8080",
						Usage:    "The port number that the local instance of the API should listen for requests at",
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

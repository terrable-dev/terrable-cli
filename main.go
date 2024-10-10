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
					tomlConfig, _ := config.ParseTerrableToml()

					filePath := cCtx.String("file")
					moduleName := cCtx.String("module")
					port := cCtx.String("port")

					if filePath == "" {
						filePath = tomlConfig.Offline.File
					}

					if moduleName == "" {
						moduleName = tomlConfig.Offline.Module
					}

					if port == "" {
						port = tomlConfig.Offline.Port
					}

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

package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"github.com/berrybyte-net/broom/service"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
)

// main is the application entrypoint.
func main() {
	app := &cli.App{
		Name:        "broom",
		Description: "scans JAR files to uncover the 29-09-2022 Minecraft malware infections",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "dir",
				Usage:   "the root directory to recursively scan",
				Aliases: []string{"d"},
				Value:   ".",
			},
		},
		Action: func(cCtx *cli.Context) error {
			foundInfected := false
			err := filepath.Walk(filepath.Clean(cCtx.String("dir")), func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if info.IsDir() {
					fmt.Printf("walking directory %s\n", path)
				} else {
					fmt.Printf("walking file %s\n", path)
				}

				if !info.IsDir() && filepath.Ext(path) == ".jar" { // found jar file
					fmt.Printf("found jar file %s\n", path)

					zf, err := zip.OpenReader(path)
					if err != nil {
						return err
					}

					for _, file := range zf.File {
						if !file.Mode().IsDir() {
							switch filepath.Base(file.Name) {
							case "plugin-config.bin":
								color.Red("found %s file in jar file %s", file.Name, path)
								foundInfected = true
							}
						}
					}

					if err := zf.Close(); err != nil {
						return err
					}
				}
				return nil
			})
			if err != nil {
				return err
			}

			foundService := false
			checker := service.GetChecker()
			if checker != nil {
				foundService, err = checker.HasService("vmd-gnu")
				if err != nil {
					color.Red(err.Error())
				}
			}

			fmt.Print("\n\n\n\n") // a bit of space before the assessment
			if foundInfected {
				color.Red("JAR files containing files that are a known part of the malware were found.")
			} else {
				color.Green("No files related to the malware were found!")
			}
			if foundService {
				color.Red("A system service that is a known part of the malware was found.")
			} else {
				color.Green("No services related to the malware were found!")
			}

			fmt.Print("Press Enter to exit...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

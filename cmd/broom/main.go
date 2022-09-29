package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc/mgr"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// SystemdSearchPaths are the base directories to search in for service files.
var SystemdSearchPaths = []string{
	// TODO: include user unit paths
	"/etc/systemd/system.control/",
	"/run/systemd/system.control/",
	"/run/systemd/transient/",
	"/run/systemd/generator.early/",
	"/etc/systemd/system/",
	"/etc/systemd/system.attached/",
	"/run/systemd/system/",
	"/run/systemd/system.attached/",
	"/run/systemd/generator/",
	"/usr/lib/systemd/system/",
	"/run/systemd/generator.late/",
}

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
			foundService := false
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

			if runtime.GOOS == "windows" {
				var s *uint16
				h, err := windows.OpenSCManager(s, nil, windows.SC_MANAGER_ENUMERATE_SERVICE)
				if err != nil {
					color.Red("could not query windows services (missing access?)")
				} else {
					svcMgr := &mgr.Mgr{Handle: h}
					services, err := svcMgr.ListServices()
					if err != nil {
						color.Red("could not query windows services")
					} else {
						for _, svc := range services {
							if svc == "vmd-gnu" {
								color.Red("found 'vmd-gnu' service")
								foundService = true
							}
						}
					}
				}
			} else {
				for _, dir := range SystemdSearchPaths {
					if info, err := os.Stat(dir); !os.IsNotExist(err) && info.IsDir() {
						err := filepath.Walk(dir, func(path string, info0 fs.FileInfo, err error) error {
							if err != nil {
								return err
							}
							if !info0.IsDir() && strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)) == "vmd-gnu" {
								color.Red("found %s service", path)
								foundService = true
							}

							return nil
						})
						if err != nil {
							color.Red("could not query services at path %s", dir)
						}
					}
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

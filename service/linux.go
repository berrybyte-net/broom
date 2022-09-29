//go:build linux || darwin

package service

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
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

func init() {
	instance = &LinuxChecker{}
}

type LinuxChecker struct {
}

func (lc *LinuxChecker) HasService(service string) (bool, error) {
	found := false
	for _, dir := range SystemdSearchPaths {
		if info, err := os.Stat(dir); !os.IsNotExist(err) && info.IsDir() {
			if found {
				break
			}
			err := filepath.Walk(dir, func(path string, info0 fs.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info0.IsDir() && strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)) == service {
					found = true
				}

				return nil
			})
			if err != nil {
				return found, errors.New(fmt.Sprintf("could not query services at path %s", dir))
			}
		}
	}
	return found, nil
}

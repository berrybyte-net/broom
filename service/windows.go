//go:build windows

package service

import (
	"errors"
	"fmt"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc/mgr"
)

func init() {
	instance = &WindowsChecker{}
}

type WindowsChecker struct {
}

func (wc *WindowsChecker) HasService(service string) (bool, error) {
	found := false
	var s *uint16
	// TODO: move to struct
	h, err := windows.OpenSCManager(s, nil, windows.SC_MANAGER_ENUMERATE_SERVICE)
	if err != nil {
		return false, errors.New(fmt.Sprintf("could not query services (%s)", err.Error()))
	} else {
		svcMgr := &mgr.Mgr{Handle: h}
		services, err := svcMgr.ListServices()
		if err != nil {
			return false, errors.New(fmt.Sprintf("could not query services (%s)", err.Error()))
		} else {
			for _, svc := range services {
				if svc == service {
					found = true
					break
				}
			}
		}
	}
	return found, nil
}

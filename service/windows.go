//go:build windows

package service

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc/mgr"
)

func init() {
	var s *uint16
	h, err := windows.OpenSCManager(s, nil, windows.SC_MANAGER_ENUMERATE_SERVICE)
	if err != nil {
		color.Red("could not query services (%s)", err.Error())
	} else {
		instance = &WindowsChecker{svcMgr: &mgr.Mgr{Handle: h}}
	}
}

type WindowsChecker struct {
	svcMgr *mgr.Mgr
}

func (wc *WindowsChecker) HasService(service string) (bool, error) {
	found := false
	services, err := wc.svcMgr.ListServices()
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
	return found, nil
}

// SPDX-FileCopyrightText: © 2014-2021 David Parsons
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
	"time"
)

func getState(s *mgr.Service) svc.State {
	status, err := s.Query()
	if err != nil {
		panic(fmt.Sprintf("Query(%s) failed: %s", s.Name, err))
	}
	return status.State
}

func waitState(s *mgr.Service, want svc.State) {
	for i := 0; ; i++ {
		have := getState(s)
		if have == want {
			return
		}
		if i > 10 {
			panic(fmt.Sprintf("%s state is=%d, waiting timeout", s.Name, have))
		}
		time.Sleep(300 * time.Millisecond)
	}
}

func startService(name string) {
	m, err := mgr.Connect()
	if err != nil {
		panic("SCM connection failed")
	}

	//goland:noinspection ALL
	defer m.Disconnect()

	s, err := m.OpenService(name)
	if err != nil {
		//println(fmt.Sprintf("Invalid service %s", name))
		return
	} else {
		println(fmt.Sprintf("Starting service %s", name))
	}

	//goland:noinspection ALL
	defer s.Close()

	if getState(s) == svc.Stopped {
		err = s.Start()
		if err != nil {
			panic(fmt.Sprintf("Control(%s) failed: %s", name, err))
		}
		waitState(s, svc.Running)
	}

	err = m.Disconnect()

}

func stopService(name string) {
	m, err := mgr.Connect()
	if err != nil {
		panic("SCM connection failed")
	}

	//goland:noinspection ALL
	defer m.Disconnect()

	s, err := m.OpenService(name)
	if err != nil {
		//println(fmt.Sprintf("Invalid service %s", name))
		return
	} else {
		println(fmt.Sprintf("Stopping service %s", name))
	}

	//goland:noinspection ALL
	defer s.Close()

	if getState(s) == svc.Running {
		_, err = s.Control(svc.Stop)
		if err != nil {
			panic(fmt.Sprintf("Control(%s) failed: %s", name, err))
		}
		waitState(s, svc.Stopped)
	}

	err = m.Disconnect()

}

func main() {
	println("Hello Windows!")
	stopService("VMAuthdService")
	stopService("VMwareHostd")
	stopService("VMUSBArbService")
	startService("VMAuthdService")
	startService("VMwareHostd")
	startService("VMUSBArbService")

}

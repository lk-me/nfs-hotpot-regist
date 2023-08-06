package main

import (
	"errors"
	hotpot "nfs-hotpot-regist/pkg/hotpot"
	"nfs-hotpot-regist/pkg/logger"
	"nfs-hotpot-regist/pkg/module/single"
	"sync"

	"github.com/godbus/dbus"
)

type Manager struct {
	conn           *dbus.Conn
	hotpot         *hotpot.Hotpot
	activated      bool
	allowedProcess string
	mu             sync.Mutex
}

func NewManager() (*Manager, error) {
	hot, err := hotpot.NewHotpot()
	if err != nil {
		logger.Fatal("Failed to new upgrade:", err)
		return nil, err
	}

	var m = &Manager{
		hotpot:         hot,
		activated:      false,
		allowedProcess: "dbus_example",
	}

	conn, err := dbus.SystemBus()
	if err != nil {
		logger.Fatal("Failed to connect dbus:", err)
		return nil, err
	}
	m.conn = conn

	return m, nil
}

func (m *Manager) BlackLists(version string, sender dbus.Sender) ([]string, *dbus.Error) {
	var list []string
	m.mu.Lock()

	if !single.SetSingleInstance() {
		return list, dbus.MakeFailedError(errors.New("process already exists"))
	}
	defer func() {
		single.Remove()
		m.mu.Unlock()
	}()
	list, err := m.hotpot.BlackLists()
	if err != nil {
		return list, dbus.MakeFailedError(err)
	}

	return list, nil
}

func (m *Manager) SaveRegistInfo(is bool, sender dbus.Sender) *dbus.Error {
	m.mu.Lock()
	if !single.SetSingleInstance() {
		return dbus.MakeFailedError(errors.New("process already exists"))
	}
	defer func() {
		single.Remove()
		m.mu.Unlock()
	}()
	b, err := m.isAllowed(sender)
	if !b || err != nil {
		return dbus.MakeFailedError(err)
	}
	err = m.hotpot.SaveRegistInfo(is)
	if err != nil {
		return dbus.MakeFailedError(err)
	}
	m.activated = is
	return nil
}

func (m *Manager) ReadRegistInfo(sender dbus.Sender) (bool, *dbus.Error) {
	m.mu.Lock()
	if !single.SetSingleInstance() {
		return false, dbus.MakeFailedError(errors.New("process already exists"))
	}
	defer func() {
		single.Remove()
		m.mu.Unlock()
	}()
	b, err := m.isAllowed(sender)
	if !b || err != nil {
		return false, dbus.MakeFailedError(err)
	}
	return m.activated, nil
}

func (m *Manager) Test(sender dbus.Sender) *dbus.Error {
	b, err := m.isAllowed(sender)
	if !b || err != nil {
		return dbus.MakeFailedError(err)
	}
	logger.Infof("%v\n", b)
	return nil
}

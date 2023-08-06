package main

import (
	"errors"
	"fmt"
	"nfs-hotpot-regist/pkg/logger"
	"os/exec"
	"strconv"
	"strings"

	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/godbus/dbus/prop"
)

const (
	dbusDest = "org.nfs.HotpotRegist1"
	dbusPath = "/org/nfs/HotpotRegist1"
	dbusIFC  = dbusDest
)

func (m *Manager) setupDBus() error {
	err := m.conn.Export(m, dbusPath, dbusIFC)
	if err != nil {
		return err
	}
	props := prop.New(m.conn, dbusPath, m.makeProps())
	node := &introspect.Node{
		Name: dbusDest,
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			prop.IntrospectData,
			{
				Name:       dbusIFC,
				Methods:    introspect.Methods(m),
				Properties: props.Introspection(dbusIFC),
			},
		},
	}
	err = m.conn.Export(introspect.NewIntrospectable(node), dbusPath,
		"org.freedesktop.DBus.Introspectable")
	if err != nil {
		return err
	}

	reply, err := m.conn.RequestName(dbusDest, dbus.NameFlagDoNotQueue)
	if err != nil {
		return err
	}

	if reply != dbus.RequestNameReplyPrimaryOwner {
		return fmt.Errorf("service %q has owned", dbusDest)
	}
	return nil
}

func (m *Manager) makeProps() map[string]map[string]*prop.Prop {
	return map[string]map[string]*prop.Prop{
		dbusIFC: {
			"Activated": &prop.Prop{
				Value:    &m.activated,
				Writable: false,
				Emit:     prop.EmitTrue,
				Callback: func(c *prop.Change) *dbus.Error {
					logger.Debugf("Running changed: %s -> %s", c.Name, c.Value)
					return nil
				},
			},
		},
	}
}

func (m *Manager) Wait() {
	select {}
}

func (m *Manager) isAllowed(sender dbus.Sender) (bool, error) {
	var pid uint32
	err := m.conn.BusObject().Call("org.freedesktop.DBus.GetConnectionUnixProcessID",
		0, string(sender)).Store(&pid)
	if err != nil {
		dbus.MakeFailedError(errors.New("process already exists"))
	}
	logger.Infof("%d\n", pid)
	// 执行 ps 命令获取进程信息
	cmd := exec.Command("ps", "-p", strconv.Itoa(int(pid)), "-o", "comm=,uid=")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("Failed to execute ps command: %v", err)
	}

	// 解析命令输出
	fields := strings.Fields(string(output))
	if len(fields) != 2 {
		return false, fmt.Errorf("Invalid ps output: %s", output)
	}

	// 第一个字段是进程名称，第二个字段是用户 ID
	name := fields[0]
	uid := fields[1]

	// 判断进程是否为 root 用户
	isRoot := uid == "0"

	logger.Infof("Process Name %s Is Root: %v\n", name, isRoot)
	if isRoot && name == m.allowedProcess {
		return true, nil
	}
	return false, errors.New("Not allowed to be called")
}

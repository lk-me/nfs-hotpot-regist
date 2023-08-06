package main

import (
	"nfs-hotpot-regist/pkg/logger"
	"os"
)

const (
	FAILED_PROCESS_EXISTS = -255
	FAILED_VERSION_EXISTS = -256
)

func main() {
	m, err := NewManager()
	logger.NewLogger("deepin-upgrade-manager", true)
	if err != nil {
		logger.Fatal("Failed to setup dbus:", err)
		os.Exit(-1)
	}
	logger.Info("start running dbus service")
	err = m.setupDBus()
	if err != nil {
		logger.Fatal("Failed to setup dbus:", err)
		os.Exit(-1)
	}
	defer func() {
		m.conn.Close()
	}()
	m.Wait()
	return
}

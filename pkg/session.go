package keychainlinker

import "github.com/godbus/dbus/v5"

type Session struct{}

func (s Session) Close() *dbus.Error {
	return nil
}

// opens

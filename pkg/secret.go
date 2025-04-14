package keychainlinker

import "github.com/godbus/dbus/v5"

type SecretStruct struct {
	Session     dbus.ObjectPath
	Parameters  []byte
	Value       []byte
	ContentType string
}

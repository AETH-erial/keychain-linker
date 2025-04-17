package main

import (
	"fmt"
	"os"

	keychainlinker "git.aetherial.dev/aeth/keychain-linker/pkg"
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
)

func main() {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	path := "/org/freedesktop/secrets"
	service := keychainlinker.NewService(dbus.ObjectPath(path))

	conn.Export(service, dbus.ObjectPath(path), "org.freedesktop.Secret.Service")
	conn.Export(service.Cache["/org/freedesktop/secrets/collection/default"], "/org/freedesktop/secrets/collection/default", "org.freedesktop.DBus.Properties")

	conn.Export(introspect.Introspectable(keychainlinker.DbusAdv), dbus.ObjectPath(path),
		"org.freedesktop.DBus.Introspectable")
	conn.Export(introspect.Introspectable(keychainlinker.DbusAdv), "/org/freedesktop/secrets/collection/default", "org.freedesktop.DBus.Introspectable")
	reply, err := conn.RequestName("org.freedesktop.secrets",
		dbus.NameFlagDoNotQueue)
	if err != nil {
		panic(err)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		fmt.Fprintln(os.Stderr, "name already taken")
		os.Exit(1)
	}
	fmt.Println("Listening on org.freedesktop.secrets / /org/freedesktop/secrets ...")
	select {}
}

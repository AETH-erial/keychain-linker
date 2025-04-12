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

	path := "/dev/aetherial/KeychainLinker"
	session := &keychainlinker.SecretService{SessionBase: "/dev/aetherial/KeychainLinker/session/",
		CollectionBase: "/dev/aetherial/KeychainLinker/collection/",
		Collections:    []dbus.ObjectPath{}}

	conn.Export(session, dbus.ObjectPath(path), "dev.aetherial.git.KeychainLinker.Service")
	conn.Export(introspect.Introspectable(keychainlinker.DbusAdv), dbus.ObjectPath(path),
		"org.freedesktop.DBus.Introspectable")

	reply, err := conn.RequestName("dev.aetherial.git.KeychainLinker.Service",
		dbus.NameFlagDoNotQueue)
	if err != nil {
		panic(err)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		fmt.Fprintln(os.Stderr, "name already taken")
		os.Exit(1)
	}
	fmt.Println("Listening on dev.aetherial.git.KeychainLinker.Service / /dev/aetherial/git/KeychainLinker/Service ...")
	select {}
}

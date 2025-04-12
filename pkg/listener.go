package keychainlinker

import (
	"fmt"
	"path"
	"strconv"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
)

const DbusAdv = `
<node>
	<interface name="dev.aetherial.git.KeychainLinker.Service">
		<method name="OpenSession">
		  <arg name="algorithm" direction="in" type="s"/>
		  <arg name="input" direction="in" type="v"/>
		  <arg name="output" direction="out" type="v"/>
		  <arg name="result" direction="out" type="o"/>
		</method>
		<method name="CreateCollection">
		  <arg name="properties" direction="in" type="a{sv}"/>
		  <arg name="alias" direction="in" type="s"/>
		  <arg name="collection" direction="out" type="o"/>
		  <arg name="prompt" direction="out" type="o"/>
		</method>
		<method name="SearchItems">
		  <arg name="attributes" direction="in" type="a{ss}"/>
		  <arg name="unlocked" direction="out" type="ao"/>
		  <arg name="locked" direction="out" type="ao"/>
		</method>
		<method name="Unlock">
		  <arg name="objects" direction="in" type="ao"/>
		  <arg name="unlocked" direction="out" type="ao"/>
		  <arg name="prompt" direction="out" type="o"/>
		</method>
		<method name="Lock">
		  <arg name="objects" direction="in" type="ao"/>
		  <arg name="locked" direction="out" type="ao"/>
		  <arg name="prompt" direction="out" type="o"/>
		</method>
		<method name="GetSecrets">
		  <arg name="items" direction="in" type="ao"/>
		  <arg name="session" direction="in" type="o"/>
		  <arg name="secrets" direction="out" type="a{o(ayays)}"/>
		</method>
		<method name="ReadAlias">
		  <arg name="name" direction="in" type="s"/>
		  <arg name="collection" direction="out" type="o"/>
		</method>
		<method name="SetAlias">
		  <arg name="name" direction="in" type="s"/>
		  <arg name="collection" direction="in" type="o"/>
		</method>
		<property name="Collections" type="ao" access="read"/>
	</interface>` + introspect.IntrospectDataString + `</node> `

type Session struct {
	Path string
	Open int
}

/*
Get the next session
*/
func (s *Session) next() int {
	s.Open = s.Open + 1
	return s.Open

}

func (s *Session) OpenSession(algorithm string, input dbus.Variant) (dbus.Variant, dbus.ObjectPath, *dbus.Error) {
	if algorithm != "PLAIN" {
		return dbus.Variant{}, dbus.ObjectPath(""), &dbus.ErrMsgInvalidArg
	}
	nextPath := path.Join(s.Path, strconv.Itoa(s.next()))
	fmt.Println("recieved algorithm: ", algorithm, "\nresponding with path: ", nextPath)
	return dbus.MakeVariant(algorithm), dbus.ObjectPath(nextPath), nil
}

type SecretStruct struct {
	Session     dbus.ObjectPath
	Parameters  []byte
	Value       []byte
	ContentType string
}

type SecretService struct {
	Collections    []dbus.ObjectPath
	SessionBase    string // e.g. "/org/freedesktop/secrets/session/"
	CollectionBase string // e.g. "/org/freedesktop/secrets/collection/"
}

type Collection struct {
}

// deletes the collection
func (c *Collection) Delete() (dbus.ObjectPath, *dbus.Error) {
	return dbus.ObjectPath("prompt"), nil
}

/*
Searches the collection for matching items

	:param attr: the attributes to attempt to match to a key in the collection
*/
func (c *Collection) SearchItems(attr map[string]string) ([]dbus.ObjectPath, *dbus.Error) {
	// implement a recursive searching thing
	return []dbus.ObjectPath{}, nil
}

/*
Creates a new item in the collection with the properties defined in 'props'.
Returns the items dbus object path, as well as a path to a dbus prompt

	:param props: a map of properties to assign to the item
	:param secret: the secret to encode into the collection
	:param replace: replace secret if a matching one is found in the store
*/
func (c *Collection) CreateItem(props map[string]dbus.Variant, secret SecretStruct, replace bool) (dbus.ObjectPath, dbus.ObjectPath) {
	return dbus.ObjectPath("/"), dbus.ObjectPath("/")
}

/*
Opens a session for the Secret Service Interface

	:param algorithm: the encryption algorithm to use with the client
	:param input: the data used when implementing more advanced encryption algos
*/
func (s *SecretService) OpenSession(algorithm string, input dbus.Variant) (dbus.Variant, dbus.ObjectPath, *dbus.Error) {
	if algorithm != "PLAIN" {
		return dbus.Variant{}, "/", dbus.MakeFailedError(fmt.Errorf("only PLAIN is supported"))
	}

	sessionPath := dbus.ObjectPath(path.Join(s.SessionBase, "1"))
	return input, sessionPath, nil
}

/*
Creates a collection with the Service object

	:param properties: a set of properties that are used by client apps
	:param alias: the shortname of the collection
*/
func (s *SecretService) CreateCollection(properties map[string]dbus.Variant, alias string) (dbus.ObjectPath, dbus.ObjectPath, *dbus.Error) {
	collPath := dbus.ObjectPath(path.Join(s.CollectionBase, "login"))
	s.Collections = append(s.Collections, collPath)
	return collPath, "/", nil
}

/*
search for items in the keychain that satisfy 'attrs'

	:param attrs: a map of search criteria
*/
func (s *SecretService) SearchItems(attrs map[string]string) ([]dbus.ObjectPath, []dbus.ObjectPath, *dbus.Error) {
	// Just return empty results for now
	return []dbus.ObjectPath{}, []dbus.ObjectPath{}, nil
}

/*
sets all dbus.Objects in 'objects' to the 'unlocked' position

	:param objects: a slice of dbus.Objects to unlock
*/
func (s *SecretService) Unlock(objects []dbus.ObjectPath) ([]dbus.ObjectPath, dbus.ObjectPath, *dbus.Error) {
	return objects, "/", nil // No prompt
}

/*
Sets all dbus.Objects in 'objects' to the 'locked' position

	:param objects: a slice of dbus.Objects to unlock
*/
func (s *SecretService) Lock(objects []dbus.ObjectPath) ([]dbus.ObjectPath, dbus.ObjectPath, *dbus.Error) {
	return objects, "/", nil // No prompt
}

/*
retrives secrets from an array of items/collections

	:param items: a slice of dbus.ObjectPath that will have their secrets returned
*/
func (s *SecretService) GetSecrets(items []dbus.ObjectPath, session dbus.ObjectPath) (map[dbus.ObjectPath]SecretStruct, *dbus.Error) {
	return map[dbus.ObjectPath]SecretStruct{}, nil
}

/*
returns the collection with the given alias 'name'

	:param name: the name of the alias to return
*/
func (s *SecretService) ReadAlias(name string) (dbus.ObjectPath, *dbus.Error) {
	return dbus.ObjectPath("/dev/aetherial/KeychainLinker/login"), nil
}

/*
sets the collections alias name to the specified value in 'name'

	:param name: the alias name to assign
	:param collection: the dbus.ObjectPath to assign the alias name to
*/
func (s *SecretService) SetAlias(name string, collection dbus.ObjectPath) *dbus.Error {
	// will implement later
	return nil
}

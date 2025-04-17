package keychainlinker

import (
	"fmt"
	"path"
	"strconv"
	"time"

	"github.com/godbus/dbus/v5"
)

const DEFAULT_COLLECTION = "/org/freedesktop/secrets/aliases/default"

type Cache struct {
	Collections map[dbus.ObjectPath]Collection
}

type Service struct {
	/*
		Working on implementing the org.freedesktop.Secret.Service interface, from their v0.2 spec:
		https://specifications.freedesktop.org/secret-service-spec/latest-single/#org.freedesktop.Secret.Service
	*/
	Cache          map[dbus.ObjectPath]*Collection
	Collections    []dbus.ObjectPath
	PathCount      int
	Alias          map[string]dbus.ObjectPath
	SessionBase    string // e.g. "/org/freedesktop/secrets/session/"
	CollectionBase string // e.g. "/org/freedesktop/secrets/collection/"
}

/*
implementing the properties interface
func (s *Service) Get(iface, property string) (dbus.Variant, *dbus.Error) {

		if iface == "org.freedesktop.Secret.Service" && property == "Collections" {
			return dbus.MakeVariant(s.Collections), nil
		}
		return dbus.Variant{}, dbus.MakeFailedError(fmt.Errorf("no such property"))
	}

	func (s *Service) GetAll(iface string) (map[string]dbus.Variant, *dbus.Error) {
		if iface == "org.freedesktop.Secret.Service" {
			return map[string]dbus.Variant{
				"Collections": dbus.MakeVariant(s.Collections),
			}, nil
		}
		return nil, dbus.MakeFailedError(fmt.Errorf("no such interface"))
	}

/*
Create a new service interface
*/
func NewService(base dbus.ObjectPath) *Service {
	return &Service{
		Cache: map[dbus.ObjectPath]*Collection{
			DEFAULT_COLLECTION: {
				Items: []dbus.ObjectPath{
					"/",
				},
				Label:     "default",
				Cache:     map[dbus.ObjectPath]CacheItem{},
				PathCount: 0,
				PathBase:  fmt.Sprintf("%s/collection/default", base),
				Private:   "true",
				Created:   uint64(time.Now().Unix()),
				Modified:  uint64(time.Now().Unix()),
			},
		},
		Collections: []dbus.ObjectPath{
			DEFAULT_COLLECTION,
		},
		PathCount: 0,
		Alias: map[string]dbus.ObjectPath{
			"default": DEFAULT_COLLECTION,
		},
		SessionBase:    fmt.Sprintf("%s/session", base),
		CollectionBase: fmt.Sprintf("%s/collection", base),
	}
}

/*
Opens a session for the Secret Service Interface

	:param algorithm: the encryption algorithm to use with the client
	:param input: the data used when implementing more advanced encryption algos
*/
func (s *Service) OpenSession(algorithm string, input dbus.Variant) (dbus.Variant, dbus.ObjectPath, *dbus.Error) {
	if algorithm != "PLAIN" {
		return dbus.Variant{}, "/", dbus.MakeFailedError(fmt.Errorf("only PLAIN is supported"))
	}
	sessionPath := dbus.ObjectPath(path.Join(s.SessionBase, "1"))
	fmt.Println(sessionPath)
	return input, sessionPath, nil
}

/*
Creates a collection with the Service object

	:param properties: a set of properties that are used by client apps
	:param alias: the shortname of the collection
*/
func (s *Service) CreateCollection(properties map[string]dbus.Variant, alias string) (dbus.ObjectPath, dbus.ObjectPath, *dbus.Error) {
	collPath := dbus.ObjectPath(path.Join(s.CollectionBase, strconv.Itoa(s.PathCount+1)))
	s.Collections = append(s.Collections, collPath)
	s.PathCount = s.PathCount + 1
	s.Alias[alias] = collPath
	return collPath, "/", nil
}

/*
search for items in the keychain that satisfy 'attrs'

	:param attrs: a map of search criteria
*/
func (s *Service) SearchItems(attrs map[string]string) ([]dbus.ObjectPath, []dbus.ObjectPath, *dbus.Error) {
	// Just return empty results for now
	found := []dbus.ObjectPath{}
	for _, v := range s.Cache {
		f, err := v.SearchItems(attrs)
		if err != nil {
			return []dbus.ObjectPath{}, []dbus.ObjectPath{}, err
		}
		found = append(found, f...)

	}
	return found, []dbus.ObjectPath{}, nil
}

/*
attempts to return secrets that were either already unlocked, or unlocked without a prompt, in addition to
a prompt path that can be used to unlock all remaining locked objects

	:param objects: a slice of dbus.Objects to unlock
*/
func (s *Service) Unlock(objects []dbus.ObjectPath) ([]dbus.ObjectPath, dbus.ObjectPath, *dbus.Error) {
	return []dbus.ObjectPath{}, dbus.ObjectPath("/"), nil // No prompt
}

/*
Sets all dbus.Objects in 'objects' to the 'locked' position

	:param objects: a slice of dbus.Objects to unlock
*/
func (s *Service) Lock(objects []dbus.ObjectPath) ([]dbus.ObjectPath, dbus.ObjectPath, *dbus.Error) {
	return []dbus.ObjectPath{}, dbus.ObjectPath("/"), nil // No prompt
}

/*
retrives secrets from an array of items/collections

	:param items: a slice of dbus.ObjectPath that will have their secrets returned
*/
func (s *Service) GetSecrets(items []dbus.ObjectPath, session dbus.ObjectPath) (map[dbus.ObjectPath]SecretStruct, *dbus.Error) {
	return map[dbus.ObjectPath]SecretStruct{}, nil
}

/*
Return a collection based on the alias name

	:param name: the alias to search for
*/
func (s *Service) ReadAlias(name string) (dbus.ObjectPath, *dbus.Error) {
	objectPath, ok := s.Alias[name]
	if !ok {
		return dbus.ObjectPath("/"), &dbus.ErrMsgNoObject
	}
	return objectPath, nil

}

/*
set the alias of the passed in collection

	:param name: the alias to set the collection to
	:param collection: the collection to modify
*/
func (s *Service) SetAlias(name string, collection dbus.ObjectPath) *dbus.Error {
	s.Alias[name] = collection
	return nil

}

package keychainlinker

import (
	"fmt"
	"strconv"

	"github.com/godbus/dbus/v5"
)

type CacheItem struct {
	Secret      SecretStruct
	Label       string
	LookupProps map[string]string
}

type Collection struct {
	/*
		Implying the org.freedesktop.Secret.Collection interface as per the v0.2 spec:
		https://specifications.freedesktop.org/secret-service-spec/latest-single/#org.freedesktop.Secret.Collection
	*/
	Items     []dbus.ObjectPath             // items in the collection
	Cache     map[dbus.ObjectPath]CacheItem // memory cache for the collection. Will be removed later
	PathCount int
	PathBase  string
	Private   string // specifies whether the collection is private or not
	Label     string //  The displayable label of this collection.
	Locked    string //  Whether the collection is locked and must be authenticated by the client application.
	Created   uint64 //  The unix time when the collection was created.
	Modified  uint64 //  The unix time when the collection was last modified.
}

func (c *Collection) Get(iface, property string) (dbus.Variant, *dbus.Error) {
	if iface != "org.freedesktop.Secret.Collection" {
		return dbus.Variant{}, dbus.MakeFailedError(fmt.Errorf("no such property"))
	}
	switch property {
	case "Label":
		return dbus.MakeVariant(c.Label), nil
	case "Locked":
		return dbus.MakeVariant(c.Locked), nil
	case "Created":
		return dbus.MakeVariant(c.Created), nil
	case "Modified":
		return dbus.MakeVariant(c.Modified), nil
	case "Items":
		return dbus.MakeVariant(c.Items), nil
	}
	return dbus.Variant{}, dbus.MakeFailedError(fmt.Errorf("no such property"))

}

func (c *Collection) Set(iface, property string, value dbus.Variant) *dbus.Error {
	if iface == "org.freedesktop.Secret.Collection" && property == "Label" {
		if label, ok := value.Value().(string); ok {
			c.Label = label
			return nil
		}
		return dbus.MakeFailedError(fmt.Errorf("invalid type"))
	}
	return dbus.MakeFailedError(fmt.Errorf("no such property"))
}

func (c *Collection) GetAll(iface string) (map[string]dbus.Variant, *dbus.Error) {
	if iface == "org.freedesktop.Secret.Collection" {
		return map[string]dbus.Variant{
			"Label":    dbus.MakeVariant(c.Label),
			"Locked":   dbus.MakeVariant(c.Locked),
			"Created":  dbus.MakeVariant(c.Created),
			"Modified": dbus.MakeVariant(c.Modified),
			"Items":    dbus.MakeVariant(c.Items),
		}, nil
	}
	return nil, dbus.MakeFailedError(fmt.Errorf("no such interface"))
}

/*
Create a path to assign the secret
*/
func (c *Collection) pathBuilder() dbus.ObjectPath {
	return dbus.ObjectPath(c.PathBase + "/" + strconv.Itoa(c.PathCount+1))

}

// deletes the collection, returning an object path tied to a prompt incase it is necessary.
func (c *Collection) Delete() (dbus.ObjectPath, *dbus.Error) {
	return dbus.ObjectPath("prompt"), nil
}

/*
Searches the collection for matching items

	:param attr: the attributes to attempt to match to a key in the collection
*/
func (c *Collection) SearchItems(attr map[string]string) ([]dbus.ObjectPath, *dbus.Error) {
	matched := []dbus.ObjectPath{}
	for path, sec := range c.Cache {
		secretAttr := sec.LookupProps
		var passedLookup bool
		passedLookup = true
		for k, v := range attr {
			got, ok := secretAttr[k]
			if !ok {
				passedLookup = false
				continue
			}
			if got != v {
				passedLookup = false
				continue
			}
		}
		if passedLookup {
			matched = append(matched, path)
		}

	}

	return matched, nil
}

/*
Creates a new item in the collection with the properties defined in 'props'.
Returns the items dbus object path, as well as a path to a dbus prompt incase it is required to edit

	:param props: a map of properties to assign to the item. Will be used to match during lookups
	:param secret: the secret to encode into the collection
	:param replace: replace secret if a matching one is found in the store
*/
func (c *Collection) CreateItem(props map[string]dbus.Variant, secret SecretStruct, replace bool) (dbus.ObjectPath, dbus.ObjectPath, *dbus.Error) {
	v := props["org.freedesktop.Secret.Item.Attributes"]
	attrs, ok := v.Value().(map[string]string)
	if !ok {
		return dbus.ObjectPath("/"), dbus.ObjectPath("/"), &dbus.ErrMsgNoObject
	}
	label, ok := props["org.freedesktop.Secret.Item.Label"].Value().(string)
	if !ok {
		// no label found
		label = ""
	}
	if !replace {
		// implement the the replace option
	}
	path := c.pathBuilder()
	c.Cache[path] = CacheItem{
		LookupProps: attrs,
		Label:       label,
		Secret:      secret,
	}

	return path, dbus.ObjectPath("/"), nil
}

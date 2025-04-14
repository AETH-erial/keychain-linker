package keychainlinker

import "github.com/godbus/dbus/v5"

type Collection struct {
	/*
		Implying the org.freedesktop.Secret.Collection interface as per the v0.2 spec:
		https://specifications.freedesktop.org/secret-service-spec/latest-single/#org.freedesktop.Secret.Collection
	*/
	Items    []dbus.ObjectPath // items in the collection
	Private  string            // specifies whether the collection is private or not
	Label    string            //  The displayable label of this collection.
	Locked   string            //  Whether the collection is locked and must be authenticated by the client application.
	Created  uint64            //  The unix time when the collection was created.
	Modified uint64            //  The unix time when the collection was last modified.
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
	// implement a recursive searching thing
	return []dbus.ObjectPath{}, nil
}

/*
Creates a new item in the collection with the properties defined in 'props'.
Returns the items dbus object path, as well as a path to a dbus prompt incase it is required to edit

	:param fields: a map of properties to assign to the item. Will be used to match during lookups
	:param secret: the secret to encode into the collection
	:param replace: replace secret if a matching one is found in the store
*/
func (c *Collection) CreateItem(fields map[string]dbus.Variant, secret SecretStruct, replace bool) (dbus.ObjectPath, dbus.ObjectPath, *dbus.Error) {
	return dbus.ObjectPath("/"), dbus.ObjectPath("/"), nil
}

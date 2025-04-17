import secretstorage
import dbus

# Connect to the session bus
bus = dbus.SessionBus()

# Connect to the Secret Service
connection = secretstorage.dbus_init()
collection = secretstorage.get_default_collection(connection)


# Define attributes and label
attributes = {
    'username': 'alice',
    'service': 'example-app'
}
label = 'My Example Secret'
secret = 'super_secret_password_123'

# Create the item
item = collection.create_item(
    label=label,
    attributes=attributes,
    secret=secret,
    replace=True  # Replace if an item with same attributes exists
)

print(f"Secret stored with label: {label}")


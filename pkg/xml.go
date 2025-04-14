package keychainlinker

import (
	"github.com/godbus/dbus/v5/introspect"
)

const DbusAdv = `
<node>
	<interface name="dev.aetherial.git.KeychainLinker.Service">
		<method name="OpenSession">
		  <arg name="algorithm" type="s" direction="in"/>
		  <arg name="input" type="v" direction="in"/>
		  <arg name="output" type="v" direction="out"/>
		  <arg name="session" type="o" direction="out"/>
		</method>
		<method name="SearchItems">
		  <arg name="attributes" type="a{ss}" direction="in"/>
		  <arg name="unlocked" type="ao" direction="out"/>
		  <arg name="locked" type="ao" direction="out"/>
		</method>
		<method name="CreateCollection">
		  <arg name="properties" type="a{sv}" direction="in"/>
		  <arg name="alias" type="s" direction="in"/>
		  <arg name="collection" type="o" direction="out"/>
		  <arg name="prompt" type="o" direction="out"/>
		</method>
		<method name="Unlock">
		  <arg name="objects" type="ao" direction="in"/>
		  <arg name="unlocked" type="ao" direction="out"/>
		  <arg name="prompt" type="o" direction="out"/>
		</method>
		<method name="Lock">
		  <arg name="objects" type="ao" direction="in"/>
		  <arg name="locked" type="ao" direction="out"/>
		  <arg name="prompt" type="o" direction="out"/>
		</method>
		<method name="GetSecrets">
		  <arg name="items" type="ao" direction="in"/>
		  <arg name="session" type="o" direction="in"/>
		  <arg name="secrets" type="a{o((oayays))}" direction="out"/>
		</method>
		<method name="ReadAlias">
		  <arg name="name" type="s" direction="in"/>
		  <arg name="collection" type="o" direction="out"/>
		</method>
		<method name="SetAlias">
		  <arg name="name" type="s" direction="in"/>
		  <arg name="collection" type="o" direction="in"/>
		</method>
		<property name="Collections" type="ao" access="read"/>	
	</interface>` + introspect.IntrospectDataString + `</node> `

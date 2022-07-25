package protocol

import "strings"

/*
The JID or Jabber Identifier refers to a specific user with XMPP, similar to an email address.
	ie. username@mydomain.com/q3er3
*/
type JID struct {
	Username   string // The Username is the way name the user registered in the server.
	Domain     string // The Domain points to the server where the user registered.
	DeviceName string // The DeviceName identifies a device where the user is logged in.
}

// JIDFromString returns a JID struct based on a string with the format username@domain/device. The device is optional.
// Returns JID and ok.
func JIDFromString(str string) (JID, bool) {
	split := strings.SplitN(str, "@", 2)
	if len(split) != 2 {
		return JID{}, false
	}

	jid := JID{}

	end := strings.SplitN(split[1], "/", 2)
	if len(end) == 2 {
		jid.DeviceName = end[1]
	}
	jid.Domain = end[0]
	jid.Username = split[0]

	return jid, true
}

// Returns the complete JID. ie user@domain.com/device
func (jid *JID) String() string {
	return jid.Username + "@" + jid.Domain + "/" + jid.DeviceName
}

// BaseJid returns the JID in string format without the device. ie user@domain.com
func (jid *JID) BaseJid() string {
	return jid.Username + "@" + jid.Domain
}

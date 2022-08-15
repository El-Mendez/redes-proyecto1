package query

// Query is just an interface to help my IDE and avoid using interface{} everywhere in an IQ. Similar to the Stanza interface.
type Query interface {
	isQuery()
}

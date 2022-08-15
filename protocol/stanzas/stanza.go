package stanzas

// Stanza is a simple interface to avoid using interface{} everywhere. It mainly helps readability and improves
// IDE performance.
type Stanza interface {
	isStanza()
}

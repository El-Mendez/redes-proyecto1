package query

type Query interface {
	isQuery() bool
}

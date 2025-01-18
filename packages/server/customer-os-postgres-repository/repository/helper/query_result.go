package helper

// Deprecated
// use postgres_entity as return type
type QueryResult struct {
	Result interface{}
	Error  error
}

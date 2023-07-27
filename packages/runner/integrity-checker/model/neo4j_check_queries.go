package model

// Query contains a single query definition
type Query struct {
	Name    string `json:"name"`
	Query   string `json:"query"`
	Message string `json:"message"`
}

// QueryList is the top-level structure
type IntegrityCheckQueries struct {
	Queries []Query `json:"queries"`
}

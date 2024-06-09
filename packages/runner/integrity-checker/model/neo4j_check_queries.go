package model

// Query contains a single query definition
type Query struct {
	Name    string `json:"name" required:"true"`
	Query   string `json:"query" required:"true"`
	Message string `json:"message"`
}

// Group is a collection of queries for a specific domain
type Group struct {
	Name    string  `json:"name"`
	Queries []Query `json:"queries"`
}

// QueryList is the top-level structure
type IntegrityCheckQueries struct {
	Groups  []Group `json:"groups"`
	Queries []Query `json:"queries"`
}

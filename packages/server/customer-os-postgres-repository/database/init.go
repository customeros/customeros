package database

func InitDatabase(dbConfig DatabaseConfig) (*DbConnections, error) {
	db, err := NewConnection(&dbConfig)
	return db, err
}

package config

type Config struct {
	ApiPort string `env:"PORT"`
	Neo4j   struct {
		Target string `env:"NEO4J_TARGET,required"`
		User   string `env:"NEO4J_AUTH_USER,required,unset"`
		Pwd    string `env:"NEO4J_AUTH_PWD,required,unset"`
		Realm  string `env:"NEO4J_AUTH_REALM"`
	}
}

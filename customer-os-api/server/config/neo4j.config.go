package config

//func ReadConfig(path string) (*Config, error) {
//	file, err := ioutil.ReadFile(path)
//	if err != nil {
//		return nil, err
//	}
//	config := Config{}
//	if err = json.Unmarshal(file, &config); err != nil {
//		return nil, err
//	}
//	return &config, nil
//}
//
//// end::readConfig[]
//
//type Config struct {
//	Uri      string `json:"neo4j://neo4j-customer-os.openline-development.svc.cluster.local:7687"`
//	Username string `json:"neo4j"`
//	Password string `json:"StrongLocalPa$$"`
//}
//
///**
// * Initiate the Neo4j Driver
// *
// * @param {Config} config   Config struct loaded from config.json
// * @returns {neo4j.Driver}	A new Driver instance
// */
//// tag::initDriver[]
//func NewDriver(settings *Config) (neo4j.Driver, error) {
//	return nil, nil
//}

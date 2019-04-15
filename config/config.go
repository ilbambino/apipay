// package config holds the const of the configuration values
// Right now it is just a thin wrapper around viper. If growing in
// use or complexity it should probably have its own structs and
// hide viper
// also depending on the use, maybe more sources of config
// TODO add tests

package config

import "github.com/spf13/viper"

const (
	// MongoHost holds the mongo server name
	MongoHost = "MongoHost"

	// MongoPort holds the mongo server name
	MongoPort = "MongoPort"

	// MongoUser holds the mongo server url
	MongoUser = "MongoUser"

	// MongoPassword holds the mongo server url
	MongoPassword = "MongoPassword"
)

// Load loads the config from the env vars. It could be extended to load also
// from a different source
// it also sets the defaults to the values if needed
func Load() error {

	viper.SetEnvPrefix("APIPAY")

	viper.SetDefault(MongoHost, "localhost")
	err := viper.BindEnv(MongoHost)
	if err != nil {
		return err
	}

	viper.SetDefault(MongoPort, 27017)
	err = viper.BindEnv(MongoPort)
	if err != nil {
		return err
	}

	err = viper.BindEnv(MongoUser)
	if err != nil {
		return err
	}
	err = viper.BindEnv(MongoPassword)
	if err != nil {
		return err
	}

	return nil

}

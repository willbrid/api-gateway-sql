package config

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Auth struct {
	Enabled  bool   `mapstructure:"enabled"`
	Username string `mapstructure:"username" validate:"required_if=Enabled true,min=2,max=25"`
	Password string `mapstructure:"password" validate:"required_if=Enabled true,min=8"`
}

type Database struct {
	Name     string        `mapstructure:"name" validate:"required,max=25"`
	Type     string        `mapstructure:"type" validate:"required,oneof=mariadb mysql postgres sqlserver sqlite"`
	Host     string        `mapstructure:"host" validate:"required_unless=Type sqlite,omitempty,ipv4"`
	Port     int           `mapstructure:"port" validate:"required_unless=Type sqlite,omitempty,min=1024,max=49151"`
	Username string        `mapstructure:"username" validate:"required_unless=Type sqlite"`
	Password string        `mapstructure:"password" validate:"required_unless=Type sqlite"`
	Dbname   string        `mapstructure:"dbname" validate:"required"`
	Sslmode  bool          `mapstructure:"sslmode"`
	Timeout  time.Duration `mapstructure:"timeout" validate:"required"`
}

type Target struct {
	Name           string `mapstructure:"name" validate:"required,max=25"`
	DataSourceName string `mapstructure:"data_source_name" validate:"required"`
	Multi          bool   `mapstructure:"multi"`
	BatchSize      int    `mapstructure:"batch_size" validate:"required_if=Multi true"`
	BufferSize     int    `mapstructure:"buffer_size" validate:"required_if=Multi true"`
	BatchFields    string `mapstructure:"batch_fields" validate:"required_if=Multi true"`
	SqlQuery       string `mapstructure:"sql" validate:"required"`
}

type Config struct {
	ApiGatewaySQL struct {
		EnableSwagger bool   `mapstructure:"enable_swagger"`
		Sqlitedb      string `mapstructure:"sqlitedb" validate:"required"`
		Auth          `mapstructure:"auth"`
		Databases     []Database `mapstructure:"databases" validate:"gt=0,required,dive"`
		Targets       []Target   `mapstructure:"targets" validate:"gt=0,required,dive"`
	} `mapstructure:"api_gateway_sql"`
}

// setConfigDefaults used to set default configuration
func setConfigDefaults(v *viper.Viper) {
	v.SetDefault("api_gateway_sql.enable_swagger", true)
	v.SetDefault("api_gateway_sql.sqlitedb", "/data/api_gateway_sql")
	v.SetDefault("api_gateway_sql.auth.enabled", false)
	v.SetDefault("api_gateway_sql.auth.username", "")
	v.SetDefault("api_gateway_sql.auth.password", "")
	v.SetDefault("api_gateway_sql.databases", make([]Database, 0))
	v.SetDefault("api_gateway_sql.targets", make([]Target, 0))
}

// ReadConfigFile reads configuration file and return viper instance
func ReadConfigFile(filename string) (*viper.Viper, error) {
	viperInstance := viper.New()

	// Load configuration file
	viperInstance.SetConfigType("yaml")
	viperInstance.SetConfigFile(filename)

	if err := viperInstance.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("configuration file '%s' not found", filename)
		} else {
			return nil, err
		}
	}

	return viperInstance, nil
}

// LoadConfig used to load configuration file
func LoadConfig(viperInstance *viper.Viper, validate *validator.Validate) (*Config, error) {

	setConfigDefaults(viperInstance)

	// Parse configuration file to Config struct
	var config Config
	if err := viperInstance.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Validate config struct
	if err := validate.Struct(config); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return nil, err
		}

		for _, err := range err.(validator.ValidationErrors) {
			return nil, fmt.Errorf("validation failed on field '%s' for condition '%s'", err.Field(), err.Tag())
		}
	}

	return &config, nil
}

// GetTargetByName is a method of Config struct for retreive target by his name
func (config *Config) GetTargetByName(targetName string) (Target, bool) {
	var target Target
	found := false

	for _, targetItem := range config.ApiGatewaySQL.Targets {
		if targetItem.Name == targetName {
			found = true
			target = targetItem
			break
		}
	}

	return target, found
}

// GetDatabaseByDataSourceName is a method of Config struct for retreive datasource by his name
func (config *Config) GetDatabaseByDataSourceName(dataSourceName string) (Database, bool) {
	var database Database
	found := false

	for _, databaseItem := range config.ApiGatewaySQL.Databases {
		if databaseItem.Name == dataSourceName {
			found = true
			database = databaseItem
			break
		}
	}

	return database, found
}

package config_test

import (
	"github.com/willbrid/api-gateway-sql/config"

	"bytes"
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func triggerTest(t *testing.T, yamlConfig []byte, expectations []string, index int) {
	v := viper.New()
	v.SetConfigType("yaml")

	if err := v.ReadConfig(bytes.NewBuffer([]byte(yamlConfig))); err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())
	_, err := config.LoadConfig(v, validate)

	expected := expectations[index]

	if err == nil {
		t.Errorf("no error returned, expected:\n%v", expected)
	}

	if err.Error() != expected {
		t.Errorf("\nexpected:\n%v\ngot:\n%v", expected, err.Error())
	}
}

func TestReadConfigFile_ReturnFileNotFoundError(t *testing.T) {
	t.Parallel()

	var filename string

	_, err := config.ReadConfigFile(filename)
	expected := "configuration file '' not found"

	if err == nil {
		t.Fatalf("no error returned, expected:\n%v", expected)
	}

	if err.Error() != expected {
		t.Errorf("\nexpected:\n%v\ngot:\n%v", expected, err.Error())
	}
}

func TestReadConfigFile_ReturnFileNotExistError(t *testing.T) {
	t.Parallel()

	var filename string = "nonexistentfile.yaml"

	_, err := config.ReadConfigFile(filename)

	expected := "open nonexistentfile.yaml: no such file or directory"

	if err == nil {
		t.Fatalf("no error returned, expected:\n%v", expected)
	}

	if err.Error() != expected {
		t.Errorf("\nexpected:\n%v\ngot:\n%v", expected, err.Error())
	}
}

func TestLoadConfig_ReturnErrorWithBadAuthFieldWhenAuthEnabled(t *testing.T) {
	t.Parallel()

	configSlices := [][]byte{
		[]byte(`---
api_gateway_sql:
  auth:
    enabled: true
    username: ""
`),
		[]byte(`---
api_gateway_sql:
  auth:
    enabled: true
    username: "x"
`),
		[]byte(`---
api_gateway_sql:
  auth:
    enabled: true
    username: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
`),
		[]byte(`---
api_gateway_sql:
  auth:
    enabled: true
    username: "xxxxx"
    password: ""
`),
		[]byte(`---
api_gateway_sql:
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxx
`),
	}

	expectations := []string{
		"validation failed on field 'Username' for condition 'required_if'",
		"validation failed on field 'Username' for condition 'min'",
		"validation failed on field 'Username' for condition 'max'",
		"validation failed on field 'Password' for condition 'required_if'",
		"validation failed on field 'Password' for condition 'min'",
	}

	for index, yamlConfig := range configSlices {
		t.Run(fmt.Sprintf("LoadConfig  #%v", index), func(subT *testing.T) {
			triggerTest(subT, yamlConfig, expectations, index)
		})
	}
}

func TestLoadConfig_ReturnErrorWithBadSqlitedbField(t *testing.T) {
	t.Parallel()

	yamlConfig := []byte(`---
api_gateway_sql:
  sqlitedb: ''
`)

	v := viper.New()
	v.SetConfigType("yaml")

	if err := v.ReadConfig(bytes.NewBuffer([]byte(yamlConfig))); err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())
	_, err := config.LoadConfig(v, validate)

	if err == nil {
		t.Errorf("no error returned")
	}

	if err.Error() == "" {
		t.Errorf("no error message found")
	}
}

func TestLoadConfig_ReturnErrorWithBadDabatasesField(t *testing.T) {
	t.Parallel()

	configSlices := [][]byte{
		[]byte(`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
`),
		[]byte(`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
`),
		[]byte(`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: ""
`),
		[]byte(`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
`),
		[]byte(`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: ""
`),
		[]byte(`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "xxxxx"
`),
		[]byte(`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "mariadb"
    host: ""
`),
		[]byte(`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "mariadb"
    host: "127.0"
`),
		[]byte(`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "mariadb"
    host: "127.0.0.1"
    port: ""
`),
		[]byte(`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "mariadb"
    host: "127.0.0.1"
    port: "1000"
`),
		[]byte(`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "mariadb"
    host: "127.0.0.1"
    port: "49152"
`),
		[]byte(`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "mariadb"
    host: "127.0.0.1"
    port: "3306"
    username: ""
`),
		[]byte(`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "mariadb"
    host: "127.0.0.1"
    port: "3306"
    username: "xxxxx"
    password: ""
`),
		[]byte(`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "mariadb"
    host: "127.0.0.1"
    port: "3306"
    username: "xxxxx"
    password: "xxxxx"
    dbname: ""
`),
		[]byte(`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "mariadb"
    host: "127.0.0.1"
    port: "3306"
    username: "xxxxx"
    password: "xxxxx"
    dbname: "xxxxx"
`),
		[]byte(`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "mariadb"
    host: "127.0.0.1"
    port: "3306"
    username: "xxxxx"
    password: "xxxxx"
    dbname: "xxxxx"
    timeout: "10"
`),
	}

	expectations := []string{
		"validation failed on field 'Databases' for condition 'gt'",
		"validation failed on field 'Databases' for condition 'gt'",
		"validation failed on field 'Name' for condition 'required'",
		"validation failed on field 'Name' for condition 'max'",
		"validation failed on field 'Type' for condition 'required'",
		"validation failed on field 'Type' for condition 'oneof'",
		"validation failed on field 'Host' for condition 'required_unless'",
		"validation failed on field 'Host' for condition 'ipv4'",
		"validation failed on field 'Port' for condition 'required_unless'",
		"validation failed on field 'Port' for condition 'min'",
		"validation failed on field 'Port' for condition 'max'",
		"validation failed on field 'Username' for condition 'required_unless'",
		"validation failed on field 'Password' for condition 'required_unless'",
		"validation failed on field 'Dbname' for condition 'required'",
		"validation failed on field 'Timeout' for condition 'required'",
		"decoding failed due to the following error(s):\n\n'api_gateway_sql.databases[0].timeout' time: missing unit in duration",
	}

	for index, yamlConfig := range configSlices {
		t.Run(fmt.Sprintf("LoadConfig  #%v", index), func(subT *testing.T) {
			triggerTest(subT, yamlConfig, expectations, index)
		})
	}
}

func TestLoadConfig_ReturnErrorWithBadTargetsField(t *testing.T) {
	t.Parallel()

	configSlices := [][]byte{
		[]byte(`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "mariadb"
    host: "127.0.0.1"
    port: 3306
    username: "xxxxx"
    password: "xxxxx"
    dbname: "xxxxx"
    sslmode: false
    timeout: "10s"
`),
		[]byte(`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "mariadb"
    host: "127.0.0.1"
    port: 3306
    username: "xxxxx"
    password: "xxxxx"
    dbname: "xxxxx"
    sslmode: false
    timeout: "10s"
  targets:
`),
		[]byte(`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "mariadb"
    host: "127.0.0.1"
    port: 3306
    username: "xxxxx"
    password: "xxxxx"
    dbname: "xxxxx"
    sslmode: false
    timeout: "10s"
  targets:
  - name: ""
`),
		[]byte(`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "mariadb"
    host: "127.0.0.1"
    port: 3306
    username: "xxxxx"
    password: "xxxxx"
    dbname: "xxxxx"
    sslmode: false
    timeout: "10s"
  targets:
  - name: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
`),
		[]byte(`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "mariadb"
    host: "127.0.0.1"
    port: 3306
    username: "xxxxx"
    password: "xxxxx"
    dbname: "xxxxx"
    sslmode: false
    timeout: "10s"
  targets:
  - name: "xxxxx"
    data_source_name: ""
`),
		[]byte(`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "mariadb"
    host: "127.0.0.1"
    port: 3306
    username: "xxxxx"
    password: "xxxxx"
    dbname: "xxxxx"
    sslmode: false
    timeout: "10s"
  targets:
  - name: "xxxxx"
    data_source_name: "xxxxx"
    sql: ""
`),
		[]byte(`---
api_gateway_sql:
  sqlitedb: "/data/api_gateway_sql"
  auth:
    enabled: true
    username: "xxxxx"
    password: xxxxxxxx
  databases:
  - name: "xxxxx"
    type: "mariadb"
    host: "127.0.0.1"
    port: 3306
    username: "xxxxx"
    password: "xxxxx"
    dbname: "xxxxx"
    sslmode: false
    timeout: "10s"
  targets:
  - name: "xxxxx"
    data_source_name: "xxxxx"
    sql: "select * from student"
    Multi: true
`),
	}

	expectations := []string{
		"validation failed on field 'Targets' for condition 'gt'",
		"validation failed on field 'Targets' for condition 'gt'",
		"validation failed on field 'Name' for condition 'required'",
		"validation failed on field 'Name' for condition 'max'",
		"validation failed on field 'DataSourceName' for condition 'required'",
		"validation failed on field 'SqlQuery' for condition 'required'",
		"validation failed on field 'BatchSize' for condition 'required_if'",
	}

	for index, yamlConfig := range configSlices {
		t.Run(fmt.Sprintf("LoadConfig  #%v", index), func(subT *testing.T) {
			triggerTest(subT, yamlConfig, expectations, index)
		})
	}
}

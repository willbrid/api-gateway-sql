package config_test

import (
	"strings"

	"github.com/willbrid/api-gateway-sql/config"

	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func triggerTest(t *testing.T, yamlConfig []byte) {
	v := viper.New()
	v.SetConfigType("yaml")

	if err := v.ReadConfig(bytes.NewBuffer([]byte(yamlConfig))); err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	_, err := config.LoadConfig(v, validate)

	var fieldErr validator.FieldError
	if !errors.As(err, &fieldErr) && !strings.Contains(err.Error(), "unable to unmarshal config struct") {
		t.Errorf("wrong error: %v", err)
	}
}

func TestReadConfigFile_ReturnFileNotFoundError(t *testing.T) {
	t.Parallel()

	var filename string

	_, err := config.ReadConfigFile(filename)

	if err == nil {
		t.Error("no error returned for file not found")
	}
}

func TestReadConfigFile_ReturnFileNotExistError(t *testing.T) {
	t.Parallel()

	filename := "nonexistentfile.yaml"
	_, err := config.ReadConfigFile(filename)

	if err == nil {
		t.Error("no error returned for file no exist")
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

	for index, yamlConfig := range configSlices {
		t.Run(fmt.Sprintf("LoadConfig  #%v", index), func(subT *testing.T) {
			triggerTest(subT, yamlConfig)
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

	validate := validator.New(validator.WithRequiredStructEnabled())
	_, err := config.LoadConfig(v, validate)

	if err == nil {
		t.Errorf("no error returned for bad sqlite field")
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

	for index, yamlConfig := range configSlices {
		t.Run(fmt.Sprintf("LoadConfig  #%v", index), func(subT *testing.T) {
			triggerTest(subT, yamlConfig)
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

	for index, yamlConfig := range configSlices {
		t.Run(fmt.Sprintf("LoadConfig  #%v", index), func(subT *testing.T) {
			triggerTest(subT, yamlConfig)
		})
	}
}

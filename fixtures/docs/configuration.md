# Configuration

### Configuration options

- **Binary mode**

|Option          |Mandatory|Description|
|----------------|---------|-----------|
`--config-file`      |yes|option to specify the location of the configuration file
`--port`|no|option to specify the port (default: `5297`)
`--enable-https`     |no|option to enable or disable TLS communication (default: `false`)
`--cert-file`|no|option to specify the location of the certificate file (required if the `--enable-https` option is set to `true`)
`--key-file`|no|option to specify the location of the private key file (required if the `--enable-https` option is set to `true`)

- **Container mode**

|Environment variable|Mandatory|Description|
|--------------------|---------|-----------|
`API_GATEWAY_SQL_CONFIG_FILE`|no|a variable that specifies the location of the configuration file within the container (default: `/etc/api-gateway-sql/config.yaml`). It can be overwritten by an external file if the latter is mounted on a volume with the same name and in the same location.
`API_GATEWAY_SQL_PORT`|no|variable to specify the port (default: `5297`)
`API_GATEWAY_SQL_ENABLE_HTTPS`|no|variable to enable or disable TLS communication (default: `true`)
`API_GATEWAY_SQL_CERT_FILE`|no|variable to specify the location of the certificate file (required if the variable `API_GATEWAY_SQL_ENABLE_HTTPS` is set to `true`, default: `/etc/api-gateway-sql/tls/server.crt`)
`API_GATEWAY_SQL_KEY_FILE`|no|variable to specify the location of the private key file (required if the variable `API_GATEWAY_SQL_ENABLE_HTTPS` is set to `true`, default: `/etc/api-gateway-sql/tls/server.key`)

### Configuration file

```
api_gateway_sql:
  # Database configuration
  sqlitedb: "api_gateway_sql"
  # Configuration to enable or disable API documentation
  enable_swagger: true
  # Authentication parameter configuration
  auth:
    # Parameter to enable or disable authentication
    enabled: true
    # Username parameter used when authentication is enabled
    username: test
    # Password parameter used when authentication is enabled
    password: test@test
  # Target database parameter configuration
  databases:
    # Target identifier parameter
  - name: school
    # DBMS type parameter
    type: mariadb
    # Database host address parameter
    host: "@HOST_IP"
    # Database port parameter
    port: 3307
    # Database user parameter
    username: "test"
    # Database user password parameter
    password: "test"
    # Database name parameter
    dbname: "school"
    # Parameter to enable or disable SSL communication mode with the database
    sslmode: false
    # Database communication timeout parameter
    timeout: 1s
  # Target parameter configuration
  targets:
    # Target name parameter
  - name: insert_batch_student
    # Target database name parameter
    data_source_name: school
    # Parameter to enable or disable bulk execution
    multi: true
    # Batch size parameter to execute. Used when bulk execution is enabled
    batch_size: 10
    # Number of blocks parameter used to split the CSV file. Used when bulk execution is enabled
    buffer_size: 50
    # Database table fields parameter. Used when bulk execution is enabled
    batch_fields: "name;address"
    # SQL query content parameter
    sql: "insert into school (name, address) values ({{name}}, {{address}})"
```
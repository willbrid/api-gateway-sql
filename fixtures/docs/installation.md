## Installation

Here we install the **api-gateway-sql** application on a Linux machine:

- via **docker**: installation tested on Ubuntu 22.04 and Ubuntu 24.04
- via **podman**: installation tested on Rocky Linux 8.9

As a prerequisite, it is necessary to install a compatible DBMS, such as **MariaDB**, **MySQL**, **PostgreSQL**, **SqlServer**, or **Sqlite**, and configure a database; or to use an existing database from a compatible DBMS that is already installed. **MariaDB** is used here as an example for a test environment. The application also allows you to configure one or more databases from one or more different compatible DBMSs. During execution, it is possible to query a specific database by consuming an API and specifying a target in the URL parameter. This target refers to a configuration that contains both the reference to the database to query and the SQL string to execute. You can, for example, choose a containerized installation depending on your operating system. The link below will guide you through setting up a **MariaDB** sandbox with the **school** database:

[https://github.com/willbrid/api-gateway-sql/blob/main/fixtures/docs/databases.md](https://github.com/willbrid/api-gateway-sql/blob/main/fixtures/docs/databases.md).

Now let's install the **api-gateway-sql** application in a container.

```
mkdir $HOME/api-gateway-sql && $HOME/api-gateway-sql/data && cd $HOME/api-gateway-sql
```

```
vi config.yaml
```

```
api_gateway_sql:
  sqlitedb: "api_gateway_sql"
  auth:
    enabled: true
    username: test
    password: test@test
  databases:
  - name: school
    type: mariadb
    host: "127.0.0.1"
    port: 3307
    username: "test"
    password: "test"
    dbname: "school"
    sslmode: false
    timeout: 1s
  targets:
  - name: list-student
    data_source_name: school
    multi: false
    sql: "select * from student"
  - name: list-school
    data_source_name: school
    multi: false
    sql: "select * from school"
  - name: find-one-student
    data_source_name: school
    multi: false
    sql: "select * from student where id = {{id}}"
  - name: find-student-with-cond
    data_source_name: school
    multi: false
    sql: "select * from student where class_id = {{class}} and school_id = {{school}} and age >= {{age}}"
  - name: insert_school
    data_source_name: school
    multi: false
    sql: "insert into school (id, name, address) values ({{id}}, {{name}}, {{address}})"
  - name: insert_batch_school
    data_source_name: school
    multi: true
    batch_size: 10
    buffer_size: 50
    batch_fields: "name;address"
    sql: "insert into school (name, address) values ({{name}}, {{address}})"
```

- **Installation without persistence of SQLite application data**

Under Ubuntu
```
docker run -d --network=host --name api_gateway_sql -v $HOME/api_gateway_sql/config.yaml:/etc/api-gateway-sql/config.yaml -e API_GATEWAY_SQL_ENABLE_HTTPS=true willbrid/api-gateway-sql:latest
```

or

Under Rocky
```
podman run -d --net=host --name api_gateway_sql -v $HOME/api_gateway_sql/config.yaml:/etc/api-gateway-sql/config.yaml:z -e API_GATEWAY_SQL_ENABLE_HTTPS=true willbrid/api-gateway-sql:latest
```

- **Installation with persistent SQLite data for the application**

Under Ubuntu
```
docker run -d --network=host --name api_gateway_sql -v $HOME/api_gateway_sql/data:/data -v $HOME/api_gateway_sql/config.yaml:/etc/api-gateway-sql/config.yaml -e API_GATEWAY_SQL_ENABLE_HTTPS=true willbrid/api-gateway-sql:latest
```

or

Under Rocky
```
podman run -d --net=host --name api_gateway_sql -v $HOME/api_gateway_sql/data:/data:z -v $HOME/api_gateway_sql/config.yaml:/etc/api-gateway-sql/config.yaml:z -e API_GATEWAY_SQL_ENABLE_HTTPS=true willbrid/api-gateway-sql:latest
```

Once the installation is complete, to open Swagger via a browser, we access its page via the URL below.

```
https://localhost:5297/swagger/index.html
```
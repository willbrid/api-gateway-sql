## Usage

The goal of this content is to demonstrate the use of the **api-gateway-sql** application through its APIs. We will use the **curl** command to illustrate concrete examples of consuming these APIs.

#### Api [POST] : /v1/api-gateway-sql/{datasource}/init

This API can be used to create the database schema and insert data into it.

For our test environment, we will use the **school_mariadb.sql** file, which can be downloaded via the link below:

[https://github.com/willbrid/api-gateway-sql/blob/main/fixtures/sql/school_mariadb.sql](https://github.com/willbrid/api-gateway-sql/blob/main/fixtures/sql/school_mariadb.sql)

```
mkdir -p $HOME/api-gateway-sql && cd $HOME/api-gateway-sql
```

```
curl -fsSL https://github.com/willbrid/api-gateway-sql/raw/main/fixtures/sql/school_mariadb.sql -o school_mariadb.sql
```

```
curl -k -v -X POST -H 'Authorization: Basic dGVzdDp0ZXN0QHRlc3Q=' -H 'accept: application/json' -F "sqlfile=@school_mariadb.sql" https://localhost:5297/v1/api-gateway-sql/school/init
```

- The value of the **Basic** header represents the **base64** encoding of the application credentials (**username:password**) specified in its configuration file.

```
echo -n test:test@test | base64
```

- **school** is the name of the database connection string on **MariaDB** that we configured in the **api_gateway_sql.databases** section of the application's configuration file.

> Note: This API is not required if we are using the application with existing databases.

#### Api [GET] : /v1/api-gateway-sql/{target}

This API allows you to execute an SQL query based on the target name (**target**), which contains the query configuration. This SQL query must not be parameterized with the **{{}}** symbol.

```
curl -k -v -X GET -H 'Authorization: Basic dGVzdDp0ZXN0QHRlc3Q=' -H 'accept: application/json' -H 'Content-Type: application/json' https://localhost:5297/v1/api-gateway-sql/list-student
```

```
curl -k -v -X GET -H 'Authorization: Basic dGVzdDp0ZXN0QHRlc3Q=' -H 'accept: application/json' -H 'Content-Type: application/json' https://localhost:5297/v1/api-gateway-sql/list-school
```

**list-student** and **list-school** are the names of the targets configured in the **api_gateway_sql.targets** section of the configuration file :

- with **list-student**, its SQL query selects all rows from the **student** table
- with **list-school**, its SQL query selects all rows from the **school** table

#### Api [POST] : /v1/api-gateway-sql/{target}

This API allows you to execute an SQL query based on the target name (**target**), which contains the query configuration. This SQL query must be parameterized using one or more parameters defined by the **{{}}** symbol. The parameter values ​​must be passed via a POST request.

```
curl -k -v -X POST -H 'Authorization: Basic dGVzdDp0ZXN0QHRlc3Q=' -H 'accept: application/json' -H 'Content-Type: application/json' -d '{"id":"1"}' https://localhost:5297/v1/api-gateway-sql/find-one-student
```

```
curl -k -v -X POST -H 'Authorization: Basic dGVzdDp0ZXN0QHRlc3Q=' -H 'accept: application/json' -H 'Content-Type: application/json' -d '{"class":"1", "school":"1", "age":"15"}' https://localhost:5297/v1/api-gateway-sql/find-student-with-cond
```

```
curl -k -v -X POST -H 'Authorization: Basic dGVzdDp0ZXN0QHRlc3Q=' -H 'accept: application/json' -H 'Content-Type: application/json' -d '{"id":"300", "name":"high-tech", "address":"Willow Ave"}' https://localhost:5297/v1/api-gateway-sql/insert_school
```

**find-one-student**, **find-student-with-cond**, and **insert_school** are the names of the targets configured in the **api_gateway_sql.targets** section of the configuration file :

- with **find-one-student**, its SQL query selects a row from the **student** table where the **id** field value is **1**
- with **find-student-with-cond**, its SQL query selects all rows from the **student** table where the **class_id** field value is **1**, the **school_id** field value is **1**, and the **age** field value is greater than or equal to **15**
- with **insert_school**, its SQL query inserts a new row into the **school** table with the **id** field value being **300**, the **name** field value being **high-tech**, and the **address** field value being **high-tech** **Willow Ave**

> Note: Each parameter name sent via POST must be identical to a parameter name configured in the SQL query.

#### Api [POST] : /v1/api-gateway-sql/{target}/batch

This API allows you to execute a batch SQL query based on the target name (**target**), which contains the query configuration. The SQL query is parameterized using one or more parameters defined by the **{{}}** symbol, and the values ​​of these parameters must be sent via a CSV file in a POST request. The target configuration includes the following :

- Enabling batch mode (**multi: true**)
- Defining the maximum size of a data block in the CSV file (**buffer_size: 50**)
- Defining the maximum size of each batch within a block (**batch_size: 10**)
- Defining the parameter fields, where each field corresponds to a column in the CSV file, in order (**batch_fields: "name;address"**)

As a test, you can generate a 100-line CSV file with two columns : the first column contains the names of the schools, and the second their addresses. This CSV file can be generated using a Bash script available in this repository: [https://github.com/willbrid/api-gateway-sql/blob/main/fixtures/generate_schools.sh](https://github.com/willbrid/api-gateway-sql/blob/main/fixtures/generate_schools.sh).

```
mkdir -p $HOME/api_gateway_sql && cd $HOME/api_gateway_sql
```

```
curl -fsSL https://github.com/willbrid/api_gateway_sql/raw/main/fixtures/generate_schools.sh -o generate_schools.sh
```

```
chmod +x generate_schools.sh
./generate_schools.sh 100
```

This script will generate a csv file at the location **/tmp/schools.csv**.

```
curl -k -v -X POST -H 'Authorization: Basic dGVzdDp0ZXN0QHRlc3Q=' -H 'accept: application/json' -F "csvfile=@/tmp/schools.csv" https://localhost:5297/v1/api-gateway-sql/insert_batch_school/batch
```

**insert_batch_school** is the name of a target configured in the **api_gateway_sql.targets** section of the configuration file. This target allows batch and parallel execution of SQL inserts, retrieving values ​​from the **/tmp/schools.csv** file.

#### Api [GET] : /v1/api-gateway-sql/stats

This API allows you to view batch query execution statistics and track their progress.

```
curl -k -v -X GET -H 'Authorization: Basic dGVzdDp0ZXN0QHRlc3Q=' -H 'accept: application/json' 'https://localhost:5297/v1/api-gateway-sql/stats?page_num=1&page_size=20'
```

The API response provides execution information, including :
- the corresponding target
- for each block: its start number, end number, number of successes, number of failures, and the range of failed rows
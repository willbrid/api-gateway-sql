# Api-gateway-sql

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/willbrid/api_gateway_sql/blob/main/LICENSE) [![Build and Release](https://github.com/willbrid/api_gateway_sql/actions/workflows/cicd.yml/badge.svg)](https://github.com/willbrid/api_gateway_sql/actions/workflows/cicd.yml)

**Api-gateway-sql** is an application that allows SQL queries to be executed through an API. Each SQL query is defined in a configuration file and associated with a target. Query execution is performed by calling the corresponding API endpoint with the specified target.
The application supports both simple and batch queries and is compatible with several popular database management systems (DBMS), including: MySQL, MariaDB, PostgreSQL, SQL Server, and SQLite.

## Features

The **api-gateway-sql** application provides several features for executing SQL queries via an API, with support for simple, parameterized, and batch queries. Below is an overview of the main features:

- **Authentication configuration**

The application allows configuring *Basic authentication* to secure access to the API.

- **Execution of SQL queries from a SQL file (POST)**

The application allows executing SQL queries defined in a SQL file. This is particularly useful for database schema initialization.

- **Execution of SQL queries without parameters (GET)**

Some SQL queries can be executed without additional parameters. The API supports direct execution of such SQL queries via a GET request.

- **Execution of parameterized SQL queries using POST (POST)**

The application allows executing parameterized SQL queries by sending parameters through a POST request. This feature is ideal for dynamic queries where column values may change with each execution.

- **Batch execution of SQL queries using values from a CSV file (POST)**

The application supports batch execution of SQL queries by retrieving parameters from a CSV file. This is useful for automating the processing of large datasets in a single operation.

- **Batch execution statistics (GET)**

For each batch execution, the application provides access to statistics about the process, such as the number of successful and failed executions, total duration, and other relevant metrics.

## Documentation

1- [Configuration file](https://github.com/willbrid/api-gateway-sql/blob/main/fixtures/docs/configuration.md) <br>
2- [Starting a test database per DBMS](https://github.com/willbrid/api-gateway-sql/blob/main/fixtures/docs/databases.md) <br>
3- [Installation](https://github.com/willbrid/api-gateway-sql/blob/main/fixtures/docs/installation.md) <br>
4- [Usage](https://github.com/willbrid/api-gateway-sql/blob/main/fixtures/docs/usage.md)

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/willbrid/api-gateway-sql/blob/main/LICENSE) file for more details.
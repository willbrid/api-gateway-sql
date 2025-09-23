# Api_gateway_sql

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/willbrid/api_gateway_sql/blob/main/LICENSE)

**api_gateway_sql** est une application permettant d'exécuter des requêtes SQL via une API. Chaque requête SQL est définie dans un fichier de configuration et associée à une cible (appelée target). L'exécution de la requête s'effectue en appelant l'API correspondante, avec la cible spécifiée. L'application supporte aussi bien les requêtes simples que les requêtes en batch, et est compatible avec plusieurs systèmes de gestion de bases de données (SGBD) populaires, notamment : MySQL, MariaDB, PostgreSQL, SQL Server et SQLite.

## Fonctionnalités

L'application **api_gateway_sql** offre plusieurs fonctionnalités pour l'exécution de requêtes SQL via une API, avec un support pour des requêtes simples, paramétrées, et en batch. Voici un aperçu des fonctionnalités principales :

- **Configuration d'authentification**

L'application permet de configurer une authentification de type **basic** pour sécuriser l'accès à l'API.

- **Exécution de requêtes SQL via un fichier sql (POST)**

L'application permet d'exécuter des requêtes SQL définies dans un fichier sql. Utile dans le cas de l'initialisation du schéma d'une base de données.

- **Exécution de requêtes SQL sans paramètres (GET)**

Certaines requêtes sql peuvent être exécutées sans passer de paramètres supplémentaires. L'API supporte l'exécution directe d'une requête SQL par une requête GET.

- **Exécution de requêtes SQL avec des paramètres envoyés en POST (POST)**

L'application permet d'exécuter des requêtes SQL paramétrées en envoyant des paramètres via une requête POST. Cette fonctionnalité est idéale pour des requêtes dynamiques où les valeurs des colonnes peuvent changer à chaque exécution.

- **Exécution en batch d'une requête SQL avec des valeurs issues d'un fichier CSV (POST)**

L'application prend en charge l'exécution de requêtes SQL en batch en récupérant les paramètres d'un fichier CSV. Cela est utile pour automatiser le traitement d'un grand nombre de données en une seule requête.

- **Statistiques d'une exécution en batch (GET)**

Pour chaque exécution en batch, l'application permet d'obtenir des statistiques sur le processus, comme le nombre d'exécutions réussies et échouées, la durée totale, et d'autres métriques pertinentes.

## Documentation

1- [Fichier de configuration](https://github.com/willbrid/api_gateway_sql/blob/main/fixtures/docs/configuration.md) <br>
2- [Démarrage d'une base de données de test par SGBD](https://github.com/willbrid/api_gateway_sql/blob/main/fixtures/docs/databases.md) <br>
3- [Installation](https://github.com/willbrid/api_gateway_sql/blob/main/fixtures/docs/installation.md) <br>
4- [Utilisation](https://github.com/willbrid/api_gateway_sql/blob/main/fixtures/docs/utilisation.md)

## Licence

Ce projet est sous licence MIT - voir le fichier [LICENSE](https://github.com/willbrid/api_gateway_sql/blob/main/LICENSE) pour plus de détails.
# Coffeehouse

Coffeehouse is an example/reference service for keeping records of your favorite coffee recipes.

Coffeehouse includes the following features:
- RESTful API: Example CRUD operations for a coffee recipe
- Redis Caching: Example caching of the coffee recipe
- SQLC and PostgreSQL: Uses SQLC for type-safe Go from SQL and PostgreSQL for the database
- Example unit and integration tests for the server
- Docker Compose for running the server, database, nats, and redis
- Kubernetes deployment files for the server, database, nats, and redis
- Logging using Zap
- Pub/Sub using NATS
- Command line digital sign that subscribes to the NATS topic
- Taskfile for running common tasks that are used in development

## Directory
```
coffeehouse
│   routes.go - all the routes for the server
│   server.go - where the server struct is defined   
│   main.go - entry point of the application
|___cmd
│   │   digitalsign - command line digital sign
|___cache - caching interface
│   │   cache.go
|___|
|   |___redis - redis implementation
│   │   recipe.go
└───database
│   │   database.go
└───logger
│   │   logger.go
```

## Code Generation

Coffeehouse uses [sqlc](https://sqlc.dev/) for generating type-safe Go from SQL. It uses moq for mocking interfaces for use in testing.

## References

- https://github.com/benbjohnson/wtf
- https://github.com/ardanlabs/service
- https://pace.dev/blog/2018/05/09/how-I-write-http-services-after-eight-years.html
- https://brandur.org/sqlc

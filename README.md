# Coffeehouse

Coffeehouse is an example/reference service for keeping records of your favorite coffee recipes.

## Directory
```
coffeehouse
│   routes.go - all the routes for the server
│   server.go - where the server struct is defined   
└───database
│   │   database.go
└───logger
│   │   logger.go
```

## Code Generation

Coffeehouse uses [sqlc](https://sqlc.dev/) for generating type-safe Go from SQL.

## References

- https://github.com/benbjohnson/wtf
- https://github.com/ardanlabs/service
- https://pace.dev/blog/2018/05/09/how-I-write-http-services-after-eight-years.html
- https://brandur.org/sqlc

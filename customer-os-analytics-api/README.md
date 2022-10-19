# customerOs Analytics API

**GraphQL APIs to browse customer' analytics data**

This module is using [gorm](https://github.com/go-gorm/gorm) and [gqlgen](https://github.com/99designs/gqlgen)

## Quick start


1. Add any missing module requirements necessary to build the current module’s packages and dependencies, and remove requirements on modules that don’t provide any relevant packages.

       go mod tidy

2. Generate graphql generate models

       go run github.com/99designs/gqlgen generate

3. Set environment variables for DB connection:
   1. DB_HOST
   2. DB_PORT
   3. DB_NAME
   4. DB_USER
   5. DB_PWD


4. Start the graphql server

       go run server.go
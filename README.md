# edgecraft-api
> **[참고자료]**
> https://github.com/create-go-app/fiber-go-template
> https://github.com/golang-standards/project-layout

<img src="https://img.shields.io/badge/Go-1.17+-00ADD8?style=for-the-badge&logo=go" alt="go version" />&nbsp;<a href="https://goreportcard.com/report/github.com/create-go-app/fiber-go-template" target="_blank"><img src="https://img.shields.io/badge/Go_report-A+-success?style=for-the-badge&logo=none" alt="go report" /></a>&nbsp;<img src="https://img.shields.io/badge/license-Apache_2.0-red?style=for-the-badge&logo=none" alt="license" />

## 📦 Used packages

| Name                                                                  | Version    | Type       |
| --------------------------------------------------------------------- | ---------- | ---------- |
| [gofiber/fiber](https://github.com/gofiber/fiber)                     | `v2.34.0`  | core       |
| [gofiber/jwt](https://github.com/gofiber/jwt)                         | `v2.2.7`   | middleware |
| [arsmn/fiber-swagger](https://github.com/arsmn/fiber-swagger)         | `v2.31.1`  | middleware |
| [stretchr/testify](https://github.com/stretchr/testify)               | `v1.7.1`   | tests      |
| [golang-jwt/jwt](https://github.com/golang-jwt/jwt)                   | `v4.4.1`   | auth       |
| [joho/godotenv](https://github.com/joho/godotenv)                     | `v1.4.0`   | config     |
| [jmoiron/sqlx](https://github.com/jmoiron/sqlx)                       | `v1.3.5`   | database   |
| [jackc/pgx](https://github.com/jackc/pgx)                             | `v4.16.1`  | database   |
| [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)         | `v1.6.0`   | database   |
| [go-redis/redis](https://github.com/go-redis/redis)                   | `v8.11.5`  | cache      |
| [swaggo/swag](https://github.com/swaggo/swag)                         | `v1.8.2`   | utils      |
| [google/uuid](https://github.com/google/uuid)                         | `v1.3.0`   | utils      |
| [go-playground/validator](https://github.com/go-playground/validator) | `v10.10.0` | utils      |


## Project structure 
![Project Structure](./docs/images/Project-Structure.png)

## 🗄 Directory structure

### ./api
**Folder with OpenAPI/Swagger 스펙들.**

### ./cmd
**Main applications for this project.**

### ./internal
**Private application and library code.**. This is the code you don't want others importing in their applications or libraries.

- `./internal/cache` folder with in-memory cache setup functions (by default, Redis)
- `./internal/database` folder with database setup functions (by default, PostgreSQL)
- `./internal/controllers` folder for functional controllers (used in routes)
- `./internal/models` folder for describe business models and methods of your project
- `./internal/queries` folder for describe queries for models of your project
- `./internal/routes` folder for describe routes of your project


### ./pkg

**Library code that's ok to use by external applications.**. This directory contains all the project-specific code tailored only for your business use case, like _configs_, _middleware_, _routes_ or _utils_.

- `./pkg/configs` folder for configuration functions
- `./pkg/middleware` folder for add middleware (Fiber built-in and yours)
- `./pkg/repository` folder for describe `const` of your project
- `./pkg/utils` folder with utility functions (server starter, error checker, etc)

### ./docs

**Folder with 사용자 문서들.**

### ./scripts


# go-rss-agg-example

## Build and Run

```sh
# Build
go build
# Run
./go-rss-agg-example
# Build && Run
go build && ./go-rss-agg-example
```

### DB Migration

[Goose](http://pressly.github.io/goose/)

```sh
cd sql/schema
source ../../.env
goose $DB_URL up
goose $DB_URL down
```

### SQL Code Generation

[SQLC](https://docs.sqlc.dev/en/latest/index.html)

```sh
sqlc generate
```

# Project ds-easy


## Getting Started

### Install dependencies

to install tailwind : https://tailwindcss.com/blog/standalone-cli

to install go-migrate : https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

to install sqlc : https://docs.sqlc.dev/en/latest/overview/install.html



Make `.env` file like following :

```bash
PORT=8080
SESSION_STORE_KEY="session_store_key"
DB_URL=./test.db
```

migrate database
```bash
migrate -database "sqlite://test.db" -path db/migrations/ up
```
OR
```bash
make migrate
```

## MakeFile

run and live reload the application
```bash
make watch
```

clean up binary from the last build
```bash
make clean
```

# Project ds-easy


## Getting Started


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

## MakeFile

run and live reload the application
```bash
make watch
```

clean up binary from the last build
```bash
make clean
```

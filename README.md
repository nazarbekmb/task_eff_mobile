# task_eff_mobile

REST API service in Go, designed to store information about people enriched through external APIs (age, gender, nationality). The data is stored in PostgreSQL.

---

## Project launch

### 🐳 Via Docker (recommended)

```
make start
```

### Makefile commands

```
# Build and start the project
make start

# Stop and delete containers
make stop

# Apply migrations
make migrate-up

# Roll back migrations
make migrate-down

# Create a new migration
make create-migrations
```

---

## Example of running manually without Docker

If you already have local PostgreSQL running:

```
DB_HOST=localhost go run cmd/main.go
```

---

## Swagger UI

Once run, it is available at:

```
http://localhost:8080/swagger/index.html
```

---

## API Features

- `POST /people` - add a person
- `GET /people` - get all people (with filters)
- `PUT /people/:id` - update person
- `DELETE /people/:id` - delete person
- `GET /swagger/*any` - Swagger documentation

---
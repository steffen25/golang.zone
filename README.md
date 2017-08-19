# golang.zone
Home of golang.zone

The config/app.json should look like this
```json
{
  "env": "local",
  "mysql": {
    "username": "root",
    "password": "root",
    "database": "database_name",
    "encoding": "utf8mb4"
  },
  "redis": {
    "host": "localhost",
    "port": 6379
  },
  "port": 8080,
  "jwt_secret": "secret"
}
```

### Prerequisites
- mysql
- redis on port 6379 (revoke JWTs)

# golang.zone
Home of golang.zone

The config/app.json should look like this
```json
{
  "env": "local",
  "database": {
    "username": "root",
    "password": "root",
    "database": "test",
    "encoding": "utf8mb4"
  },
  "port": 8080,
  "jwt_secret": "secret"
}
```

# golang.zone
Home of golang.zone

This repository holds the files for the REST API of https://golang.zone

You need to make a config file in the config directory - below is an example of what it could look like
```json
{
  "env": "local",
  "mysql": {
    "host": "app_mysql",
    "username": "root",
    "password": "root",
    "database": "database_name",
    "encoding": "utf8mb4",
    "port": "3306"
  },
  "redis": {
    "host": "app_redis",
    "port": 6379
  },
  "jwt": {
      "secret": "secret",
      "public_key_path": "config/api.rsa.pub", // openssl genrsa -out api.rsa keysize
      "private_key_path": "config/api.rsa" // openssl rsa -in api.rsa -pubout > api.rsa.pub
  ,
  "port": 8080
}
```

### Prerequisites
- Go
- Mysql
- Redis (To revoke JWTs)

### Endpoints


Users:

| URL        								| Method           	| Info  |
| ------------- 							|:-------------:	| -----:|
| http://127.0.0.1:8080/api/v1/users      			| GET 				| Returns an array of users |
| http://127.0.0.1:8080/api/v1/users      			| POST 				| Endpoint to create a user |
| http://127.0.0.1:8080/api/v1/users/{id}      		| GET 				| Returns a single user specificed by a id |
| http://127.0.0.1:8080/api/v1/users/{id}/posts    	| GET 				| Returns an array of posts created by an user |


Posts:

| URL        									| Method           	| Info  |
| ------------- 								|:-------------:	| -----:|
| http://127.0.0.1:8080/api/v1/posts      		| GET 				| Returns an array of posts |
| http://127.0.0.1:8080/api/v1/posts      		| POST 				| Endpoint to create a post - proctected by a auth middleware that requires the user to be authenticated and be an admin |
| http://127.0.0.1:8080/api/v1/posts/{id}      	| PUT 				| Endpoint to update a specific post - proctected by a auth middleware that requires the user to be authenticated and be an admin |

Auth:

| URL        												| Method           	| Info  |
| ------------- 											|:-------------:	| -----:|
| http://127.0.0.1:8080/api/v1/auth/login      				| POST 				| Endpoint to authenticate a user- returns a JWT that last for 24 hours |
| http://127.0.0.1:8080/api/v1/auth/update      			| PUT 				| Endpoint to update the current user that is logged in - proctected by a auth middleware |
| http://127.0.0.1:8080/api/v1/auth/refresh      			| GET 				| Endpoint to refresh a JWT - proctected by a auth middleware |
| http://127.0.0.1:8080/api/v1/auth/logout      			| GET 				| Logout the user from the API by revoking the user's JWT - proctected by a auth middleware |
| http://127.0.0.1:8080/api/v1/auth/logout/all      		| GET 				| Logout the user from all the places that he/she is logged in - proctected by a auth middleware |


Misc:

| URL        								| Method           	| Info  |
| ------------- 							|:-------------:	| -----:|
| http://127.0.0.1:8080       				| GET 				| A test endpoint to see the API is working (hello world) |
| http://127.0.0.1:8080/api/v1/protected    | GET 				| A test endpoint to see the id of the user that is logged in - proctected by a auth middleware |


### Example requests

todo


### Docker Development
Run the following commands only before boostrapping the application.
```sh
# Migrate database
docker-compose run --rm --name app -p 8080:8080 app_api bash
# Inside the app container
# TODO: Clean up the migration
goose -dir ./db/migrations/ mysql "app:password@(app_mysql:3306)/app" up
```

Development - Since linking the local folder to container we can edit code and run container.
```sh
docker-compose up
```
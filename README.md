<p align="center">
    <img src="https://i.imgur.com/rZY34qI.png" height="130">
</p>
<p align="center">

  <a href="https://travis-ci.org/steffen25/golang.zone">
  	<img src="https://travis-ci.org/steffen25/golang.zone.svg?branch=master"
  alt="">
  </a>

  <a href="https://goreportcard.com/report/github.com/steffen25/golang.zone">
  	<img src="https://goreportcard.com/badge/github.com/steffen25/golang.zone" alt="">
  </a>
  
  <a href="https://coveralls.io/github/steffen25/golang.zone?branch=master">
  	<img src="https://coveralls.io/repos/github/steffen25/golang.zone/badge.svg?branch=master" alt="">
  </a>

  <a href="https://golang.zone">
    <img src="https://img.shields.io/website-up-down-green-red/http/shields.io.svg?label=golang.zone" >
  </a>

  <a href="LICENSE">
  <img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt=""></a>

</p>


# golang.zone

The frontend app for this project is made with Vue.js and is open source 

[golang.zone-frontend](https://github.com/steffen25/golang.zone-frontend)

This repository holds the files for the REST API of https://golang.zone

You need to make a config file in the config directory - below is an example of what it could look like. You can also directly use the example file in the config directory just remember to remove the .example extension
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
      "public_key_path": "config/api.rsa.pub",
      "private_key_path": "config/api.rsa"
  },
  "port": 8080
}
```
Generate private key

- openssl genrsa -out api.rsa keysize(2048 or 4096)

Generate public key
- openssl rsa -in api.rsa -pubout > api.rsa.pub


### Prerequisites
- Go 1.7+
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
| http://127.0.0.1:8080/api/v1/images/upload    | POST 				| Endpoint to upload images to a post using a wysiwyg editor - proctected by a auth middleware that require the user to be authenticated and be an admin |
| http://127.0.0.1:8080/api/v1/posts/{id}      	| GET 				| Endpoint to retrieve a post specified by an id |
| http://127.0.0.1:8080/api/v1/posts/{slug}     | GET 				| Endpoint to retrieve a post specified by a slug |
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

#### Paw
<i>Mac users only</i>

[golang.zone.paw](examples/golang.zone.paw)

Open the file with Paw then set the environment in the upper left corner

#### Postman
1. Import the [golang.zone.postman_collection.json](examples/golang.zone.postman_collection.json) in Postman
2. Setup your environments by clicking the 🔩 icon just next to the 👁️ icon top in the top right corner
3. You can also import the environments created by me, 
[golang.zone.Local.postman_environment.json](examples/golang.zone.Local.postman_environment.json), [golang.zone.Prod.postman_environment.json](examples/golang.zone.Prod.postman_environment.json)
continue at step 8.
4. I suggest creating 2 envs on called golang.zone Local and one called golang.zone Prod
5. Create 3 keys BASE_URL, ACCESS_TOKEN and REFRESH_TOKEN
6. Set the value of BASE_URL to http://localhost:8080 for the local env and https://golang.zone for the Prod env, leave the tokens value empty for now
7. Click Add/Update.
8. Go to the Login request and open up the Tests tab.
9. Insert the code below then save and finally click send.
```javascript
var jsonData = JSON.parse(responseBody);
postman.setEnvironmentVariable("ACCESS_TOKEN", jsonData.data.tokens.accessToken);
postman.setEnvironmentVariable("REFRESH_TOKEN", jsonData.data.tokens.refreshToken);
```
10. Now when you call endpoints that require authorization it will automatically insert the token value inside the Authorization header




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

# To do
- [ ] Add categories for posts
- [ ] Add tags for posts
- [ ] Add featured image for posts

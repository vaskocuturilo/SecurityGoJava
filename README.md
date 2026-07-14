##      

The demo project for Security and JWT auth uses Golang and Java.

### This project was created to demonstrate how to work with the Security.

- Create a REST API security service (/auth/register/ auth/login /auth/refresh /auth/logout /tasks(GET|POST)).
- Add TTL(access token - 5 minutes, refresh token - 7 days).
- Use secret from env.
- Add Correct Global exception.
- The refresh token should be saved to the database.
- Add Logout functionality.
- Add Rate Limiting.
- Add Bcrypt cost (it should be config)
- Add Role-based functionality.
- Add Docker file.
- Add Docker Compose.
- Add unit(repository, service and rest) and integration tests with testcontainers.

You will need the following technologies available to try it out:

* Git
* Spring Boot 4+
* Gradle 9+
* JDK 24+
* Spring Security 7+
* Golang 1.25+
* Any HTTP router (net/http, chi, gin)
* Any JWT library(golang-jwt/jwt, jsonwebtoken)
* Docker
* Docker compose
* IDE of your choice

### How to run via Spring Boot.

``` ./gradlew bootRun ```

``` docker compose -f "docker-compose-java.yml" up --detach ```

### How to run via Golang.

``` go run .```

``` docker compose -f "docker-compose-golang.yml" up --detach ```

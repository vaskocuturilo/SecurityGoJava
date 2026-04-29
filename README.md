##      

The demo project for Security and JWT auth working with Golang and Java.

### This project was created to demonstrate how to work with the Securuty.

- Create a REST API security service (/auth/register/ auth/login /auth/refresh /tasks).
- Add TTL(access token - between 10 and 15 minutes, refresh token - approximately 7 days).
- Use secret from env.
- Add Correct Global exception.
- The refresh token should be saved to the database.
- Add Logout functionality.
- Add Role-based
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

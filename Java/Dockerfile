## Build stage
FROM gradle:jdk-21-and-24 AS build

COPY --chown=gradle:gradle . /home/gradle/src

WORKDIR /home/gradle/src

RUN gradle clean build -x test --no-daemon --refresh-dependencies || return 0

## Package stage
FROM openjdk:24-ea-1-jdk-slim

ENV ARTIFACT_NAME=api-countries-1.0.0.jar

ENV APP_HOME=/home/gradle/src/

RUN mkdir /app

COPY --from=build $APP_HOME/build/libs/$ARTIFACT_NAME app.jar

ENTRYPOINT ["java", "-jar","-Dspring.profiles.active=docker", "/app.jar"]

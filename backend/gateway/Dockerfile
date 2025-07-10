FROM maven:3.9.6-eclipse-temurin-17-focal AS build

WORKDIR /app

COPY pom.xml .
COPY gateway/pom.xml ./gateway/

RUN mvn -B -pl gateway -am dependency:go-offline

COPY . .

RUN mvn -B clean install -pl gateway -am


FROM eclipse-temurin:17-jre-jammy

WORKDIR /app

COPY --from=build /app/gateway/target/*.jar ./app.jar

EXPOSE 8080

ENTRYPOINT ["java", "-jar", "app.jar"]

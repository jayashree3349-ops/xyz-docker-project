# Stage 1: Build the Java application
FROM maven:3.9-eclipse-temurin-17 AS builder
WORKDIR /build
COPY pom.xml .
COPY src ./src
RUN mvn -B package

# Stage 2: Lightweight Java runtime
FROM eclipse-temurin:17-jre
WORKDIR /app
COPY --from=builder /build/target/java-service-1.0.0.jar app.jar
USER 10001
CMD ["java", "-jar", "app.jar"]

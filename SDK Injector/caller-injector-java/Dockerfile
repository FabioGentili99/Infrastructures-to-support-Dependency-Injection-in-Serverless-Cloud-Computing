FROM openjdk:21-slim-buster
EXPOSE 8080
COPY target/caller-java-sdk-0.0.1-SNAPSHOT.jar app.jar
ENTRYPOINT ["java","-jar","/app.jar"]
package com.example;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.cache.annotation.EnableCaching;
import org.springframework.data.mongodb.repository.config.EnableMongoRepositories;


@SpringBootApplication
@EnableCaching
@EnableMongoRepositories
public class App {
  public static void main(String[] args) {
    SpringApplication.run(App.class, args);
  }
}

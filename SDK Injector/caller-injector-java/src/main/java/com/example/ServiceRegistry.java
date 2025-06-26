package com.example;
import org.springframework.data.mongodb.repository.MongoRepository;
import org.springframework.data.mongodb.repository.Query;

import java.util.Optional;

public interface ServiceRegistry extends MongoRepository<Service, String> {
  @Query("{id:'?0'}")
  Optional<Service> findServiceById(String id);
}
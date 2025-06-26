package com.example;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Component;
import reactor.core.publisher.Mono;

import java.util.Optional;

@Component
public class Injector {

  private static final Logger log = LoggerFactory.getLogger(Injector.class);

  private final ServiceRegistry serviceConfigRepository;

  public Injector(ServiceRegistry serviceConfigRepository) {
    this.serviceConfigRepository = serviceConfigRepository;
  }

  /**
   * Retrieves a service address from MongoDB based on a configuration ID.
   * Results are cached to avoid repeated database lookups for the same ID.
   *
   * @param id The ID of the service configuration to retrieve (e.g., "helloWorldService").
   */
  @Cacheable(value = "serviceAddresses", key = "#id")
  public String getServiceAddress(String id) {
    //log.info("Attempting to retrieve service configuration for ID '{}' from MongoDB (or cache miss).", configId); // Log when actual DB call happens
    Optional<Service> service = serviceConfigRepository.findServiceById(id);
    if (service.isEmpty()){
      return null;
    } else {
      return service.get().getServiceAddress();
    }
  }

  
}
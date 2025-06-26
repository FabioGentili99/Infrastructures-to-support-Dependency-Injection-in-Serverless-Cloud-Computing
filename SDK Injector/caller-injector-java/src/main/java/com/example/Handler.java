package com.example;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.context.annotation.Bean;
import org.springframework.stereotype.Component;
import org.springframework.web.reactive.function.client.WebClient;
import org.springframework.web.util.UriBuilder;
import reactor.core.publisher.Mono;

import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.time.Duration;
import java.time.Instant;
import java.time.LocalDateTime;
import java.time.ZoneOffset;
import java.util.function.Function;

@Component
public class Handler {

  private static final Logger log = LoggerFactory.getLogger(Handler.class);

  private final Injector injector;
  private final WebClient webClient;

  public Handler(Injector injector, WebClient.Builder webClientBuilder) {
    this.injector = injector;
    // Build WebClient. Consider adding retries, timeouts, etc., for production.
    this.webClient = webClientBuilder.build();
  }

  @Bean
  public Function<Request, Response> processTimestamp() {
    return request -> {
      String timestampMessage = request.getMessage();
      long millis = Long.parseLong(timestampMessage);
      Instant begin = Instant.ofEpochMilli(millis);
      //log.info("Received request with timestamp message: {}", timestampMessage);

      

      Instant now = Instant.now();
      long start = now.getEpochSecond() * 1_000_000 + now.getNano() / 1_000;
      String service = injector.getServiceAddress("hello");
      now = Instant.now();
      long end = now.getEpochSecond() * 1_000_000 + now.getNano() / 1_000;

      log.info(LocalDateTime.now(ZoneOffset.UTC) + ", " + "CALLER, info, " + "Service retrieved in " + (end - start) / 1000.0 + " ms");
      //log.info("service address: " + service);

      now = Instant.now();
      long start_invoke = now.getEpochSecond() * 1_000_000 + now.getNano() / 1_000;

      String externalResponse = invoke(service);

      now = Instant.now();
      long end_invoke = now.getEpochSecond() * 1_000_000 + now.getNano() / 1_000;
      Instant end2 = Instant.now();
      Duration duration = Duration.between(begin, end2);
      log.info(LocalDateTime.now(ZoneOffset.UTC) + ", " + "CALLER, info, " + "Service invoked in " + (end_invoke - start_invoke) / 1000.0 + " ms");
      log.info(LocalDateTime.now(ZoneOffset.UTC) + ", " + "CALLER, info, " + "Total latency is " + duration.toMillis() + " ms");

      return new Response("SUCCESS", externalResponse);
    };
  }


  public String invoke(String uri) {
    HttpClient client = HttpClient.newHttpClient();
    HttpRequest request = HttpRequest.newBuilder()
      .uri(URI.create(uri))
      .GET()
      .build();

    HttpResponse<String> response = null;
    try {
      response = client.send(request, HttpResponse.BodyHandlers.ofString());
    } catch (IOException e) {
      throw new RuntimeException(e);
    } catch (InterruptedException e) {
      throw new RuntimeException(e);
    }
    return response.body();
  }

}
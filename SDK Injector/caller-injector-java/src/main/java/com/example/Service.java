package com.example;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;
import org.bson.types.ObjectId;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.Document;
import org.springframework.data.mongodb.core.mapping.Field;

@Data
@NoArgsConstructor
@AllArgsConstructor
@Document(collection = "services") // Collection name in MongoDB
public class Service {
  @Id
  private ObjectId _id;
  @Field(name = "id")
  private String id; // e.g., "hello"
  @Field(name = "ServiceName")
  private String ServiceName;
  @Field(name = "ServiceAddress")
  private String ServiceAddress; // e.g., "http://your-hello-world-service.com/hello"
}
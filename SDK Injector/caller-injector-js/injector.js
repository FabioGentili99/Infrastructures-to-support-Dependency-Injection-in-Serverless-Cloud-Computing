const { MongoClient } = require('mongodb');
const { createLogger, format, transports } = require("winston");

const mongoURI = process.env.MONGO_URI || 'mongodb://localhost:27017';
const cache = new Map();


class Injector {


    constructor() {
        this.logger = createLogger({
            level: 'info',
            format: format.combine(
                format.timestamp(),
                format.printf(({ timestamp, level, message }) => {
                return `${new Date(timestamp).toISOString()},INJECTOR,${level},${message}`;
            })),
            transports: [
                new transports.Console()
            ]
        });
        this.dbUrl = mongoURI;
        this.dbName = 'services';
        this.collectionName = 'services';
        this.client = new MongoClient(mongoURI);
        this.connect();
    }

    connect() {
        if (!this.client.isConnected) {
            this.client.connect();
        }
        this.db = this.client.db(this.dbName);
        this.collection = this.db.collection(this.collectionName);
      }


    /**
   * Registers a new service in the database.
   * @param {string} id - The id of the service.
   * @param {string} name - The name of the service.
   * @param {string} topic - The topic of the service.
   */
  async registerService(id, name, address) {
    const service = { id: id,
                      ServiceName: name,
                      ServiceAddress: address };
    await this.collection.insertOne(service, function(err, res) {
        if (err) throw err;
        console.log("1 document inserted");
    });
  }


    /**
   * Retrieves a service by its ID.
   * @param {string} id - The ID of the service.
   * @returns {Promise<Object|null>} The service document or null if not found.
   */
  async getServiceById(id) {
    
    this.logger.info(`Fetching service with ID: ${id}`);

    const start = Date.now();
    //console.log(start, "start time")

    // Check if the service is in cache
    if (cache.has(id)) {
        const service = cache.get(id);
        const end = Date.now();
        this.logger.info(`Service retrieved in ${(end - start).toFixed(3)} ms`);
        return service
    }

    try {
        const service = await this.collection.findOne({ id: id });
        if (!service) {
            this.logger.info(`Error finding service with id '${id}': service not found`);
            return service
        }

        const end = Date.now();
        this.logger.info(`Service retrieved in ${(end - start).toFixed(3)} ms`);
        // Store in cache
        cache.set(id, service);
        return service
    } catch (err) {
        this.logger.info(`Error finding service with id '${id}': ${err}`);
        return null;
    }
  }

}

module.exports = Injector;
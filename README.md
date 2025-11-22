
# BSF Engine (Bahrain Speed Finder)

The **BSF Engine** is a Go-based backend service that provides speed limit information for roads in Bahrain using OpenStreetMap data. Given a latitude and longitude, it returns the speed limit for the nearest road segment. The service is optimized for fast lookups using Redis and is containerized for easy deployment.

---

## Features

- **Speed Limit Lookup:** Returns the speed limit for a given latitude and longitude.
- **OpenStreetMap Data:** Uses preprocessed OSM data for Bahrain.
- **Redis Caching:** Fast geospatial queries using Redis GEO features.
- **REST API:** Simple HTTP endpoint for integration.
- **Dockerized:** Easy to run locally or in production with Docker Compose.

---

## Project Structure

```
.
├── assets/                # Contains the Bahrain OSM road data (bahrain_roads.json)
├── entities/              # Data models and utility functions
├── handlers/              # HTTP handler logic
├── loader/                # Data loading logic (if any)
├── receiver/              # Redis interaction and geospatial logic
├── routes/                # API route definitions
├── main.go                # Application entry point
├── Dockerfile             # Docker build instructions
├── docker-compose.yml     # Multi-container setup (app + Redis)
├── go.mod                 # Go module definition
└── README.md              # Project documentation
```

---

## How It Works

1. **Startup:** On launch, the service loads OSM road data from `assets/bahrain_roads.json` into Redis and an in-memory cache.
2. **API Endpoint:** The `/get-speed-limit` endpoint accepts `lat` and `lon` query parameters.
3. **Geospatial Query:** The service queries Redis for the nearest road segment and returns its speed limit.
4. **Periodic Reload:** Road data is reloaded into Redis every 12 hours to ensure freshness.

---

## API Usage

### Endpoint

```
GET /get-speed-limit?lat=<latitude>&lon=<longitude>
```

#### Query Parameters

- `lat`: Latitude (required)
- `lon`: Longitude (required)

#### Example Request

```
GET /get-speed-limit?lat=26.2285&lon=50.5860
```

#### Example Response

```json
{
  "speed": 80
}
```

If no data is found, a default speed of 50 is returned with an error message.

---

## Running the Service

### Prerequisites

- [Docker](https://www.docker.com/get-started) and [Docker Compose](https://docs.docker.com/compose/)
- (Optional) Go 1.25+ if running natively

### Using Docker Compose (Recommended)

1. **Build and Start:**

	```sh
	docker-compose up --build
	```

	This will start both the Redis service and the Go application.

2. **Access the API:**

	The service will be available at [http://localhost:9090/get-speed-limit](http://localhost:9090/get-speed-limit).

### Native Go Run (Advanced)

1. **Start Redis:**

	```sh
	docker run -p 6379:6379 redis:7
	```

2. **Set Environment Variables (optional):**

	- `PORT` (default: 9090)
	- `REDIS_HOST` (default: localhost)
	- `REDIS_PORT` (default: 6379)

3. **Run the App:**

	```sh
	go run main.go
	```

---

## Environment Variables

- `PORT`: Port for the HTTP server (default: 9090)
- `REDIS_HOST`: Redis server hostname (default: localhost)
- `REDIS_PORT`: Redis server port (default: 6379)

---

## Data Preparation

- The `assets/bahrain_roads.json` file must be present and contain the preprocessed OSM road data for Bahrain.
- This file is loaded into Redis at startup and periodically refreshed.

---

## Development

- Code is organized by responsibility (entities, handlers, receiver, etc.).
- Uses the Gin web framework for HTTP routing.
- Redis is used for fast geospatial queries and caching.

---

## License

MIT License

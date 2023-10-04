# Feed service

- This service provides endpoints for clients to retrieve news articles. It fetches the stored news data from MongoDB and presents it to clients in JSON format

## Project Structure

### Directory structure

```text
├── build             # Docker files required for building the project
├── cmd               # Application entry points
│   └── app
├── config            # Application configuration files
├── deployment        # Deployment-related files
├── internal          # Internal application code
│   ├── config        # Code for configuration
│   ├── model         # Code defining data structures
│   ├── repository    # Code for data storage (in this case, "feed")
│   │   └── feed
│   ├── service       # Code for business logic (in this case, "feed")
│   │   └── feed
│   └── transport     # Code for handling transport mechanisms
├── postman           # Postman collections
├── go.mod            # Go module file for managing dependencies
├── go.sum            # Go module checksum file
├── Makefile          # File defining command-line automation
└── README.md         # Project description and instructions
```

## Local Setup

### In docker

Docker-compose is used for deployment

1. Create an environment variables file in the root directory
```text
cp ./.env.dict ./.env
```

2. Change environment variables if you need
3. Build the local service image
```text
make build-image
```
4. Run the project using Docker Compose
```text
make deploy
```
5. Check containers
```text
docker ps
```
6. You can use the [postman collection](postman) or cURL below to test API methods 
```text
curl --location --request GET 'http://localhost:8084/api/feed/v1/news?limit=5'

curl --location --request GET 'http://localhost:8084/api/feed/v1/news?limit=5&cursor=2023-08-16T08:00:00Z'

curl --location --request GET 'http://localhost:8084/api/feed/v1/news/64dcc2326c1e91ddcad83046'
```
7. You can watch docker logs in real time
```text
docker logs -f feed-service
```
8. You can remove services after testing
```text
make delete
```

### Locally

1. You can create you own configuration file and specify it in environment variables
```text
cp ./config/default.yaml ./config/local.yaml

CONFIG_FILE=./config/local.yaml
```
2. Run external services like db, broker, cache etc. using docker compose file
```text
make deploy services='mongodb'
```
3. It is ready, you can run this application, for example, via IntelliJ IDEA.

## Testing
- You can run unit tests using cmd ```make test-unit```

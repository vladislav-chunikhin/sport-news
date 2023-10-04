# Feed transformer service

- It listens to RabbitMQ for incoming news data. Once received, it processes and transforms the XML data into the desired model, then stores it in MongoDB

## Project Structure

### Directory structure

```text
├── build                 # Docker files required for building the project
├── cmd                   # Application entry points
│   └── app
├── config                # Application configuration files
├── deployment            # Deployment-related files
├── internal              # Internal application code
│   ├── config            # Code for configuration
│   ├── mapper            # Code for data mapping/transformations
│   ├── provider          # Code providing data or services
│   │   └── htafc
│   ├── rabbitmq          # Code for RabbitMQ interaction (message consumer)
│   ├── repository        # Code for data storage
│   └── transformer       # Code for data transformation
├── go.mod                # Go module file for managing dependencies
├── go.sum                # Go module checksum file
├── Makefile              # File defining command-line automation
└── README.md             # Project description and instructions
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
6. You can use the rabbitMQ UI to inspect messages in queues
```text
http://localhost:15672/
```
7. You can watch docker logs in real time
```text
docker logs -f feed-transformer
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
make deploy services='rabbitmq mongodb redis'
```
3. It is ready, you can run this application, for example, via IntelliJ IDEA.

## Testing
- You can run unit tests using cmd ```make test-unit```


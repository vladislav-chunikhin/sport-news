# Improvements

- Broker: Feed transformer service can handle multiple messages concurrently from rabbitmq
- Broker: Regulate message consumption rate to balance processing and prevent overload
- Broker: We can utilize a dead letter queue to process failed messages at a later time
- Broker: Retry mechanisms for failed message deliveries
- Deployment: We can use k8s for orchestrating containers in production, but docker-compose is a good solution for local deployment
- Testing: Cover more code with unit and integration tests to automate testing and accelerate the delivery to production
- Monitoring and Metrics: We can implement Prometheus metrics to collect them and use Grafana to monitor consumption
- Caching: Use cache for frequently used requests from clients
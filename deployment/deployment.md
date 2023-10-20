## Application Deployment Guide

This guide provides step-by-step instructions to successfully deploy and run the Sport News application.

### Prerequisites

Before you begin, ensure you have the following installed:

- Docker
- Docker-Compose

### Algorithm

1. Deploy sport news application
    1. run cmd ```make deploy-local``` to deploy application
2. After testing, you can remove containers
    1. run cmd ```make delete-local``` to delete containers

**Note**: Each service has its own ```docker-compose.yaml```.
It is necessary to enable independent testing and deployment if you only need to test a single service.

### API

```text
curl --location --request GET 'http://localhost:8084/api/feed/v1/news?limit=5'

curl --location --request GET 'http://localhost:8084/api/feed/v1/news?limit=5&cursor=2023-08-16T08:00:00Z'

curl --location --request GET 'http://localhost:8084/api/feed/v1/news/64dcc2326c1e91ddcad83046'
```

- limit=5: This parameter indicates that the API should return a maximum of 5 news items. It's used to specify the
  number of items you want to retrieve from the API.
- cursor=2023-08-16T08:00:00Z: This parameter (published) might be used to specify a starting point for fetching news
  items. The term "cursor" often refers to a pointer or indicator that shows the position or point of reference in a
  dataset. In this case, the value 2023-08-16T08:00:00Z is a timestamp (in ISO 8601 format) that could be used to
  indicate the point in time from which you want to start retrieving news items.

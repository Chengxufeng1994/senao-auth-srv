# Senao Authenticate Server

## Introduce

It will provide APIs for the frontend to do the following things:

* Create and verify accounts.

## Setup local development

* Docker desktop
* Golang

## Setup infrastructure

* Create the network
  ```bash
  docker network create senao-network
  ```

* Create the redis

  ```bash
  make redis
  ```

## How to run

1. Local
   ```bash
   make server
   ```

2. Docker

    * Build docker image

      ```bash
      docker build -t <image-name> --no-cache -f Dockerfile .
      ```

    * Start docker container

      ```bash
      docker run --name <container-name> --it -p 8000:8000 <image-naem>
      ```

3. Docker-compose

    ```bash
    docker compose up --build
    ```

4. From Docker hub
    * Pull Image
    
       ```bash
       docker pull benny0329/senao-auth-srv
       ```
      
    * Start docker container

      ```bash
      docker run --name <container-name> --it -p 8000:8000 --network=senao-network \
      -e REDIS_HOST=<redis_container_name>
      <image-name>
       ```

### Environment variables

| variables      | Description            |
|----------------|------------------------|
| SERVER_HOST    | Specify server host    |
| SERVER_PORT    | Specify server port    |
| REDIS_HOST     | Specify redis host     |
| REDIS_PORT     | Specify redis port     |
| REDIS_PASSWORD | Specify redis password |

## Documentation

following this address: http://localhost:8000/swagger/index.html

## Basic Requirement

* [x] RESTful API
* [x] Error handling and Simple input validation
* [x] Data storage
    * Timeout retry mechanism: using redis to set timeout
* [x] Containerize
* [x] Open Swagger API
* [x] User Guide (README.md)

## TODO

* [x] Support docker-compose
* [ ] Relational database
* [ ] Microservice
* [ ] Kubernetes (deployment, service)
* [ ] View Log (ELK)
* [ ] Monitoring (Grafana, Prometheus)
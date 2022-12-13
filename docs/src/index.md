# Installation

Work in progress!  
Not production ready! Use at your own risk!

## <a name="docker-images"></a> Docker images

Lookup latest release at [GHCR](https://github.com/lostinsoba/ninhydrin/pkgs/container/ninhydrin)

## <a name="configuration"></a> Configuration

Docker image services entrypoints:

| Service   | Entrypoint              | Ports                         |
| --------- | ----------------------- | ----------------------------- |
| API       | `./ninhydrin/api`       | 8080 (API), 8081 (Monitoring) |
| Scheduler | `./ninhydrin/scheduler` | 8081 (Monitoring)             |

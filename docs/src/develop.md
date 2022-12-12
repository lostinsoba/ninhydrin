# Development environment

1. [Installation](#installation)  
1.1 [Docker Compose](#installation-docker-compose)  

2. [Debugging](#debugging)  

## <a name="installation"></a> Installation

Fork the repository and/or clone the `develop` branch:

```
git clone -b develop https://github.com/lostinsoba/ninhydrin
```

### <a name="installation-docker-compose"></a> Docker Compose

```
make develop-compose
```

## <a name="debugging"></a> Debugging

List of running services:

| Service                     | Description                          | Network ports                 |
| --------------------------- | ------------------------------------ | ----------------------------- |
| ninhydrin-api               | Ninhydrin API                        | 8080 (API), 8081 (Monitoring) |
| ninhydrin-scheduler         | Ninhydrin Scheduler                  | 8082 (Monitoring)             |
| ninhydrin-storage           | Storage (default: PostgreSQL)        | 5432                          |
| ninhydrin-monitoring-source | Metric source (default: Prometheus)  | 9090                          |
| ninhydrin-monitoring-ui     | Monitoring Web UI (default: Grafana) | 3000                          |

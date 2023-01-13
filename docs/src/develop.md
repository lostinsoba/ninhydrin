# Development environment

## <a name="installation"></a> Installation

Fork the repository and/or clone the `develop` branch:

```
git clone -b develop https://github.com/lostinsoba/ninhydrin
```

### <a name="installation-docker-compose"></a> Docker Compose

```
make develop-compose arg=value [arg...]
```

List of arguments:

| Argument        | Description         | Values                                  |
|-----------------|---------------------|-----------------------------------------|
| DEVELOP_STORAGE | Storage kind to use | postgres, redis (redis, redis.sentinel) |

## <a name="debugging"></a> Debugging

List of running services:

| Service                     | Description                          | Network ports                 |
| --------------------------- | ------------------------------------ | ----------------------------- |
| ninhydrin-api               | Ninhydrin API                        | 8080 (API), 8081 (Monitoring) |
| ninhydrin-scheduler         | Ninhydrin Scheduler                  | 8082 (Monitoring)             |
| ninhydrin-storage           | Storage (default: PostgreSQL)        | 5432                          |
| ninhydrin-monitoring-source | Metric source (default: Prometheus)  | 9090                          |
| ninhydrin-monitoring-ui     | Monitoring Web UI (default: Grafana) | 3000                          |

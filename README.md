# Email-sender

## How to run

1. Each service have "env" file, you need to fill it before running the project. There are two running modes: pipes and events. If you run events mode, you will need to run all services. If you run pipes mode, you will need only rest-api service, which demonstrates pipes-and-filters.
2. `docker-compose up -d` run all necessary containers
3. It is possible to see all metrics, default port is 3000 for Grafana
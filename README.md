# Journeys API for Ozon Voice Assistant project
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/ozonva/ova-journey-api.svg)](https://github.com/ozonva/ova-journey-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/ozonva/ova-journey-api)](https://goreportcard.com/report/github.com/ozonva/ova-journey-api)
[![GitHub license](https://img.shields.io/github/license/ozonva/ova-journey-api.svg)](https://github.com/ozonva/ova-journey-api/LICENSE)
[![GitHub pull-requests](https://img.shields.io/github/issues-pr/ozonva/ova-journey-api.svg)](https://github.com/ozonva/ova-journey-api.svg/pull/)
[![Test & Lint](https://github.com/ozonva/ova-journey-api/actions/workflows/test_and_lint.yml/badge.svg?branch=main)](https://github.com/ozonva/ova-journey-api/actions/workflows/test_and_lint.yml)

Ozon Voice Assistant Project Documentation is here https://github.com/ozonva/docs

## Services
### API
#### GRPC
http://localhost:8081
### JSON Gateway
http://localhost:8080
#### Metrics for Prometheus
http://localhost:9100

### With UI
#### Swagger UI 
http://localhost:8080
#### Jaeger UI
http://localhost:16686
#### Kafka UI
http://localhost:8082
#### Prometheus UI
http://localhost:9090
#### Graylog UI
http://localhost:9000

## Commands
+ ```make all``` - clean from binary builds, run all tests and linters, then build the application with code generation
+ ```make run```- run the application
+ ```make build``` - build the application into `bin` directory
+ ```make deps``` - install all project dependencies
+ ```make bin-deps``` - install binary tools (mockgen, proto generators) into `bin` directory
+ ```make clean``` - clean from binary builds, delete `bin` directory
+ ```make lint``` - check code by golangci-lint and gosec linters
+ ```make vendor-proto``` - download external proto files (Google)
+ ```make generate``` - use code generation 
+ ```make test``` - run all tests
+ ```make documentation``` - make documentation and run documentation server on http://localhost:6060/pkg/github.com/ozonva/ova-journey-api/?m=all
+ ```make docker-build``` - build docker application environment using `docker-compose`



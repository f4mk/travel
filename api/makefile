SHELL := /bin/bash

CONFIG_PATH = config/.env

include ${CONFIG_PATH}

GO_PROVIDE_CONFIG = -ldflags="-X 'main.configPath=${CONFIG_PATH}'"

#CREATE SERVICE IMAGE
image:
	docker build \
		-f ./dockerfile \
		-t travel-api \
		--build-arg BUILD_REF=`git rev-parse HEAD` \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		--build-arg CONFIG_PATH=${CONFIG_PATH} \
		--build-arg API_CERT_FILE=${API_CERT_FILE} \
		--build-arg API_KEY_FILE=${API_KEY_FILE} \
		--build-arg AUTH_PRIVATE_KEY_FILE=${AUTH_PRIVATE_KEY_FILE} \
		--build-arg AUTH_PUBLIC_KEY_FILE=${AUTH_PUBLIC_KEY_FILE} \
		.
#DEKETE ALL DOCKER DANGLING IMAGES
docker-rmi:
	docker rmi `docker images -f dangling=true -q`

#RUN DOCKER IMAGE
run-docker:
	docker run -p ${HOST_PORT}:${API_PORT} app

#UPDATE SERVICE DEPENDENCIES
tidy:
	go mod tidy

#BUILD SERVICE
build:
	go build ${GO_PROVIDE_CONFIG} -o app ./cmd/app/main.go

#RUN SERVICE LOCALLY
run:
	go run ${GO_PROVIDE_CONFIG} ./cmd/app/main.go

#PULL DOCKER IMAGES
pull: 
	docker pull postgres:15.3
	docker pull openzipkin/zipkin:2.24.2
	docker pull alpine:3.18.2
	docker pull golang:1.20.6-alpine3.18

#START APP FROM SCRATCH
all: pull image
	docker compose --env-file ${CONFIG_PATH} up

#RUN APP WITH REBUILD
up: image
	docker compose --env-file ${CONFIG_PATH} up

#RUN APP
compose:
	docker compose --env-file ${CONFIG_PATH} up

#GENERATE SSL CERTIFICATE FOR HTTPS
cert:
	mkdir -p certs
	openssl genrsa -out ${API_KEY_FILE} 4096
	openssl req -new -x509 -days 1825 -key ${API_KEY_FILE} -out ${API_CERT_FILE}

#MIGRATE (brew install golang-migrate)
migrate:
	go run ${GO_PROVIDE_CONFIG}  ./cmd/admin/main.go -migrate
# migrate -path ./cmd/admin/migration/sql -database \
# "postgresql://postgres:password@0.0.0.0:5432/travel?sslmode=disable" -verbose up

#GENERATE JWT TOKEN FOR SPECIFIED KEY.PEM
gen-jwt:
	go run ${GO_PROVIDE_CONFIG} ./cmd/admin/main.go -tokengen \
		-kid="9911fa97-db1d-4e15-a683-a85f860f565b" \
		-role="ADMIN"

#GENERATE JWT ALL TOKENS FOR KEY.PEM FILES
gen-jwt-all:
	go run ${GO_PROVIDE_CONFIG} ./cmd/admin/main.go -tokengen-all -role="ADMIN"

#GENERATE JWT TOKEN/PRIVATE KEY FOR AUTH
gen-key:
	go run ${GO_PROVIDE_CONFIG} ./cmd/admin/main.go -keygen
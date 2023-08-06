SHELL := /bin/bash

API_CONFIG_PATH := $(shell cd ./api && $(MAKE) -s config-path)

api-image:
	$(MAKE) -C ./api image

front-image:
	$(MAKE) -C ./frontend image

#DELETE ALL DOCKER DANGLING IMAGES
docker-rmi:
	docker rmi `docker images -f dangling=true -q`


#PULL DOCKER IMAGES
pull: 
	docker pull postgres:15.3
	docker pull openzipkin/zipkin:2.24.2
	docker pull alpine:3.18.2
	docker pull golang:1.20.6-alpine3.18
	docker pull node:18.17

.PHONY: compose
compose:
	docker-compose \
		--env-file api/${API_CONFIG_PATH} \
		-f docker-compose.yml \
		-f api/docker-compose.yml \
		-f frontend/docker-compose.yml \
		up

#START APP FROM SCRATCH
all: pull api-image front-image compose

#RUN APP WITH REBUILD
up: api-image front-image compose

#GENERATE SSL CERTIFICATE FOR HTTPS
cert:
	mkdir -p certs
	openssl genrsa -out ${API_KEY_FILE} 4096
	openssl req -new -x509 -days 1825 -key ${API_KEY_FILE} -out ${API_CERT_FILE}

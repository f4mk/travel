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
	$(MAKE) -C ./api pull
	docker pull node:18.17-alpine

.PHONY: compose-up
compose-up:
	docker-compose \
		--env-file api/${API_CONFIG_PATH} \
		-f haproxy/docker-compose.yml \
		-f api/docker-compose.yml \
		-f frontend/docker-compose.yml \
		up

.PHONY: compose-down
compose-down:
	docker-compose \
		--env-file api/${API_CONFIG_PATH} \
		-f haproxy/docker-compose.yml \
		-f api/docker-compose.yml \
		-f frontend/docker-compose.yml \
		up

#START APP FROM SCRATCH
all: pull api-image front-image compose-up

#RUN APP WITH REBUILD
up: api-image front-image compose-up

#GENERATE SSL CERTIFICATE FOR HTTPS
cert:
	mkdir -p certs
	openssl genrsa -out ${API_KEY_FILE} 4096
	openssl req -new -x509 -days 1825 -key ${API_KEY_FILE} -out ${API_CERT_FILE}

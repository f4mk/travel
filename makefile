SHELL := /bin/bash

API_CONFIG_PATH := $(shell cd ./backend && $(MAKE) -s config-path)

backend-image:
	$(MAKE) -C ./backend image


front-image:
	$(MAKE) -C ./frontend image

#DELETE ALL DOCKER DANGLING IMAGES
docker-rmi:
	docker rmi `docker images -f dangling=true -q`


#PULL DOCKER IMAGES
pull: 
	$(MAKE) -C ./backend pull
	docker pull node:18.17-alpine
	docker pull rabbitmq:3.12:management

.PHONY: compose-up
compose-up:
	docker compose \
		--env-file backend/${API_CONFIG_PATH} \
		-f haproxy/docker-compose.yml \
		-f backend/docker-compose.yml \
		-f frontend/docker-compose.yml \
		-f rabbit/docker-compose.yml \
		up

.PHONY: compose-down
compose-down:
	docker compose \
		--env-file backend/${API_CONFIG_PATH} \
		-f haproxy/docker-compose.yml \
		-f backend/docker-compose.yml \
		-f frontend/docker-compose.yml \
		up

#START APP FROM SCRATCH
all: pull backend-image front-image compose-up

#RUN APP WITH REBUILD
up: backend-image front-image compose-up

#GENERATE SSL CERTIFICATE FOR HTTPS
cert:
	mkdir -p certs
	openssl genrsa -out ${API_KEY_FILE} 4096
	openssl req -new -x509 -days 1825 -key ${API_KEY_FILE} -out ${API_CERT_FILE}

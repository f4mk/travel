SHELL := /bin/bash

image:
	docker rmi travel-front:latest
	docker build -t travel-front -f ./dockerfile .

up: image
	docker compose up
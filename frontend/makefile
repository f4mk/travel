SHELL := /bin/bash

image:
	docker rmi travel-static:latest || true
	docker build -t travel-static -f ./dockerfile .

up: image
	docker compose up
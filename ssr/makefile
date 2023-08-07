SHELL := /bin/bash

image:
	docker rmi travel-ssr:latest || true
	docker build -t travel-ssr -f ./dockerfile .

up: image
	docker compose up
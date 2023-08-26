SHELL := /bin/bash

image:
	docker rmi haproxy-volume:latest || true
	docker build \
		-f ./dockerfile \
		-t haproxy-volume .

SHELL := /bin/bash

# Dont forget to add $GOPATH and sourse you .zshrc
# .zshrc content:
# export GOPATH=$HOME/go
# export PATH=/usr/local/bin:$PATH:$GOPATH/bin
# eval $(/opt/homebrew/bin/brew shellenv)
# command:
# source ~/.zshrc

# For load testing use hey:
# brew install hey
# hey -n 100000 -c 200 http://localhost:8080/

# For monitoring service resourse utilization install expvarmon
# go install github.com/divan/expvarmon@latest
# expvarmon -ports="8081"
# expvarmon -ports="8081" -vars="goroutines,mem:memstats.Alloc"

# For typegen install oapi-codegen globally
# go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest


################# CONFIG ##################

CONFIG_PATH = travel-api/config/.env
GO_API_PREFIX_PATH = travel-api

include ${CONFIG_PATH}

# TODO: refactor this
GO_PROVIDE_CONFIG = -ldflags="-X 'main.configPath=${CONFIG_PATH}'"
GO_PROVIDE_MIGRATION = -ldflags="-X 'main.configPath=${CONFIG_PATH}' \
	-X 'main.migrationsPath=file://${GO_API_PREFIX_PATH}/cmd/admin/migration/sql'"

config-path:
	@echo ${CONFIG_PATH}


###########################################

################# IMAGES ##################
#CREATE SERVICE IMAGE
image:
	docker rmi travel-api:latest || true
	docker build \
		-f ./dockerfile.api \
		-t travel-api \
		--build-arg BUILD_REF=`git rev-parse HEAD` \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		--build-arg CONFIG_PATH=${CONFIG_PATH} \
		--build-arg MAIL_PRIVATE_KEY_PATH=${MAIL_PRIVATE_KEY_PATH} \
		--build-arg MAIL_PUBLIC_KEY_PATH=${MAIL_PUBLIC_KEY_PATH} \
		--build-arg AUTH_PRIVATE_KEY_FILE=${AUTH_PRIVATE_KEY_FILE} \
		--build-arg AUTH_PUBLIC_KEY_FILE=${AUTH_PUBLIC_KEY_FILE} \
		.

#CREATE CRON IMAGE
image-cron:
	docker rmi travel-api-cron:latest || true
	docker build \
		-f ./dockerfile.cron \
		-t travel-api-cron \
		--build-arg CONFIG_PATH=${CONFIG_PATH} \
		.

#CREATE METRICS IMAGE
image-metrics:
	docker rmi travel-api-metrics:latest || true
	docker build \
		-f ./dockerfile.metrics \
		-t travel-api-metrics \
		--build-arg BUILD_REF=`git rev-parse HEAD` \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		--build-arg CONFIG_PATH=${CONFIG_PATH} \
		.

#DELETE ALL DOCKER DANGLING IMAGES
docker-rmi:
	docker rmi `docker images -f dangling=true -q`

#PULL DOCKER IMAGES
pull: 
	docker pull postgres:15.3
	docker pull alpine:3.18.2
	docker pull golang:1.20.6-alpine3.18
	docker pull redis:6.2-alpine

###########################################

################# GO CODE #################

#UPDATE SERVICE DEPENDENCIES
tidy:
	go mod tidy

###########################################

################## TOOLS ##################

#GENERATE SSL CERTIFICATE FOR HTTPS
cert:
	mkdir -p certs
	openssl genrsa -out ${API_KEY_FILE} 4096
	openssl req -new -x509 -days 1825 -key ${API_KEY_FILE} -out ${API_CERT_FILE}

#MIGRATE (brew install golang-migrate)
migrate:
	go run ${GO_PROVIDE_MIGRATION}  ./travel-api/cmd/admin/main.go -migrate
# migrate -path ./travel-api/cmd/admin/migration/sql -database \
# "postgresql://postgres:password@0.0.0.0:5432/travel?sslmode=disable" -verbose up

#GENERATE JWT TOKEN FOR SPECIFIED KEY.PEM
gen-jwt:
	go run ${GO_PROVIDE_CONFIG} ./travel-api/cmd/admin/main.go -tokengen \
		-kid="9911fa97-db1d-4e15-a683-a85f860f565b" \
		-role="ADMIN"

#GENERATE JWT ALL TOKENS FOR KEY.PEM FILES
gen-jwt-all:
	go run ${GO_PROVIDE_CONFIG} ./travel-api/cmd/admin/main.go -tokengen-all -role="ADMIN"

#GENERATE JWT TOKEN/PRIVATE KEY FOR AUTH
gen-key:
	go run ${GO_PROVIDE_CONFIG} ./travel-api/cmd/admin/main.go -keygen

#GENERATE TYPES FROM OPENAPI SPEC
#https://pkg.go.dev/github.com/go-playground/validator/v10#section-readme
#npm install -g @apidevtools/swagger-cli

types:
	go run travel-api/cmd/tools/typegen/main.go

#LINT CODE
lint: 
	golangci-lint run
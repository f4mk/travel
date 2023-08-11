# TRAVEL API SERVICE

## HOW TO RUN

**Prereqs**
- To run this service you need to install the following:
```brew install docker```
```brew install docker-compose```
- You also must be able to run makefile commands

**Run**
- for production: run ```make all``` for pulling all the necessary images, building project image itself, and launching it with docker-compose as if it would be in production environment
- for development: 
  - run ```make up``` to rebuild and launch the whole application with docker-compose
  - *run ```make up``` to launch the whole application without rebuilding images*
  - run ```make run``` to launch only the Go server itself

## HOW TO CONFIGURE

**Config files**
- all the configuration resides in ```config/.env``` file by default
- use ```config/.env.example``` as a template
- path to ```.env``` file is set in ```makefile``` in the root folder, DO NOT MOVE IT
as it depends on the upper layer confing

**Certificates**
- path to certificates and JWT keys is set in ```config/.env``` by default
- ```make cert``` command will generate certificates for TLS (you prolly wont need it, as TLS is terminated in haproxy in a normal use case)
- to be able to generate JWT tokens you need to create private key by running ```make gen``` command, in case you need an admin token